package interceptors

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type Interceptor struct {
	gin.ResponseWriter // inheritent the gin.ResponseWriter
	originalBody       *bytes.Buffer
}

func (w *Interceptor) Write(b []byte) (int, error) {
	// we use `originalBody` to write the original content from the controllers
	return w.originalBody.Write(b)
}

func (w *Interceptor) WriteString(s string) (int, error) {
	// we use `originalBody` to write the original content string from the controllers
	return w.originalBody.WriteString(s)
}
