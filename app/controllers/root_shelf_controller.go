package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type RootShelfControllerInterface interface {
	GetMyRootShelfById(ctx *gin.Context, reqDto *dtos.GetMyRootShelfByIdReqDto)
	SearchRecentRootShelves(ctx *gin.Context, reqDto *dtos.SearchRecentRootShelvesReqDto)
	CreateRootShelf(ctx *gin.Context, reqDto *dtos.CreateRootShelfReqDto)
	UpdateMyRootShelfById(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto)
	RestoreMyRootShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto)
	RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto)
	DeleteMyRootShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto)
	DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto)
}

type RootShelfController struct {
	rootShelfService services.RootShelfServiceInterface
}

func NewRootShelfController(service services.RootShelfServiceInterface) RootShelfControllerInterface {
	return &RootShelfController{
		rootShelfService: service,
	}
}

/* ============================== Controller ============================== */

// with AuthMiddleware()
func (c *RootShelfController) GetMyRootShelfById(ctx *gin.Context, reqDto *dtos.GetMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.GetMyRootShelfById(reqDto)
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
func (c *RootShelfController) SearchRecentRootShelves(ctx *gin.Context, reqDto *dtos.SearchRecentRootShelvesReqDto) {
	resDto, exception := c.rootShelfService.SearchRecentRootShelves(reqDto)
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
func (c *RootShelfController) CreateRootShelf(ctx *gin.Context, reqDto *dtos.CreateRootShelfReqDto) {
	resDto, exception := c.rootShelfService.CreateRootShelf(reqDto)
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
func (c *RootShelfController) UpdateMyRootShelfById(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.UpdateMyRootShelfById(reqDto)
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
func (c *RootShelfController) RestoreMyRootShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.RestoreMyRootShelfById(reqDto)
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
func (c *RootShelfController) RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto) {
	resDto, exception := c.rootShelfService.RestoreMyRootShelvesByIds(reqDto)
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
func (c *RootShelfController) DeleteMyRootShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.DeleteMyRootShelfById(reqDto)
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
func (c *RootShelfController) DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto) {
	resDto, exception := c.rootShelfService.DeleteMyRootShelvesByIds(reqDto)
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
