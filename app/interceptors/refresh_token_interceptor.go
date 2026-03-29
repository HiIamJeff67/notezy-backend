package interceptors

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	ratelimiter "notezy-backend/app/lib/ratelimiter"
	logs "notezy-backend/app/monitor/logs"
	traces "notezy-backend/app/monitor/traces"
	types "notezy-backend/shared/types"
)

// use the reusable buffer pool for refreshing the access token
var refreshAccessTokenReusableBufferPool *ratelimiter.ReusableBufferPool = ratelimiter.NewReusableBufferPool()

// To rewrite the response with adding additional field of `newAccessToken` and `newCSRFToken`,
// Note : It should be placed below the `AuthMiddleware`,
// so that it can access the `AccessToken` and `CSRFToken` in the context field
func RefreshTokenInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentBufferPool := refreshAccessTokenReusableBufferPool.Get()
		defer func() {
			currentBufferPool.Reset()
			refreshAccessTokenReusableBufferPool.Put(currentBufferPool)
		}()

		writer := ratelimiter.NewResponseWriter(ctx.Writer, currentBufferPool)
		ctx.Writer = writer // replace the response writer with the declared writer here
		// so that we can re-write the response after the controller sent the response !!
		// we can successfully do this since the interceptor inheritent the gin.ResponseWriter interface,
		// and it also implement Write() and WriteString() methods.
		// Note: they write the content into the `originalBody`,
		// so the field of `originalBody` is the original content from the controllers

		ctx.Next() // execute the following first

		if writer.IsTimeout {
			return
		}

		if writer.ResponseWriter.Written() || writer.Status() >= 400 {
			writer.FlushToOriginalWriter()
			return
		}

		if ctx.Writer.Status() >= 400 {
			writer.FlushToOriginalWriter()
			return
		}

		IsNewTokens, exception := contexts.GetAndConvertContextFieldToBoolean(ctx, types.ContextFieldName_IsNewTokens)
		if exception != nil || IsNewTokens == nil || !*IsNewTokens {
			writer.FlushToOriginalWriter()
			return
		}

		var originalResponse map[string]interface{}
		if err := json.Unmarshal(writer.Body.Bytes(), &originalResponse); err != nil {
			writer.FlushToOriginalWriter()
			return
		}

		accessToken, exists := ctx.Get(types.ContextFieldName_AccessToken.String())
		if !exists {
			writer.FlushToOriginalWriter()
			return
		}
		accessTokenStr, ok := accessToken.(string)
		if !ok {
			writer.FlushToOriginalWriter()
			return
		}
		csrfToken, exist := ctx.Get(types.ContextFieldName_CSRFToken.String())
		if !exist {
			writer.FlushToOriginalWriter()
			return
		}
		csrfTokenStr, ok := csrfToken.(string)
		if !ok {
			writer.FlushToOriginalWriter()
			return
		}

		cookies.AccessTokenCookieHandler.Set(ctx, accessTokenStr)
		originalResponse[types.RefreshableResponseFieldName_NewAccessToken.String()] = accessTokenStr
		originalResponse[types.RefreshableResponseFieldName_NewCSRFToken.String()] = csrfTokenStr
		logs.Info(traces.GetTrace(0).FileLineString(), "new CSRF token: ", csrfTokenStr)
		modifiedResponse, err := json.Marshal(originalResponse)
		if err != nil {
			writer.FlushToOriginalWriter()
			return
		}

		writer.Mutex.Lock()
		defer writer.Mutex.Unlock()

		destination := writer.ResponseWriter.Header()
		for key, val := range writer.Headers {
			destination[key] = val
		}

		if writer.Code > 0 {
			writer.ResponseWriter.WriteHeader(writer.Code)
		}

		writer.ResponseWriter.Header().Set("Content-Length", string(rune(len(modifiedResponse))))
		writer.ResponseWriter.Write(modifiedResponse)
	}
}
