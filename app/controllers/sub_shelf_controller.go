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
	GetAllSubShelvesByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllSubShelvesByRootShelfIdReqDto)
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

/* ============================== Controller ============================== */

// with AuthMiddleware()
func (c *SubShelfController) GetMySubShelfById(ctx *gin.Context, reqDto *dtos.GetMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.GetMySubShelfById(reqDto)
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
func (c *SubShelfController) GetAllSubShelvesByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllSubShelvesByRootShelfIdReqDto) {
	resDto, exception := c.subShelfService.GetAllSubShelvesByRootShelfId(reqDto)
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
func (c *SubShelfController) CreateSubShelfByRootShelfId(ctx *gin.Context, reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) {
	resDto, exception := c.subShelfService.CreateSubShelfByRootShelfId(reqDto)
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
func (c *SubShelfController) UpdateMySubShelfById(ctx *gin.Context, reqDto *dtos.UpdateMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.UpdateMySubShelfById(reqDto)
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
func (c *SubShelfController) MoveMySubShelf(ctx *gin.Context, reqDto *dtos.MoveMySubShelfReqDto) {
	resDto, exception := c.subShelfService.MoveMySubShelf(reqDto)
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
func (c *SubShelfController) MoveMySubShelves(ctx *gin.Context, reqDto *dtos.MoveMySubShelvesReqDto) {
	resDto, exception := c.subShelfService.MoveMySubShelves(reqDto)
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
func (c *SubShelfController) DeleteMySubShelfById(ctx *gin.Context, reqDto *dtos.DeleteMySubShelfByIdReqDto) {
	resDto, exception := c.subShelfService.DeleteMySubShelfById(reqDto)
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
	resDto, exception := c.subShelfService.RestoreMySubShelfById(reqDto)
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
	resDto, exception := c.subShelfService.RestoreMySubShelvesByIds(reqDto)
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
func (c *SubShelfController) DeleteMySubShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMySubShelvesByIdsReqDto) {
	resDto, exception := c.subShelfService.DeleteMySubShelvesByIds(reqDto)
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
