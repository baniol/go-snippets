# Bytes, buffers, io

* https://medium.com/golangspec/introduction-to-bufio-package-in-golang-ad7d1877f762 - !!
* https://medium.com/golangspec/in-depth-introduction-to-bufio-scanner-in-golang-55483bb689b4
* https://medium.com/stupid-gopher-tricks/streaming-data-in-go-without-buffering-3285ddd2a1e5 - !!
* https://hashrocket.com/blog/posts/go-performance-observations - benchmarks
* https://zensrc.wordpress.com/2016/06/15/golang-buffer-vs-non-buffer-in-io/
* https://groups.google.com/forum/#!topic/golang-nuts/mOvX0bmJoeI - seek
* https://stackoverflow.com/questions/39319024/builtin-append-vs-bytes-buffer-write - benchmarks
* https://syslog.ravelin.com/bytes-buffer-i-thought-you-were-my-friend-4148fd001229
* https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/
* https://github.com/stabbycutyou/buffstreams
* https://github.com/aybabtme/portproxy/blob/master/portproxy.go - simple proxy

`bytes.Buffer` is a wrapper around `[]byte`. It adds convenience by implementing popular interfaces like `io.Reader`, ` ` etc.


====
https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/r

Go has a function called ioutil.Readall, which is defined:

`func ReadAll(r io.Reader) ([]byte, error)`

Use of ioutil.ReadAll is almost always a mistake.

===
https://medium.com/go-walkthrough/go-walkthrough-io-package-8ac5e95a9fbd

Go is a programming language built for working with bytes. Whether you have lists of bytes, streams of bytes, or individual bytes, Go makes it easy to process. From these simple primitives we build our abstractions and services.

The basic construct for reading bytes from a stream is the `Reader interface`:
```go
type Reader interface {
        Read(p []byte) (n int, err error)
}
```
This interface is implemented throughout the standard library by everything from network connections to files to wrappers for in-memory slices.

One problem with the Reader interface is that it comes with some subtle rules. First, it returns an io.EOF error as a normal part of usage when the stream is done.

To read 8 bytes you only need to do this:

```go
buf := make([]byte, 8)
if _, err := io.ReadFull(r, buf); err == io.EOF {
        return io.ErrUnexpectedEOF
} else if err != nil {
        return err
}
```

## Multiwriter

The MultiWriter comes in handy in this case:

`func MultiWriter(writers ...Writer) Writer`

The name is a bit confusing since it’s not the writer version of MultiReader. Whereas MultiReader concatenates several readers into one, the MultiWriter returns a writer that duplicates each write to multiple writers.

I use MultiWriter extensively in unit tests where I need to assert that a service is logging properly:

```go
type MyService struct {
        LogOutput io.Writer
}
...
var buf bytes.Buffer
var s MyService
s.LogOutput = io.MultiWriter(&buf, os.Stderr)
```
Using a MultiWriter allows me to verify the contents of buf while also seeing the full log output in my terminal for debugging.

### Optimizing string writes
There are a lot of writers in the standard library that have a WriteString() method which can be used to improve write performance by not requiring an allocation when converting a string to a byte slice. You can take advantage of this optimization by using the io.WriteString() function.

The function is simple. It first checks if the writer implements a WriteString() method and uses it if available. Otherwise it falls back to copying the string to a byte slice and using the Write() method.

### Moving around within streams
Streams are usually a continuous flow of bytes from beginning to end but there are a few exceptions. A file, for example, can be operated on as a stream but you can also jump to a specific position within the file.

The Seeker interface is provided to jump around in a stream:
```go
type Seeker interface {
        Seek(offset int64, whence int) (int64, error)
}
```
There are 3 ways to jump around: move from on the current position, move from the beginning, and move from the end. You specify the mode of movement using the whence argument. The offset argument specifies how many bytes to move by.