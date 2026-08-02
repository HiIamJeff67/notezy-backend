package interceptors

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	cookies "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/cookies"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

// To add additional field to the response with adding additional field of `newAccessToken` and `newCSRFToken`,
// Note : It should be placed below the `AuthMiddleware`,
// so that it can access the `AccessToken` and `CSRFToken` in the context field
func RefreshTokenInterceptor(responseWriterKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var writer *responsewriter.ResponseWriter
		existingWriter, exist := ctx.Get(responseWriterKey)
		if !exist || existingWriter == nil {
			exceptions.New(
				"WrongInterceptorOrder",
				"Context",
				"Interceptor",
				"Cannot find the existing response writer, "+
					"please make sure to call the ShareableResponseWriterInterceptor() and pass RefreshTokenInterceptor() as one of the parameters",
				http.StatusInternalServerError,
				true,
			)
			return
		}
		writer = existingWriter.(*responsewriter.ResponseWriter)

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
		if exceptionOfAccessToken != nil || exceptionOfCSRFToken != nil || accessToken == nil || csrfToken == nil {
			return
		}

		cookies.AccessTokenCookieHandler.Set(ctx, *accessToken)
		ctx.Header("X-CSRF-Token", *csrfToken)
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
