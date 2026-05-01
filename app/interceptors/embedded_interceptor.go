package interceptors

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	responsewriter "notezy-backend/shared/lib/responsewriter"
	types "notezy-backend/shared/types"
)

// To add additional field to the response with possibly embedded data that is required for the frontend.
// ex. the frontend require a publicId to indicate the user in their local database across APIs
func EmbeddedInterceptor(responseWriterKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var writer *responsewriter.ResponseWriter
		existingWriter, exist := ctx.Get(responseWriterKey)
		if !exist || existingWriter == nil {
			exceptions.Context.MissPlacingOrWrongInterceptorOrder(
				"Cannot find the existing response writer, " +
					"please make sure to call the ShareableResponseWriterInterceptor() and pass EmbeddedInterceptor() as one of the parameters",
			).Log()
			return
		}
		writer = existingWriter.(*responsewriter.ResponseWriter)

		ctx.Next()

		if writer.IsTimeout {
			return
		}

		if writer.ResponseWriter.Written() || writer.Status() >= 400 {
			return
		}

		if ctx.Writer.Status() >= 400 {
			return
		}

		var originalResponse map[string]interface{}
		if err := json.Unmarshal(writer.Body.Bytes(), &originalResponse); err != nil {
			return
		}

		publicId, exception := contexts.GetAndConvertContextFieldToString(ctx, types.ContextFieldName_User_PublicId)
		if exception != nil || publicId == nil {
			return
		}

		originalResponse[types.AdditionalResponseFieldDomainName_Embedded.String()] = gin.H{
			types.EmbeddedAuthorizedResponseFieldName_PublicId.String(): *publicId,
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
