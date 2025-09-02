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
	GetMyShelfById(ctx *gin.Context)
	GetRecentShelves(ctx *gin.Context)
	CreateShelf(ctx *gin.Context)
	SynchronizeShelves(ctx *gin.Context)
	RestoreMyShelf(ctx *gin.Context)
	RestoreMyShelves(ctx *gin.Context)
	DeleteMyShelf(ctx *gin.Context)
	DeleteMyShelves(ctx *gin.Context)
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

// with AuthMiddleware()
func (c *ShelfController) GetMyShelfById(ctx *gin.Context) {
	var reqDto dtos.GetMyShelfByIdReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.GetMyShelfById(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) GetRecentShelves(ctx *gin.Context) {
	var reqDto dtos.GetRecentShelvesReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
		exception.Log()
		exceptions.User.InvalidInput().WithError(err).ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.GetRecentShelves(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) CreateShelf(ctx *gin.Context) {
	var reqDto dtos.CreateShelfReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.CreateShelf(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
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
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.SynchronizeShelves(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) RestoreMyShelf(ctx *gin.Context) {
	var reqDto dtos.RestoreMyShelfReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.RestoreMyShelf(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) RestoreMyShelves(ctx *gin.Context) {
	var reqDto dtos.RestoreMyShelvesReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.RestoreMyShelves(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) DeleteMyShelf(ctx *gin.Context) {
	var reqDto dtos.DeleteMyShelfReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.DeleteMyShelf(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *ShelfController) DeleteMyShelves(ctx *gin.Context) {
	var reqDto dtos.DeleteMyShelvesReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.OwnerId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Shelf.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.shelfService.DeleteMyShelves(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
