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

type MaterialControllerInterface interface {
	GetMyMaterialById(ctx *gin.Context)
	SearchMyMaterialsByShelfId(ctx *gin.Context)
	CreateTextbookMaterial(ctx *gin.Context)
	RestoreMyMaterialById(ctx *gin.Context)
	RestoreMyMaterialsByIds(ctx *gin.Context)
	DeleteMyMaterialById(ctx *gin.Context)
	DeleteMyMaterialsByIds(ctx *gin.Context)
}

type MaterialController struct {
	materialService services.MaterialServiceInterface
}

func NewMaterialController(service services.MaterialServiceInterface) MaterialControllerInterface {
	return &MaterialController{
		materialService: service,
	}
}

/* ============================== Controller ============================== */

// with AuthMiddleware
func (c *MaterialController) GetMyMaterialById(ctx *gin.Context) {
	var reqDto dtos.GetMyMaterialByIdReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.GetMyMaterialById(&reqDto)
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

// with AuthMiddleware
func (c *MaterialController) SearchMyMaterialsByShelfId(ctx *gin.Context) {
	var reqDto dtos.SearchMyMaterialsByShelfIdReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
		exception.Log()
		exceptions.User.InvalidInput().WithError(err).ResponseWithJSON(ctx)
		return
	}
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.SearchMyMaterialsByShelfId(&reqDto)
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

// NOT yet done, have to prepare the storage first
// with AuthMiddleware
func (c *MaterialController) CreateTextbookMaterial(ctx *gin.Context) {
	var reqDto dtos.CreateMaterialReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}
}

// with AuthMiddleware
func (c *MaterialController) RestoreMyMaterialById(ctx *gin.Context) {
	var reqDto dtos.RestoreMyMaterialByIdReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.RestoreMyMaterialById(&reqDto)
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

// with AuthMiddleware
func (c *MaterialController) RestoreMyMaterialsByIds(ctx *gin.Context) {
	var reqDto dtos.RestoreMyMaterialsByIdsReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.RestoreMyMaterialsByIds(&reqDto)
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

// with AuthMiddleware
func (c *MaterialController) DeleteMyMaterialById(ctx *gin.Context) {
	var reqDto dtos.DeleteMyMaterialByIdReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.DeleteMyMaterialById(&reqDto)
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

// with AuthMiddleware
func (c *MaterialController) DeleteMyMaterialsByIds(ctx *gin.Context) {
	var reqDto dtos.DeleteMyMaterialsByIdsReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Material.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.materialService.DeleteMyMaterialsByIds(&reqDto)
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
