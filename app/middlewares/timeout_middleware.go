package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	lib "notezy-backend/app/lib"
	logs "notezy-backend/app/logs"
	types "notezy-backend/shared/types"
)

// use reusable buffer pool for timeout response writer to storing the current response of the handlers
var timeoutReusableBufferPool *lib.ReusableBufferPool = lib.NewReusableBufferPool()

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		originalWriter := ctx.Writer

		currentBufferPool := timeoutReusableBufferPool.Get()
		currentBufferPool.Reset()
		defer func() {
			currentBufferPool.Reset()
			timeoutReusableBufferPool.Put(currentBufferPool)
		}()

		writer := lib.NewResponseWriter(originalWriter, currentBufferPool)
		ctx.Writer = writer

		ctxCopy := ctx.Copy()
		ctxCopy.Writer = writer

		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer cancel()

		ctxCopy.Request = ctxCopy.Request.WithContext(timeoutCtx)

		// ctx     uses originalWriter
		// ctxCopy uses writer (passing to the next handlers to make them write response in the buffer)

		done := make(chan struct{}, 1)
		panicChannel := make(chan types.PanicInfo, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChannel <- types.PanicInfo{
						Value: p,
						Stack: debug.Stack(),
					}
				}
			}()

			ctx.Next()
			done <- struct{}{}
		}()

		select {
		case panicInfo := <-panicChannel:
			writer.Mutex.Lock()
			writer.FreeBuffer() // clear the buffer, this will destroy the context field stored by other middlewares
			writer.Mutex.Unlock()
			ctx.Writer = originalWriter

			if gin.IsDebugging() {
				ctx.Error(fmt.Errorf("%v", panicInfo.Value))
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(ctx.Writer, "panic caught: %v\n", panicInfo.Value)
				ctx.Writer.Write([]byte("Panic stack trace:\n"))
				ctx.Writer.Write(panicInfo.Stack)
			}

			panic(panicInfo.Value)
		case <-done:
			writer.Mutex.Lock()
			defer writer.Mutex.Unlock()
			if writer.IsTimeout || writer.ResponseWriter.Written() {
				return
			}

			destination := writer.ResponseWriter.Header()
			for key, val := range writer.Headers {
				destination[key] = val
			}

			if writer.Code != 0 {
				writer.ResponseWriter.WriteHeader(writer.Code)
			}

			if currentBufferPool.Len() > 0 {
				if _, err := writer.ResponseWriter.Write(currentBufferPool.Bytes()); err != nil {
					panic(err)
				}
			}

			return
		case <-timeoutCtx.Done():
			logs.Info("Timeout (timeoutCtx.Done())")
			writer.Mutex.Lock()
			writer.IsTimeout = true
			writer.FreeBuffer() // clear the buffer, this will destroy the context field stored by other middlewares
			writer.Mutex.Unlock()

			if !writer.ResponseWriter.Written() {
				writer.ResponseWriter.Header().Set("X-Request-Timeout", timeout.String())
				exception := exceptions.Timeout(timeout)
				writer.ResponseWriter.WriteHeader(exception.HTTPStatusCode)
				timeoutResponseBody, err := exception.GetResponseJSONBytes()
				if err != nil {
					panic(err)
				}
				if _, err := writer.ResponseWriter.Write(timeoutResponseBody); err != nil {
					panic(err)
				}
			}

			return
		case <-time.After(timeout):
			logs.Info("Timeout (time.After)")
			writer.Mutex.Lock()
			writer.IsTimeout = true
			writer.FreeBuffer() // clear the buffer, this will destroy the context field stored by other middlewares
			writer.Mutex.Unlock()

			if !writer.ResponseWriter.Written() {
				writer.ResponseWriter.Header().Set("X-Request-Timeout", timeout.String())
				exception := exceptions.Timeout(timeout)
				writer.ResponseWriter.WriteHeader(exception.HTTPStatusCode)
				timeoutResponseBody, err := exception.GetResponseJSONBytes()
				if err != nil {
					panic(err)
				}
				if _, err := writer.ResponseWriter.Write(timeoutResponseBody); err != nil {
					panic(err)
				}
			}

			return
		}
	}
}
