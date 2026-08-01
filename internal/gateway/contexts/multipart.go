package contexts

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func GetAndConvertContextToMultipartFileHeaders(ctx *gin.Context) ([]*multipart.FileHeader, *exceptions.Exception) {
	value, exists := ctx.Get(types.ContextFieldName_FormDataFileHeaders.String())
	if !exists {
		return nil, exceptions.New(
			"ContextFieldMissing",
			"Gateway",
			"ReadFormData",
			"The request context does not contain multipart file headers",
			http.StatusInternalServerError,
			true,
		)
	}

	fileHeaders, ok := value.([]*multipart.FileHeader)
	if !ok {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadFormData",
			"The multipart file headers have an invalid type",
			http.StatusInternalServerError,
			true,
		)
	}

	return fileHeaders, nil
}
