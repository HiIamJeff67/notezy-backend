package contexts

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
)

func GetAndConvertContextFieldToUUID(ctx *gin.Context, name constants.ContextFieldName) (*uuid.UUID, *exceptions.Exception) {
	value, exist := ctx.Get(name.String())
	if !exist {
		return nil, exceptions.Context.FailedToGetContextFieldOfSpecificName(name.String())
	}

	valueString, ok := value.(string)
	if !ok {
		return nil, exceptions.Context.FailedToConvertContextFieldToSpecificType("string")
	}

	id, err := uuid.Parse(valueString)
	if err != nil {
		return nil, exceptions.Context.FailedToConvertContextFieldToSpecificType("uuid")
	}

	return &id, nil
}

func GetAndConvertContextToGinContext(ctx context.Context) (*gin.Context, *exceptions.Exception) {
	ginCtx, ok := ctx.Value(constants.ContextFieldName_GinContext).(*gin.Context)
	if !ok {
		return nil, exceptions.Context.FailedToConvertContextToGinContext()
	}
	return ginCtx, nil
}

func GetAndConvertContextToMultipartFileHeaders(ctx *gin.Context) ([]*multipart.FileHeader, *exceptions.Exception) {
	fileHeadersInterface, exist := ctx.Get(constants.ContextFieldName_FormDataFileHeaders.String())
	if exist {
		return fileHeadersInterface.([]*multipart.FileHeader), nil
	}
	return nil, exceptions.Context.FailedToConvertContextFieldToSpecificType("[]*multipart.FileHeader")
}
