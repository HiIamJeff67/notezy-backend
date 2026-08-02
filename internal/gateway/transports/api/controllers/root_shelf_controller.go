package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/root-shelves"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RootShelfControllerInterface interface {
	GetMyRootShelfById(ctx *gin.Context, request *rootshelvesdto.GetMyRootShelfByIdRequestDto)
	CreateRootShelf(ctx *gin.Context, request *rootshelvesdto.CreateRootShelfRequestDto)
	CreateRootShelves(ctx *gin.Context, request *rootshelvesdto.CreateRootShelvesRequestDto)
	UpdateMyRootShelfById(ctx *gin.Context, request *rootshelvesdto.UpdateMyRootShelfByIdRequestDto)
	UpdateMyRootShelvesByIds(ctx *gin.Context, reqDto *rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto)
	RestoreMyRootShelfById(ctx *gin.Context, reqDto *rootshelvesdto.RestoreMyRootShelfByIdRequestDto)
	RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto)
	DeleteMyRootShelfById(ctx *gin.Context, reqDto *rootshelvesdto.DeleteMyRootShelfByIdRequestDto)
	DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto)

	GetMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.GetMyRootShelfPermissionRequestDto)
	CreateMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.CreateMyRootShelfPermissionRequestDto)
	UpsertMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionRequestDto)
	UpsertMyRootShelfPermissions(ctx *gin.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto)
	UpdateMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.UpdateMyRootShelfPermissionRequestDto)
	TransferMyRootShelfOwnership(ctx *gin.Context, requestDto *rootshelvesdto.TransferMyRootShelfOwnershipRequestDto)
	DeleteMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionRequestDto)
	DeleteMyRootShelfPermissions(ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto)
	LeaveMyRootShelf(ctx *gin.Context, requestDto *rootshelvesdto.LeaveMyRootShelfRequestDto)
	LeaveMyRootShelves(ctx *gin.Context, requestDto *rootshelvesdto.LeaveMyRootShelvesRequestDto)
}

type RootShelfController struct {
	coreClient *coreadapters.CoreClient
}

func NewRootShelfController(coreClient *coreadapters.CoreClient) RootShelfControllerInterface {
	return &RootShelfController{
		coreClient: coreClient,
	}
}

/* ============================== RootShelf Controller Methods ============================== */

func (c *RootShelfController) GetMyRootShelfById(ctx *gin.Context, request *rootshelvesdto.GetMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.GetMyRootShelfByIdRequestDto,
		rootshelvesdto.GetMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		rootshelvesdto.GetMyRootShelfByIdOperation,
		"/core/v1/root-shelves/get-by-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) CreateRootShelf(ctx *gin.Context, request *rootshelvesdto.CreateRootShelfRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.CreateRootShelfRequestDto,
		rootshelvesdto.CreateRootShelfResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		rootshelvesdto.CreateRootShelfOperation,
		"/core/v1/root-shelves/create",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) CreateRootShelves(ctx *gin.Context, request *rootshelvesdto.CreateRootShelvesRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.CreateRootShelvesRequestDto,
		rootshelvesdto.CreateRootShelvesResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		rootshelvesdto.CreateRootShelvesOperation,
		"/core/v1/root-shelves/create-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) UpdateMyRootShelfById(ctx *gin.Context, request *rootshelvesdto.UpdateMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.UpdateMyRootShelfByIdRequestDto,
		rootshelvesdto.UpdateMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		rootshelvesdto.UpdateMyRootShelfByIdOperation,
		"/core/v1/root-shelves/update",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) UpdateMyRootShelvesByIds(ctx *gin.Context, requestDto *rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto,
		rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.UpdateMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/update-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) RestoreMyRootShelfById(ctx *gin.Context, requestDto *rootshelvesdto.RestoreMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.RestoreMyRootShelfByIdRequestDto,
		rootshelvesdto.RestoreMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.RestoreMyRootShelfByIdOperation,
		"/core/v1/root-shelves/restore",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) RestoreMyRootShelvesByIds(ctx *gin.Context, requestDto *rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto,
		rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.RestoreMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/restore-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) DeleteMyRootShelfById(ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.DeleteMyRootShelfByIdRequestDto,
		rootshelvesdto.DeleteMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.DeleteMyRootShelfByIdOperation,
		"/core/v1/root-shelves/delete",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) DeleteMyRootShelvesByIds(ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto,
		rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.DeleteMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/delete-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) GetMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.GetMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[rootshelvesdto.GetMyRootShelfPermissionRequestDto, rootshelvesdto.GetMyRootShelfPermissionResponseDto](ctx, c.coreClient, requestDto, rootshelvesdto.GetMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/get")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) CreateMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.CreateMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[rootshelvesdto.CreateMyRootShelfPermissionRequestDto, rootshelvesdto.CreateMyRootShelfPermissionResponseDto](ctx, c.coreClient, requestDto, rootshelvesdto.CreateMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/create")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) UpsertMyRootShelfPermission(
	ctx *gin.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionRequestDto,
) {
	response, exception := coreadapters.CallSecurly[rootshelvesdto.UpsertMyRootShelfPermissionRequestDto, rootshelvesdto.UpsertMyRootShelfPermissionResponseDto](ctx, c.coreClient, requestDto, rootshelvesdto.UpsertMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/upsert")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) UpsertMyRootShelfPermissions(
	ctx *gin.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto,
		rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.UpsertMyRootShelfPermissionsOperation,
		"/core/v1/root-shelves/permissions/upsert-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) UpdateMyRootShelfPermission(ctx *gin.Context, requestDto *rootshelvesdto.UpdateMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[rootshelvesdto.UpdateMyRootShelfPermissionRequestDto, rootshelvesdto.UpdateMyRootShelfPermissionResponseDto](ctx, c.coreClient, requestDto, rootshelvesdto.UpdateMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/update")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) TransferMyRootShelfOwnership(
	ctx *gin.Context, requestDto *rootshelvesdto.TransferMyRootShelfOwnershipRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		rootshelvesdto.TransferMyRootShelfOwnershipRequestDto,
		rootshelvesdto.TransferMyRootShelfOwnershipResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.TransferMyRootShelfOwnershipOperation,
		"/core/v1/root-shelves/ownership/transfer",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RootShelfController) DeleteMyRootShelfPermission(
	ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		rootshelvesdto.DeleteMyRootShelfPermissionRequestDto,
		rootshelvesdto.DeleteMyRootShelfPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.DeleteMyRootShelfPermissionOperation,
		"/core/v1/root-shelves/permissions/delete",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) DeleteMyRootShelfPermissions(
	ctx *gin.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto,
		rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.DeleteMyRootShelfPermissionsOperation,
		"/core/v1/root-shelves/permissions/delete-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelf(ctx *gin.Context, requestDto *rootshelvesdto.LeaveMyRootShelfRequestDto) {
	_, exception := coreadapters.CallSecurly[
		rootshelvesdto.LeaveMyRootShelfRequestDto,
		rootshelvesdto.LeaveMyRootShelfResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.LeaveMyRootShelfOperation,
		"/core/v1/root-shelves/memberships/leave",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelves(ctx *gin.Context, requestDto *rootshelvesdto.LeaveMyRootShelvesRequestDto) {
	_, exception := coreadapters.CallSecurly[
		rootshelvesdto.LeaveMyRootShelvesRequestDto,
		rootshelvesdto.LeaveMyRootShelvesResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		rootshelvesdto.LeaveMyRootShelvesOperation,
		"/core/v1/root-shelves/memberships/leave-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
