package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type ShelfControllerInterface interface {
	GetMyShelfById(ctx *gin.Context, reqDto *dtos.GetMyShelfByIdReqDto)
	SearchRecentShelves(ctx *gin.Context, reqDto *dtos.SearchRecentShelvesReqDto)
	CreateShelf(ctx *gin.Context, reqDto *dtos.CreateShelfReqDto)
	SynchronizeShelves(ctx *gin.Context, reqDto *dtos.SynchronizeShelvesReqDto)
	RestoreMyShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyShelfByIdReqDto)
	RestoreMyShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyShelvesByIdsReqDto)
	DeleteMyShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyShelfByIdReqDto)
	DeleteMyShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyShelvesByIdsReqDto)
}

type ShelfController struct {
	shelfService services.ShelfServiceInterface
}

func NewShelfController(service services.ShelfServiceInterface) ShelfControllerInterface {
	return &ShelfController{
		shelfService: service,
	}
}

/* ============================== Controller ============================== */

// with AuthMiddleware()
func (c *ShelfController) GetMyShelfById(ctx *gin.Context, reqDto *dtos.GetMyShelfByIdReqDto) {
	resDto, exception := c.shelfService.GetMyShelfById(reqDto)
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
func (c *ShelfController) SearchRecentShelves(ctx *gin.Context, reqDto *dtos.SearchRecentShelvesReqDto) {
	resDto, exception := c.shelfService.SearchRecentShelves(reqDto)
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
func (c *ShelfController) CreateShelf(ctx *gin.Context, reqDto *dtos.CreateShelfReqDto) {
	resDto, exception := c.shelfService.CreateShelf(reqDto)
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
func (c *ShelfController) SynchronizeShelves(ctx *gin.Context, reqDto *dtos.SynchronizeShelvesReqDto) {
	resDto, exception := c.shelfService.SynchronizeShelves(reqDto)
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
func (c *ShelfController) RestoreMyShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyShelfByIdReqDto) {
	resDto, exception := c.shelfService.RestoreMyShelfById(reqDto)
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
func (c *ShelfController) RestoreMyShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyShelvesByIdsReqDto) {
	resDto, exception := c.shelfService.RestoreMyShelvesByIds(reqDto)
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
func (c *ShelfController) DeleteMyShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyShelfByIdReqDto) {
	resDto, exception := c.shelfService.DeleteMyShelfById(reqDto)
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
func (c *ShelfController) DeleteMyShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyShelvesByIdsReqDto) {
	resDto, exception := c.shelfService.DeleteMyShelvesByIds(reqDto)
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
