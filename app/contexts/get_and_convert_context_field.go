package contexts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
)

func GetAndConvertContextFieldToUUID(ctx *gin.Context, name string) (*uuid.UUID, *exceptions.Exception) {
	value, exist := ctx.Get(name)
	if !exist {
		return nil, exceptions.Context.FailedToGetContextFieldOfSpecificName(name)
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
