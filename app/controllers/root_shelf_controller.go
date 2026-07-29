package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RootShelfControllerInterface interface {
	GetMyRootShelfById(ctx *gin.Context, reqDto *dtos.GetMyRootShelfByIdReqDto)
	CreateRootShelf(ctx *gin.Context, reqDto *dtos.CreateRootShelfReqDto)
	CreateRootShelves(ctx *gin.Context, reqDto *dtos.CreateRootShelvesReqDto)
	UpdateMyRootShelfById(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto)
	UpdateMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelvesByIdsReqDto)
	RestoreMyRootShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto)
	RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto)
	DeleteMyRootShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto)
	DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto)

	GetMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.GetMyRootShelfPermissionReqDto)
	CreateMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.CreateMyRootShelfPermissionReqDto)
	UpsertMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto)
	UpsertMyRootShelfPermissions(ctx *gin.Context, reqDto *dtos.UpsertMyRootShelfPermissionsReqDto)
	UpdateMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfPermissionReqDto)
	TransferMyRootShelfOwnership(ctx *gin.Context, reqDto *dtos.TransferMyRootShelfOwnershipReqDto)
	DeleteMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto)
	DeleteMyRootShelfPermissions(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfPermissionsReqDto)
	LeaveMyRootShelf(ctx *gin.Context, reqDto *dtos.LeaveMyRootShelfReqDto)
	LeaveMyRootShelves(ctx *gin.Context, reqDto *dtos.LeaveMyRootShelvesReqDto)
}

type RootShelfController struct {
	rootShelfService services.RootShelfServiceInterface
}

func NewRootShelfController(service services.RootShelfServiceInterface) RootShelfControllerInterface {
	return &RootShelfController{
		rootShelfService: service,
	}
}

func (c *RootShelfController) GetMyRootShelfById(ctx *gin.Context, reqDto *dtos.GetMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.GetMyRootShelfById(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) CreateRootShelf(ctx *gin.Context, reqDto *dtos.CreateRootShelfReqDto) {
	resDto, exception := c.rootShelfService.CreateRootShelf(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) CreateRootShelves(ctx *gin.Context, reqDto *dtos.CreateRootShelvesReqDto) {
	resDto, exception := c.rootShelfService.CreateRootShelves(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) UpdateMyRootShelfById(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.UpdateMyRootShelfById(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) UpdateMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelvesByIdsReqDto) {
	resDto, exception := c.rootShelfService.UpdateMyRootShelvesByIds(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) RestoreMyRootShelfById(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.RestoreMyRootShelfById(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto) {
	resDto, exception := c.rootShelfService.RestoreMyRootShelvesByIds(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) DeleteMyRootShelfById(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto) {
	resDto, exception := c.rootShelfService.DeleteMyRootShelfById(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto) {
	resDto, exception := c.rootShelfService.DeleteMyRootShelvesByIds(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) GetMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.GetMyRootShelfPermissionReqDto) {
	resDto, exception := c.rootShelfService.GetMyRootShelfPermission(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) CreateMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.CreateMyRootShelfPermissionReqDto) {
	resDto, exception := c.rootShelfService.CreateMyRootShelfPermission(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) UpsertMyRootShelfPermission(
	ctx *gin.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto,
) {
	resDto, exception := c.rootShelfService.UpsertMyRootShelfPermission(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) UpsertMyRootShelfPermissions(
	ctx *gin.Context, reqDto *dtos.UpsertMyRootShelfPermissionsReqDto,
) {
	resDto, exception := c.rootShelfService.UpsertMyRootShelfPermissions(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) UpdateMyRootShelfPermission(ctx *gin.Context, reqDto *dtos.UpdateMyRootShelfPermissionReqDto) {
	resDto, exception := c.rootShelfService.UpdateMyRootShelfPermission(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) TransferMyRootShelfOwnership(
	ctx *gin.Context, reqDto *dtos.TransferMyRootShelfOwnershipReqDto,
) {
	resDto, exception := c.rootShelfService.TransferMyRootShelfOwnership(ctx.Request.Context(), reqDto)
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

func (c *RootShelfController) DeleteMyRootShelfPermission(
	ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto,
) {
	if exception := c.rootShelfService.DeleteMyRootShelfPermission(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) DeleteMyRootShelfPermissions(
	ctx *gin.Context, reqDto *dtos.DeleteMyRootShelfPermissionsReqDto,
) {
	if exception := c.rootShelfService.DeleteMyRootShelfPermissions(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelf(ctx *gin.Context, reqDto *dtos.LeaveMyRootShelfReqDto) {
	if exception := c.rootShelfService.LeaveMyRootShelf(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelves(ctx *gin.Context, reqDto *dtos.LeaveMyRootShelvesReqDto) {
	if exception := c.rootShelfService.LeaveMyRootShelves(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
