package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/materials"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type MaterialControllerInterface interface {
	GetMyMaterialById(ctx *gin.Context, requestDto *apicontract.GetMyMaterialByIdRequestDto)
	GetMyMaterialAndItsParentById(ctx *gin.Context, requestDto *apicontract.GetMyMaterialAndItsParentByIdRequestDto)
	GetMyMaterialsByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMyMaterialsByParentSubShelfIdRequestDto)
	GetAllMyMaterialsByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMyMaterialsByRootShelfIdRequestDto)
	CreateMyMaterial(ctx *gin.Context, requestDto *apicontract.CreateMyMaterialRequestDto)
	UpdateMyMaterialById(ctx *gin.Context, requestDto *apicontract.UpdateMyMaterialByIdRequestDto)
	SaveMyMaterialById(ctx *gin.Context, requestDto *apicontract.SaveMyMaterialByIdRequestDto)
	MoveMyMaterialById(ctx *gin.Context, requestDto *apicontract.MoveMyMaterialByIdRequestDto)
	MoveMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.MoveMyMaterialsByIdsRequestDto)
	RestoreMyMaterialById(ctx *gin.Context, requestDto *apicontract.RestoreMyMaterialByIdRequestDto)
	RestoreMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyMaterialsByIdsRequestDto)
	DeleteMyMaterialById(ctx *gin.Context, requestDto *apicontract.DeleteMyMaterialByIdRequestDto)
	DeleteMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyMaterialsByIdsRequestDto)
}

type MaterialController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewMaterialController(coreClient *coreadapters.CoreAdapter) MaterialControllerInterface {
	return &MaterialController{
		coreClient: coreClient,
	}
}

func (c *MaterialController) GetMyMaterialById(ctx *gin.Context, requestDto *apicontract.GetMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyMaterialByIdRequestDto,
		apicontract.GetMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyMaterialByIdOperation,
		"/core/v1/materials/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) GetMyMaterialAndItsParentById(ctx *gin.Context, requestDto *apicontract.GetMyMaterialAndItsParentByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyMaterialAndItsParentByIdRequestDto,
		apicontract.GetMyMaterialAndItsParentByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyMaterialAndItsParentByIdOperation,
		"/core/v1/materials/get-and-parent-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) GetMyMaterialsByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMyMaterialsByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyMaterialsByParentSubShelfIdRequestDto,
		apicontract.GetMyMaterialsByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyMaterialsByParentSubShelfIdOperation,
		"/core/v1/materials/get-by-parent-sub-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) GetAllMyMaterialsByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMyMaterialsByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetAllMyMaterialsByRootShelfIdRequestDto,
		apicontract.GetAllMyMaterialsByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetAllMyMaterialsByRootShelfIdOperation,
		"/core/v1/materials/get-all-by-root-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) CreateMyMaterial(ctx *gin.Context, requestDto *apicontract.CreateMyMaterialRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateMyMaterialRequestDto,
		apicontract.CreateMyMaterialResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.CreateMyMaterialOperation,
		"/core/v1/materials/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *MaterialController) UpdateMyMaterialById(ctx *gin.Context, requestDto *apicontract.UpdateMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyMaterialByIdRequestDto,
		apicontract.UpdateMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.UpdateMyMaterialByIdOperation,
		"/core/v1/materials/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) SaveMyMaterialById(ctx *gin.Context, requestDto *apicontract.SaveMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.SaveMyMaterialByIdRequestDto,
		apicontract.SaveMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.SaveMyMaterialByIdOperation,
		"/core/v1/materials/save",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) MoveMyMaterialById(ctx *gin.Context, requestDto *apicontract.MoveMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMyMaterialByIdRequestDto,
		apicontract.MoveMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.MoveMyMaterialByIdOperation,
		"/core/v1/materials/move",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) MoveMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.MoveMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMyMaterialsByIdsRequestDto,
		apicontract.MoveMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.MoveMyMaterialsByIdsOperation,
		"/core/v1/materials/move-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) RestoreMyMaterialById(ctx *gin.Context, requestDto *apicontract.RestoreMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyMaterialByIdRequestDto,
		apicontract.RestoreMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.RestoreMyMaterialByIdOperation,
		"/core/v1/materials/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) RestoreMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyMaterialsByIdsRequestDto,
		apicontract.RestoreMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.RestoreMyMaterialsByIdsOperation,
		"/core/v1/materials/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) DeleteMyMaterialById(ctx *gin.Context, requestDto *apicontract.DeleteMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyMaterialByIdRequestDto,
		apicontract.DeleteMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.DeleteMyMaterialByIdOperation,
		"/core/v1/materials/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *MaterialController) DeleteMyMaterialsByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyMaterialsByIdsRequestDto,
		apicontract.DeleteMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.DeleteMyMaterialsByIdsOperation,
		"/core/v1/materials/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
