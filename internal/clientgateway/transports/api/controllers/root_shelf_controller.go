package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/root-shelves"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type RootShelfControllerInterface interface {
	GetMyRootShelfById(ctx *gin.Context, request *apicontract.GetMyRootShelfByIdRequestDto)
	CreateRootShelf(ctx *gin.Context, request *apicontract.CreateRootShelfRequestDto)
	CreateRootShelves(ctx *gin.Context, request *apicontract.CreateRootShelvesRequestDto)
	UpdateMyRootShelfById(ctx *gin.Context, request *apicontract.UpdateMyRootShelfByIdRequestDto)
	UpdateMyRootShelvesByIds(ctx *gin.Context, reqDto *apicontract.UpdateMyRootShelvesByIdsRequestDto)
	RestoreMyRootShelfById(ctx *gin.Context, reqDto *apicontract.RestoreMyRootShelfByIdRequestDto)
	RestoreMyRootShelvesByIds(ctx *gin.Context, reqDto *apicontract.RestoreMyRootShelvesByIdsRequestDto)
	DeleteMyRootShelfById(ctx *gin.Context, reqDto *apicontract.DeleteMyRootShelfByIdRequestDto)
	DeleteMyRootShelvesByIds(ctx *gin.Context, reqDto *apicontract.DeleteMyRootShelvesByIdsRequestDto)

	GetMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.GetMyRootShelfPermissionRequestDto)
	CreateMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.CreateMyRootShelfPermissionRequestDto)
	UpsertMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.UpsertMyRootShelfPermissionRequestDto)
	UpsertMyRootShelfPermissions(ctx *gin.Context, requestDto *apicontract.UpsertMyRootShelfPermissionsRequestDto)
	UpdateMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.UpdateMyRootShelfPermissionRequestDto)
	TransferMyRootShelfOwnership(ctx *gin.Context, requestDto *apicontract.TransferMyRootShelfOwnershipRequestDto)
	DeleteMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelfPermissionRequestDto)
	DeleteMyRootShelfPermissions(ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelfPermissionsRequestDto)
	LeaveMyRootShelf(ctx *gin.Context, requestDto *apicontract.LeaveMyRootShelfRequestDto)
	LeaveMyRootShelves(ctx *gin.Context, requestDto *apicontract.LeaveMyRootShelvesRequestDto)
}

type RootShelfController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewRootShelfController(coreAdapter *coreadapters.CoreAdapter) RootShelfControllerInterface {
	return &RootShelfController{
		coreAdapter: coreAdapter,
	}
}

/* ============================== RootShelf Controller Methods ============================== */

func (c *RootShelfController) GetMyRootShelfById(ctx *gin.Context, request *apicontract.GetMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyRootShelfByIdRequestDto,
		apicontract.GetMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.GetMyRootShelfByIdOperation,
		"/core/v1/root-shelves/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) CreateRootShelf(ctx *gin.Context, request *apicontract.CreateRootShelfRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateRootShelfRequestDto,
		apicontract.CreateRootShelfResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.CreateRootShelfOperation,
		"/core/v1/root-shelves/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RootShelfController) CreateRootShelves(ctx *gin.Context, request *apicontract.CreateRootShelvesRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateRootShelvesRequestDto,
		apicontract.CreateRootShelvesResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.CreateRootShelvesOperation,
		"/core/v1/root-shelves/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RootShelfController) UpdateMyRootShelfById(ctx *gin.Context, request *apicontract.UpdateMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyRootShelfByIdRequestDto,
		apicontract.UpdateMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.UpdateMyRootShelfByIdOperation,
		"/core/v1/root-shelves/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) UpdateMyRootShelvesByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyRootShelvesByIdsRequestDto,
		apicontract.UpdateMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) RestoreMyRootShelfById(ctx *gin.Context, requestDto *apicontract.RestoreMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyRootShelfByIdRequestDto,
		apicontract.RestoreMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyRootShelfByIdOperation,
		"/core/v1/root-shelves/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) RestoreMyRootShelvesByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyRootShelvesByIdsRequestDto,
		apicontract.RestoreMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) DeleteMyRootShelfById(ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyRootShelfByIdRequestDto,
		apicontract.DeleteMyRootShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRootShelfByIdOperation,
		"/core/v1/root-shelves/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) DeleteMyRootShelvesByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyRootShelvesByIdsRequestDto,
		apicontract.DeleteMyRootShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRootShelvesByIdsOperation,
		"/core/v1/root-shelves/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) GetMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.GetMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetMyRootShelfPermissionRequestDto, apicontract.GetMyRootShelfPermissionResponseDto](ctx, c.coreAdapter, requestDto, apicontract.GetMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/get")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) CreateMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.CreateMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.CreateMyRootShelfPermissionRequestDto, apicontract.CreateMyRootShelfPermissionResponseDto](ctx, c.coreAdapter, requestDto, apicontract.CreateMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/create")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RootShelfController) UpsertMyRootShelfPermission(
	ctx *gin.Context, requestDto *apicontract.UpsertMyRootShelfPermissionRequestDto,
) {
	response, exception := coreadapters.CallSecurly[apicontract.UpsertMyRootShelfPermissionRequestDto, apicontract.UpsertMyRootShelfPermissionResponseDto](ctx, c.coreAdapter, requestDto, apicontract.UpsertMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/upsert")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) UpsertMyRootShelfPermissions(
	ctx *gin.Context, requestDto *apicontract.UpsertMyRootShelfPermissionsRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpsertMyRootShelfPermissionsRequestDto,
		apicontract.UpsertMyRootShelfPermissionsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpsertMyRootShelfPermissionsOperation,
		"/core/v1/root-shelves/permissions/upsert-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) UpdateMyRootShelfPermission(ctx *gin.Context, requestDto *apicontract.UpdateMyRootShelfPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.UpdateMyRootShelfPermissionRequestDto, apicontract.UpdateMyRootShelfPermissionResponseDto](ctx, c.coreAdapter, requestDto, apicontract.UpdateMyRootShelfPermissionOperation, "/core/v1/root-shelves/permissions/update")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) TransferMyRootShelfOwnership(
	ctx *gin.Context, requestDto *apicontract.TransferMyRootShelfOwnershipRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.TransferMyRootShelfOwnershipRequestDto,
		apicontract.TransferMyRootShelfOwnershipResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.TransferMyRootShelfOwnershipOperation,
		"/core/v1/root-shelves/ownership/transfer",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RootShelfController) DeleteMyRootShelfPermission(
	ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelfPermissionRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyRootShelfPermissionRequestDto,
		apicontract.DeleteMyRootShelfPermissionResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRootShelfPermissionOperation,
		"/core/v1/root-shelves/permissions/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) DeleteMyRootShelfPermissions(
	ctx *gin.Context, requestDto *apicontract.DeleteMyRootShelfPermissionsRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyRootShelfPermissionsRequestDto,
		apicontract.DeleteMyRootShelfPermissionsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRootShelfPermissionsOperation,
		"/core/v1/root-shelves/permissions/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelf(ctx *gin.Context, requestDto *apicontract.LeaveMyRootShelfRequestDto) {
	_, exception := coreadapters.CallSecurly[
		apicontract.LeaveMyRootShelfRequestDto,
		apicontract.LeaveMyRootShelfResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LeaveMyRootShelfOperation,
		"/core/v1/root-shelves/memberships/leave",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *RootShelfController) LeaveMyRootShelves(ctx *gin.Context, requestDto *apicontract.LeaveMyRootShelvesRequestDto) {
	_, exception := coreadapters.CallSecurly[
		apicontract.LeaveMyRootShelvesRequestDto,
		apicontract.LeaveMyRootShelvesResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LeaveMyRootShelvesOperation,
		"/core/v1/root-shelves/memberships/leave-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
