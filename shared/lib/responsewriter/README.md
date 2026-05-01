# ResponseWriter Library

## Overview

`app/lib/responsewriter` provides a Gin-compatible buffered response writer.

It wraps `gin.ResponseWriter` so handlers/middlewares can write to memory first,
then decide when and how to flush to the original writer. This is useful for
response interception, timeout handling, and post-processing.

## Files

- `response_writer.go`: buffered writer implementation and flush utilities.

## Public API

```go
type ResponseWriter
func NewResponseWriter(
	responseWriter gin.ResponseWriter,
	Body *bytes.Buffer,
) *ResponseWriter
func (rw *ResponseWriter) Header() http.Header
func (rw *ResponseWriter) WriteHeader(code int)
func (rw *ResponseWriter) Write(data []byte) (int, error)
func (rw *ResponseWriter) WriteString(s string) (int, error)
func (rw *ResponseWriter) WriteHeaderNow()
func (rw *ResponseWriter) Status() int
func (rw *ResponseWriter) Size() int
func (rw *ResponseWriter) FlushToOriginalWriter() error
func (rw *ResponseWriter) FreeBuffer()
```

## Behavior Notes

- `Write` and `WriteString` buffer output into memory (`Body`) instead of writing directly.
- `WriteHeader` is guarded to avoid duplicate writes and ignores writes after timeout.
- `Status` is overridden so downstream middleware can read the true status code.
- `FreeBuffer` clears and detaches the current buffer for memory reuse.
- Mutex-protected methods help prevent race conditions around writes/flush states.

## Example (Gin Middleware)

```go
func BufferedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		buf := &bytes.Buffer{}
		writer := responsewriter.NewResponseWriter(ctx.Writer, buf)
		ctx.Writer = writer

		ctx.Next()

		writer.Mutex.Lock()
		defer writer.Mutex.Unlock()

		// append an audit header before final flush
		writer.Header().Set("X-Buffered", "true")
		_ = writer.FlushToOriginalWriter()
	}
}
```

## Project Usage Example

This package is used by timeout and interceptor flows in:

- `app/middlewares/timeout_middleware.go`
- `app/interceptors/shareable_response_writer_interceptor.go`
- `app/interceptors/embedded_interceptor.go`
- `app/interceptors/refresh_token_interceptor.go`

## File Structure

```text
app/lib/responsewriter/
├── README.md
├── LICENSE.md
└── response_writer.go
```
