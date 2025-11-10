package lib

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Body             *bytes.Buffer
	Headers          http.Header
	Mutex            sync.Mutex
	IsTimeout        bool
	HasWrotenHeaders bool
	Code             int
	BufferedSize     int
}

func NewResponseWriter(responseWriter gin.ResponseWriter, Body *bytes.Buffer) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: responseWriter,
		Body:           Body,
		Headers:        make(http.Header),
	}
}

func (rw *ResponseWriter) WriteHeaderNow() {
	if !rw.HasWrotenHeaders {
		if rw.Code == 0 {
			rw.Code = http.StatusOK
		}

		destination := rw.ResponseWriter.Header()
		for key, val := range rw.Headers {
			destination[key] = val
		}
	}
}

func (rw *ResponseWriter) WriteHeader(Code int) {
	rw.Mutex.Lock()
	defer rw.Mutex.Unlock()

	if rw.IsTimeout || rw.HasWrotenHeaders || Code == -1 {
		return
	}

	if Code < 100 || Code > 999 {
		panic(fmt.Sprintf("Invalid http status Code: %d", Code))
	}

	destination := rw.ResponseWriter.Header()
	for key, val := range rw.Headers {
		destination[key] = val
	}

	rw.HasWrotenHeaders = true
	rw.Code = Code
	rw.ResponseWriter.WriteHeader(Code)
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.Headers
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	rw.Mutex.Lock()
	defer rw.Mutex.Unlock()

	if rw.IsTimeout || rw.Body == nil {
		return 0, nil
	}

	n, err := rw.Body.Write(data)
	rw.BufferedSize += n

	return n, err
}

func (rw *ResponseWriter) WriteString(s string) (int, error) {
	n, err := rw.Write([]byte(s))
	rw.BufferedSize += n
	return n, err
}

func (rw *ResponseWriter) Size() int {
	return rw.BufferedSize
}

func (rw *ResponseWriter) FreeBuffer() {
	rw.Body.Reset() // reset the Body first
	rw.BufferedSize = -1
	rw.Body = nil
}

// we must override Status function here,
// otherwise the http status Code returned by gin.Context.ResponseWriter.Status()
// will always be 200 in other custom gin middlewares.
func (rw *ResponseWriter) Status() int {
	if rw.Code == 0 || rw.IsTimeout {
		return rw.ResponseWriter.Status()
	}
	return rw.Code
}

func (rw *ResponseWriter) FlushToOriginalWriter() error {
	rw.Mutex.Lock()
	defer rw.Mutex.Unlock()

	dst := rw.ResponseWriter.Header()
	for key, val := range rw.Headers {
		dst[key] = val
	}

	if rw.Code == 0 {
		rw.ResponseWriter.WriteHeader(http.StatusOK)
	} else if rw.Code > 200 && rw.Code <= 999 {
		rw.ResponseWriter.WriteHeader(rw.Code)
	}

	if rw.Body.Len() > 0 {
		_, err := rw.ResponseWriter.Write(rw.Body.Bytes())
		return err
	}

	return nil
}
