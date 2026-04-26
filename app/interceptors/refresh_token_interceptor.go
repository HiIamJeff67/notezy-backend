package interceptors

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	exceptions "notezy-backend/app/exceptions"
	ratelimiter "notezy-backend/app/lib/ratelimiter"
	types "notezy-backend/shared/types"
)

// To add additional field to the response with adding additional field of `newAccessToken` and `newCSRFToken`,
// Note : It should be placed below the `AuthMiddleware`,
// so that it can access the `AccessToken` and `CSRFToken` in the context field
func RefreshTokenInterceptor(responseWriterKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var writer *ratelimiter.ResponseWriter
		existingWriter, exist := ctx.Get(responseWriterKey)
		if !exist || existingWriter == nil {
			exceptions.Context.MissPlacingOrWrongInterceptorOrder(
				"Cannot find the existing response writer, " +
					"please make sure to call the ShareableResponseWriterInterceptor() and pass RefreshTokenInterceptor() as one of the parameters",
			).Log()
			return
		}
		writer = existingWriter.(*ratelimiter.ResponseWriter)

		ctx.Next() // execute the following first

		if writer.IsTimeout {
			return
		}

		if writer.ResponseWriter.Written() || writer.Status() >= 400 {
			return
		}

		if ctx.Writer.Status() >= 400 {
			return
		}

		IsNewTokens, exception := contexts.GetAndConvertContextFieldToBoolean(ctx, types.ContextFieldName_IsNewTokens)
		if exception != nil || IsNewTokens == nil || !*IsNewTokens {
			return
		}

		var originalResponse map[string]interface{}
		if err := json.Unmarshal(writer.Body.Bytes(), &originalResponse); err != nil {
			return
		}

		accessToken, exceptionOfAccessToken := contexts.GetAndConvertContextFieldToString(ctx, types.ContextFieldName_AccessToken)
		csrfToken, exceptionOfCSRFToken := contexts.GetAndConvertContextFieldToString(ctx, types.ContextFieldName_CSRFToken)
		if (exceptionOfAccessToken != nil && exceptionOfCSRFToken != nil) || (accessToken == nil && csrfToken == nil) {
			return
		}

		cookies.AccessTokenCookieHandler.Set(ctx, *accessToken)
		originalResponse[types.AdditionalResponseFieldDomainName_RefreshableTokens.String()] = gin.H{
			types.RefreshableResponseFieldName_NewAccessToken.String(): *accessToken,
			types.RefreshableResponseFieldName_NewCSRFToken.String():   *csrfToken,
		}
		modifiedResponse, err := json.Marshal(originalResponse)
		if err != nil {
			return
		}

		writer.Mutex.Lock()
		writer.Body.Reset()
		writer.Body.Write(modifiedResponse)
		writer.Mutex.Unlock()
	}
}
