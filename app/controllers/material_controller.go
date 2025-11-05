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
	GetMyMaterialAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto)
	GetAllMyMaterialsByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto)
	GetAllMyMaterialsByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto)
	CreateTextbookMaterial(ctx *gin.Context, reqDto *dtos.CreateTextbookMaterialReqDto)
	CreateNotebookMaterial(ctx *gin.Context, reqDto *dtos.CreateNotebookMaterialReqDto)
	UpdateMyMaterialById(ctx *gin.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto)
	SaveMyNotebookMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto)
	MoveMyMaterialById(ctx *gin.Context, reqDto *dtos.MoveMyMaterialByIdReqDto)
	MoveMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto)
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
func (c *MaterialController) GetMyMaterialAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto) {
	resDto, exception := c.materialService.GetMyMaterialAndItsParentById(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) GetAllMyMaterialsByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto) {
	resDto, exception := c.materialService.GetAllMyMaterialsByParentSubShelfId(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) GetAllMyMaterialsByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto) {
	resDto, exception := c.materialService.GetAllMyMaterialsByRootShelfId(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) CreateTextbookMaterial(ctx *gin.Context, reqDto *dtos.CreateTextbookMaterialReqDto) {
	resDto, exception := c.materialService.CreateTextbookMaterial(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) CreateNotebookMaterial(ctx *gin.Context, reqDto *dtos.CreateNotebookMaterialReqDto) {
	resDto, exception := c.materialService.CreateNotebookMaterial(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) UpdateMyMaterialById(ctx *gin.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.UpdateMyMaterialById(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.SaveMyTextbookMaterialById(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) SaveMyNotebookMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.SaveMyNotebookMaterialById(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.MoveMyMaterialById(ctx.Request.Context(), reqDto)
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
func (c *MaterialController) MoveMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.MoveMyMaterialsByIds(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.RestoreMyMaterialById(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.RestoreMyMaterialsByIds(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.DeleteMyMaterialById(ctx.Request.Context(), reqDto)
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
	resDto, exception := c.materialService.DeleteMyMaterialsByIds(ctx.Request.Context(), reqDto)
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
