# HTTP

__1) http.HandleFunc__ [01_handlefunc.go]

```go
http.HandleFunc("/", handlerFunction)
http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
```

`ListenAndServe` - Handler is typically nil, in which case the `DefaultServeMux` is used.

In core `server.go`:

```go
// HandleFunc registers the __handler function__ for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

__2) http.HandlerFunc__ [02_handlerfunc.go]

```go
// The HandlerFunc type is an __adapter__ to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ResponseWriter, *Request)
```

`ListenAndServe` accepts type `HandlerFunc` as the second agument.

`HandlerFunc` type implements `ServeHTTP` func, therefore `HandlerFunc` satisfies the `Handler` interface:

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

The `http.Handler` interface:

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

Downside of this approach is the limited routing - we'd have to implement it within the `h`... 
Better - `ServeMux`

__3) http.ServeMux__ [03_servemux.go]

The `http.ServeMux` is itself an `http.Handler`, so it can be passed into `http.ListenAndServe`. When it receives a request it will check if the request’s path is prefixed by any of its known paths, choosing the longest prefix match it can find. We use the / endpoint as a catch-all to catch any requests to unknown endpoints. 

`http.ServeMux` has both `Handle` and `HandleFunc` methods. These do the same thing, except that Handle takes in an `http.Handler` while `HandleFunc` merely takes in a function, implicitly wrapping it just as `http.HandlerFunc` does.

### Other muxes

There are numerous replacements for `http.ServeMux` like `gorilla/mux` which give you things like automatically pulling variables out of paths, easily asserting what http methods are allowed on an endpoint, and more. Most of these replacements will implement `http.Handler` like `http.ServeMux` does, and accept `http.Handlers` as arguments.

__4) ServeMux__

`numberDumper` implements `http.Handler`, and can be passed into the `http.ServeMux` multiple times to serve multiple endpoints.

__5) middleware__

`curl localhost:8080/helloworld -d '{"name": "john"}'`

https://www.safaribooksonline.com/library/view/building-microservices-with/9781786468666/30f61396-c1f9-47ab-b5d1-1ed431ce69db.xhtml

Split request handling (validation & returning a response)

```go
type validationHandler struct {
	next http.Handler
}
```
Both `validationHandler` and `helloWorldHandler` implement the `ServeHTTP` method satisfying the `Handler` interface.
The `validationHandler` must have a reference to the next in the chain: 
```go
handler := newValidationHandler(newHelloWorldHandler())
http.Handle("/helloworld", handler)
```

__6) context__

https://www.safaribooksonline.com/library/view/building-microservices-with/9781786468666/fbc92048-e0df-4a3a-a986-4735e7f6003a.xhtml



> TODO, next: https://www.safaribooksonline.com/library/view/security-with-go/9781788627917/f0fe0d1e-470a-4cbf-a1f0-2c0c160ec4ce.xhtml