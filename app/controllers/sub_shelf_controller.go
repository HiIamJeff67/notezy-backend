package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type SubShelfControllerInterface interface {
	GetMySubShelfById(ctx *gin.Context, reqDto *dtos.GetMySubShelfByIdReqDto)
	GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, reqDto *dtos.GetMySubShelvesByPrevSubShelfIdReqDto)
	GetAllMySubShelvesByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMySubShelvesByRootShelfIdReqDto)
	CreateSubShelfByRootShelfId(ctx *gin.Context, reqDto *dtos.CreateSubShelfByRootShelfIdReqDto)
	UpdateMySubShelfById(ctx *gin.Context, reqDto *dtos.UpdateMySubShelfByIdReqDto)
	MoveMySubShelf(ctx *gin.Context, reqDto *dtos.MoveMySubShelfReqDto)
	MoveMySubShelves(ctx *gin.Context, reqDto *dtos.MoveMySubShelvesReqDto)
	RestoreMySubShelfById(ctx *gin.Context, reqDto *dtos.RestoreMySubShelfByIdReqDto)
	RestoreMySubShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMySubShelvesByIdsReqDto)
	DeleteMySubShelfById(ctx *gin.Context, reqDto *dtos.DeleteMySubShelfByIdReqDto)
	DeleteMySubShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMySubShelvesByIdsReqDto)
}

type SubShelfController struct {
	subShelfService services.SubShelfServiceInterface
}

func NewSubShelfController(service services.SubShelfServiceInterface) SubShelfControllerInterface {
	return &SubShelfController{
		subShelfService: service,
	}
}

/* ============================== Implementations ============================== */

func (c *SubShelfController) GetMySubShelfById(ctx *gin.Context, reqDto *dtos.GetMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.GetMySubShelfById(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, reqDto *dtos.GetMySubShelvesByPrevSubShelfIdReqDto) {
	resDto, exception := c.subShelfService.GetMySubShelvesByPrevSubShelfId(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) GetAllMySubShelvesByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMySubShelvesByRootShelfIdReqDto) {
	resDto, exception := c.subShelfService.GetAllMySubShelvesByRootShelfId(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) CreateSubShelfByRootShelfId(ctx *gin.Context, reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) {
	resDto, exception := c.subShelfService.CreateSubShelfByRootShelfId(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) UpdateMySubShelfById(ctx *gin.Context, reqDto *dtos.UpdateMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.UpdateMySubShelfById(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) MoveMySubShelf(ctx *gin.Context, reqDto *dtos.MoveMySubShelfReqDto) {
	resDto, exception := c.subShelfService.MoveMySubShelf(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) MoveMySubShelves(ctx *gin.Context, reqDto *dtos.MoveMySubShelvesReqDto) {
	resDto, exception := c.subShelfService.MoveMySubShelves(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) DeleteMySubShelfById(ctx *gin.Context, reqDto *dtos.DeleteMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.DeleteMySubShelfById(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) RestoreMySubShelfById(ctx *gin.Context, reqDto *dtos.RestoreMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.RestoreMySubShelfById(ctx.Request.Context(), reqDto)
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
func (c *SubShelfController) RestoreMySubShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMySubShelvesByIdsReqDto) {
	resDto, exception := c.subShelfService.RestoreMySubShelvesByIds(ctx.Request.Context(), reqDto)
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

func (c *SubShelfController) DeleteMySubShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMySubShelvesByIdsReqDto) {
	resDto, exception := c.subShelfService.DeleteMySubShelvesByIds(ctx.Request.Context(), reqDto)
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
