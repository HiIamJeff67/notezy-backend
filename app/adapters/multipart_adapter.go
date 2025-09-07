package adapters

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
)

func MultipartAdapter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			exceptions.Adapter.
				InvalidMultipartForm().
				Log().SafelyResponseWithJSON(ctx)
			ctx.Abort()
			return
		}

		jsonData := make(map[string]interface{})
		var fileHeaders []*multipart.FileHeader

		for key, values := range form.Value {
			if len(values) > 0 {
				valueStr := values[0]
				if intVal, err := strconv.Atoi(valueStr); err == nil {
					jsonData[key] = intVal
				} else if boolVal, err := strconv.ParseBool(valueStr); err == nil {
					jsonData[key] = boolVal
				} else {
					jsonData[key] = valueStr
				}
			}
		}

		for _, fileHeadersSlice := range form.File {
			for _, fileHeader := range fileHeadersSlice {
				if fileHeader.Size > constants.MaxNonVideoFileSize {
					exceptions.Adapter.
						FileTooLarge(fileHeader.Size, constants.MaxNonVideoFileSize).
						Log().SafelyResponseWithJSON(ctx)
					ctx.Abort()
					return
				}
				fileHeaders = append(fileHeaders, fileHeader)
			}
		}

		if len(jsonData) > 0 {
			jsonBytes, _ := json.Marshal(jsonData)
			ctx.Request.Body = io.NopCloser(bytes.NewReader(jsonBytes))
		}

		if len(fileHeaders) > 0 {
			ctx.Set(constants.ContextFieldName_FormDataFileHeaders.String(), fileHeaders)
		}

		ctx.Next()
	}
}
