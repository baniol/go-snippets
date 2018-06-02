# Go snippets (http)

## References

* https://www.alexedwards.net/blog/a-recap-of-request-handling
* https://www.safaribooksonline.com/library/view/building-microservices-with/9781786468666/
* https://www.safaribooksonline.com/library/view/security-with-go/9781788627917/
* https://cryptic.io/go-http/
* https://www.rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/

## http bits and pieces

`handleFunc`
```go
// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

```go
// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }
```

`handlerFunc`
```go
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

```go
// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// Except for reading the body, handlers should not modify the
// provided Request.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and either closes the network connection or sends an HTTP/2
// RST_STREAM, depending on the HTTP protocol. To abort a handler so
// the client sees an interrupted response but the server doesn't log
// an error, panic with the value ErrAbortHandler.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

> Processing HTTP requests with Go is primarily about two things: __ServeMuxes__ and __Handlers__.
> Handlers are responsible for writing response headers and bodies.
> Go's HTTP package ships with a few functions to generate common handlers, such as `FileServer`, `NotFoundHandler` and `RedirectHandler`.

## Handler

__handler_1.go__

New type `helloHandler` implements `http.Handler` interface with `ServeHTTP` method and can be passed to `http.ListenAndServe` or used with the `Go muxer` (`http.Handle(pattern, handler)` function).

__handler_2.go__

The http package provides a helper function, `http.HandlerFunc`, which wraps a function which has the signature `func(w http.ResponseWriter, r *http.Request)`, returning an `http.Handler` which will call it in all cases.

This example behaves exactly like the previous one, but uses `http.HandlerFunc` instead of defining a new type.

`http.HandlerFun`c type is used to coerce a normal function into satisfying the `Handler` interface.

Any function which has the signature `func(http.ResponseWriter, *http.Request)` can be converted into a `HandlerFunc` type. This is useful because `HandleFunc` objects come with an inbuilt `ServeHTTP` method which – rather cleverly and conveniently – executes the content of the original function.

__handleFunc.go__

```go
mux := http.NewServeMux()
h := http.HandlerFunc(mainHandler)
mux.Handle("/", h)
http.ListenAndServe(":8080", mux)
```

In fact, converting a function to a `HandlerFunc` type and then adding it to a `ServeMux` like this is so common that Go provides a shortcut: the `mux.HandleFunc` method:

```go
http.HandleFunc("/", mainHandler)
http.ListenAndServe(":8080", nil)
```

Here, `http.HandleFunc` takes an http __handler function__: `func(ResponseWriter, *Request)` as a second parameter. The `http.HandleFunc` uses `ServeMux`: `DefaultServeMux.HandleFunc(pattern, handler)`.

Inside http package there is default `ServeMux` implementation that is stored as variable `DefaultServeMux`. Using this `DefaultServeMux` we can skip to define our own `ServeMux` implementation and directly utilize it by calling function `Handle` or `HandleFunc` inside http package.

The `ServeMux` type also has a ServeHTTP method, meaning that it too satisfies the Handler interface (and can be passed to `http.ListenAndServe`).

For me it simplifies things to think of a `ServeMux` as just being a special kind of handler, which instead of providing a response itself passes the request on to a second handler.

__handler_closure.go__

If we want to pass some variable into the handler, we can put the handler function logic into a closure.

The `timeHandler` function now has a subtly different role. Instead of coercing the function into a handler (like we did previously), we are now using it to return a handler.

First it creates `fn`, an anonymous function which accesses ‐ or closes over – the `format` variable forming a closure. Regardless of what we do with the closure it will always be able to access the variables that are local to the scope it was created in – which in this case means it'll always have access to the `format` variable.

Secondly our closure has the signature func(http.ResponseWriter, *http.Request). As you may remember from earlier, this means that we can convert it into a HandlerFunc type (so that it satisfies the Handler interface). Our timeHandler function then returns this converted closure.

In this example we've just been passing a simple string to a handler. But in a real-world application you could use this method to pass database connection, template map, or any other application-level context. It's a good alternative to using global variables, and has the added benefit of making neat self-contained handlers for testing.

You might also see this same pattern written as:

```go
func timeHandler(format string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(format)
    w.Write([]byte("The time is: " + tm))
  })
}
```

Or using an implicit conversion to the `HandlerFunc` type on return:

```go
func timeHandler(format string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(format)
    w.Write([]byte("The time is: " + tm))
  }
}
```

__servemux.go__

`http.ServeMux` maps the url to handler.

The `http.ServeMux` is itself an `http.Handler`, so it can be passed into `http.ListenAndServe`. When it receives a request it will check if the request’s path is prefixed by any of its known paths, choosing the longest prefix match it can find. We use the / endpoint as a catch-all to catch any requests to unknown endpoints. 

`http.ServeMux` has both `Handle` and `HandleFunc` methods. These do the same thing, except that Handle takes in an `http.Handler` while `HandleFunc` merely takes in a function, implicitly wrapping it just as `http.HandlerFunc` does.

__servemux2.go__

When I say that the http package is composable I mean that it is very easy to create re-usable pieces of code and glue them together into a new working application. The `http.Handler` interface is the way all pieces communicate with each other.

`numberDumper` implements `http.Handler`, and can be passed into the `http.ServeMux` multiple times to serve multiple endpoints.

## Middleware

> https://www.alexedwards.net/blog/making-and-using-middleware !

Into the closure we can pass the next handler in the chain as a variable, and then transfer control to this next handler by calling it's `ServeHTTP()` method.

__Pattern for constructing middleware__

```go
func exampleMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Our middleware logic goes here...
    next.ServeHTTP(w, r)
  })
}
```

__post_body.go__

```go
http.HandleFunc("/city", mainHandler)
```

Here, `http.HandleFunc` takes an http __handler function__: `func(ResponseWriter, *Request)` as a second parameter.

__post_body_middleware.go__

For middleware, we have to wrap the handler function, so it can be passed to `http.Handle`, which takes the __`Handler` interface__ as a second parameter. This handler can be wrapped in middlewares, like `http.Handle("/city", contentTypeValidator(mainLogic))`.