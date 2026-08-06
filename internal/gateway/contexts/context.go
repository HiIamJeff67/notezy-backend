package contexts

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
)

func GetAndConvertContextFieldToBoolean(ctx *gin.Context, name sharedcontexts.ContextFieldName) (*bool, *exceptions.Exception) {
	value, exists := ctx.Get(name.String())
	if !exists {
		return nil, exceptions.New(
			"ContextFieldMissing",
			"Gateway",
			"ReadContextField",
			"The required request context field is missing",
			http.StatusInternalServerError,
			true,
		)
	}

	valueBoolean, ok := value.(bool)
	if !ok {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadContextField",
			"The request context field has an invalid type",
			http.StatusInternalServerError,
			true,
		)
	}

	return &valueBoolean, nil
}

func GetAndConvertContextFieldToString(ctx *gin.Context, name sharedcontexts.ContextFieldName) (*string, *exceptions.Exception) {
	value, exists := ctx.Get(name.String())
	if !exists {
		return nil, exceptions.New(
			"ContextFieldMissing",
			"Gateway",
			"ReadContextField",
			"The required request context field is missing",
			http.StatusInternalServerError,
			true,
		)
	}

	valueString, ok := value.(string)
	if !ok {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadContextField",
			"The request context field has an invalid type",
			http.StatusInternalServerError,
			true,
		)
	}

	return &valueString, nil
}

func GetAndConvertContextFieldToUUID(ctx *gin.Context, name sharedcontexts.ContextFieldName) (*uuid.UUID, *exceptions.Exception) {
	value, exists := ctx.Get(name.String())
	if !exists {
		return nil, exceptions.New(
			"ContextFieldMissing",
			"Gateway",
			"ReadContextField",
			"The required request context field is missing",
			http.StatusInternalServerError,
			true,
		)
	}

	if valueUUID, ok := value.(uuid.UUID); ok {
		return &valueUUID, nil
	}

	valueString, ok := value.(string)
	if !ok {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadContextField",
			"The request context field has an invalid type",
			http.StatusInternalServerError,
			true,
		)
	}

	id, err := uuid.Parse(valueString)
	if err != nil {
		return nil, exceptions.New(
			"InvalidInput",
			"Gateway",
			"ReadContextField",
			"The request context field contains an invalid UUID",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	return &id, nil
}

func GetAndConvertContextToGinContext(ctx context.Context) (*gin.Context, *exceptions.Exception) {
	ginCtx, ok := ctx.Value(sharedcontexts.ContextFieldName_GinContext).(*gin.Context)
	if !ok {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadGinContext",
			"The request context does not contain a Gin context",
			http.StatusInternalServerError,
			true,
		)
	}

	return ginCtx, nil
}
