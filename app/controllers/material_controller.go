package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type MaterialControllerInterface interface {
	GetMyMaterialById(ctx *gin.Context, reqDto *dtos.GetMyMaterialByIdReqDto)
	GetMyMaterialAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto)
	GetMyMaterialsByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetMyMaterialsByParentSubShelfIdReqDto)
	GetAllMyMaterialsByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto)
	CreateMyMaterial(ctx *gin.Context, reqDto *dtos.CreateMyMaterialReqDto)
	UpdateMyMaterialById(ctx *gin.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto)
	SaveMyMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto)
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

func (c *MaterialController) GetMyMaterialById(ctx *gin.Context, reqDto *dtos.GetMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.GetMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) GetMyMaterialAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto) {
	resDto, exception := c.materialService.GetMyMaterialAndItsParentById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) GetMyMaterialsByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetMyMaterialsByParentSubShelfIdReqDto) {
	resDto, exception := c.materialService.GetMyMaterialsByParentSubShelfId(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) GetAllMyMaterialsByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto) {
	resDto, exception := c.materialService.GetAllMyMaterialsByRootShelfId(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) CreateMyMaterial(ctx *gin.Context, reqDto *dtos.CreateMyMaterialReqDto) {
	resDto, exception := c.materialService.CreateMyMaterial(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) UpdateMyMaterialById(ctx *gin.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.UpdateMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) SaveMyMaterialById(ctx *gin.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.SaveMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) MoveMyMaterialById(ctx *gin.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.MoveMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) MoveMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.MoveMyMaterialsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) RestoreMyMaterialById(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.RestoreMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) RestoreMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.RestoreMyMaterialsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) DeleteMyMaterialById(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto) {
	resDto, exception := c.materialService.DeleteMyMaterialById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *MaterialController) DeleteMyMaterialsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto) {
	resDto, exception := c.materialService.DeleteMyMaterialsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
