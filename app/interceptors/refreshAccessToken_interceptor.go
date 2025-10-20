package interceptors

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	constants "notezy-backend/shared/constants"
)

// To rewrite the response with adding additional field of `newAccessToken`,
// Note : It should be placed below the `AuthMiddleware`,
// so that it can access the `AccessToken` in the context field
func RefreshAccessTokenInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		interceptor := &Interceptor{
			ResponseWriter: ctx.Writer,
			originalBody:   &bytes.Buffer{},
		}
		ctx.Writer = interceptor // replace the response writer with the declared writer here
		// so that we can re-write the response after the controller send the response !!
		// we can successfully do this since the interceptor inheritent the gin.ResponseWriter interface,
		// and it also implement Write() and WriteString() methods.
		// Note: they write the content into the `originalBody`,
		// so the field of `originalBody` is the original content from the controllers

		ctx.Next() // execute the following first

		if ctx.Writer.Status() >= 400 {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}

		isNew, exception := contexts.GetAndConvertContextFieldToBoolean(ctx, constants.ContextFieldName_IsNewAccessToken)
		if exception != nil || isNew == nil || !*isNew {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}

		var originalResponse map[string]interface{}
		if err := json.Unmarshal(interceptor.originalBody.Bytes(), &originalResponse); err != nil {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}

		accessToken, exists := ctx.Get(constants.ContextFieldName_AccessToken.String())
		if !exists {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}

		accessTokenStr, ok := accessToken.(string)
		if !ok {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}

		cookies.AccessToken.SetCookie(ctx, accessTokenStr)
		originalResponse["newAccessToken"] = accessTokenStr
		modifiedResponse, err := json.Marshal(originalResponse)
		if err != nil {
			interceptor.ResponseWriter.Write(interceptor.originalBody.Bytes())
			return
		}
		ctx.Writer.Header().Set("Content-Length", string(rune(len(modifiedResponse))))

		interceptor.ResponseWriter.Write(modifiedResponse)
	}
}
