package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Instance ============================== */

type ShelfControllerInterface interface {
	CreateShelf(ctx *gin.Context)
	SynchronizeShelves(ctx *gin.Context)
}

type ShelfController struct {
	shelfService services.ShelfServiceInterface
}

func NewShelfController(service services.ShelfServiceInterface) ShelfControllerInterface {
	return &ShelfController{
		shelfService: service,
	}
}

/* ============================== Controllers ============================== */

// with AuthMiddleware
func (c *ShelfController) CreateShelf(ctx *gin.Context) {
	var reqDto dtos.CreateShelfReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log()
		exception = exceptions.Shelf.InternalServerWentWrong(nil)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.OwnerId = *userId

	resDto, exception := c.shelfService.CreateShelf(&reqDto)
	if exception != nil {
		exception.Log()
		if exception.IsInternal {
			exception = exceptions.Shelf.InternalServerWentWrong(nil)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) SynchronizeShelves(ctx *gin.Context) {
	var reqDto dtos.SynchronizeShelvesReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log()
		exception = exceptions.Shelf.InternalServerWentWrong(nil)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.OwnerId = *userId

	resDto, exception := c.shelfService.SynchronizeShelves(&reqDto)
	if exception != nil {
		exception.Log()
		if exception.IsInternal {
			exception = exceptions.Shelf.InternalServerWentWrong(nil)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
