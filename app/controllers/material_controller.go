package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type MaterialControllerInterface interface {
	GetMyMaterialById(ctx *gin.Context, reqDto *dtos.GetMyMaterialByIdReqDto)
	SearchMyMaterialsByShelfId(ctx *gin.Context, reqDto *dtos.SearchMyMaterialsByShelfIdReqDto)
	CreateTextbookMaterial(ctx *gin.Context, reqDto *dtos.CreateMaterialReqDto)
	SaveMyTextbookMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto)
	MoveMyMaterialById(ctx *gin.Context, reqDto *dtos.MoveMyMaterialByIdReqDto)
	RestoreMyMaterialById(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto)
	RestoreMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto)
	DeleteMyMaterialById(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto)
	DeleteMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto)
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
func (c *MaterialController) GetMyMaterialById(ctx *gin.Context, reqDto *dtos.GetMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.GetMyMaterialById(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) SearchMyMaterialsByShelfId(ctx *gin.Context, reqDto *dtos.SearchMyMaterialsByShelfIdReqDto) {
	resDto, exception := c.materialService.SearchMyMaterialsByShelfId(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) CreateTextbookMaterial(ctx *gin.Context, reqDto *dtos.CreateMaterialReqDto) {
	resDto, exception := c.materialService.CreateTextbookMaterial(ctx, reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": false,
	})
}

// with AuthMiddleware, MultipartAdapter
func (c *MaterialController) SaveMyTextbookMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.SaveMyTextbookMaterialById(ctx, reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": false,
	})
}

// with AuthMiddleware, MultipartAdapter
func (c *MaterialController) MoveMyMaterialById(ctx *gin.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.MoveMyMaterialById(ctx, reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": false,
	})
}

// with AuthMiddleware
func (c *MaterialController) RestoreMyMaterialById(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.RestoreMyMaterialById(reqDto)
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
func (c *MaterialController) RestoreMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.RestoreMyMaterialsByIds(reqDto)
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
func (c *MaterialController) DeleteMyMaterialById(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.DeleteMyMaterialById(reqDto)
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
func (c *MaterialController) DeleteMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.DeleteMyMaterialsByIds(reqDto)
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
