package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	materialsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/materials"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type MaterialControllerInterface interface {
	GetMyMaterialById(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialByIdRequestDto)
	GetMyMaterialAndItsParentById(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialAndItsParentByIdRequestDto)
	GetMyMaterialsByParentSubShelfId(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto)
	GetAllMyMaterialsByRootShelfId(ctx *gin.Context, requestDto *materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto)
	CreateMyMaterial(ctx *gin.Context, requestDto *materialsdto.CreateMyMaterialRequestDto)
	UpdateMyMaterialById(ctx *gin.Context, requestDto *materialsdto.UpdateMyMaterialByIdRequestDto)
	SaveMyMaterialById(ctx *gin.Context, requestDto *materialsdto.SaveMyMaterialByIdRequestDto)
	MoveMyMaterialById(ctx *gin.Context, requestDto *materialsdto.MoveMyMaterialByIdRequestDto)
	MoveMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.MoveMyMaterialsByIdsRequestDto)
	RestoreMyMaterialById(ctx *gin.Context, requestDto *materialsdto.RestoreMyMaterialByIdRequestDto)
	RestoreMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.RestoreMyMaterialsByIdsRequestDto)
	DeleteMyMaterialById(ctx *gin.Context, requestDto *materialsdto.DeleteMyMaterialByIdRequestDto)
	DeleteMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.DeleteMyMaterialsByIdsRequestDto)
}

type MaterialController struct {
	coreClient *coreadapters.CoreClient
}

func NewMaterialController(coreClient *coreadapters.CoreClient) MaterialControllerInterface {
	return &MaterialController{
		coreClient: coreClient,
	}
}

func (c *MaterialController) GetMyMaterialById(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.GetMyMaterialByIdRequestDto,
		materialsdto.GetMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.GetMyMaterialByIdOperation,
		"/core/v1/materials/get-by-id",
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

func (c *MaterialController) GetMyMaterialAndItsParentById(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialAndItsParentByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.GetMyMaterialAndItsParentByIdRequestDto,
		materialsdto.GetMyMaterialAndItsParentByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.GetMyMaterialAndItsParentByIdOperation,
		"/core/v1/materials/get-and-parent-by-id",
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

func (c *MaterialController) GetMyMaterialsByParentSubShelfId(ctx *gin.Context, requestDto *materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto,
		materialsdto.GetMyMaterialsByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.GetMyMaterialsByParentSubShelfIdOperation,
		"/core/v1/materials/get-by-parent-sub-shelf-id",
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

func (c *MaterialController) GetAllMyMaterialsByRootShelfId(ctx *gin.Context, requestDto *materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto,
		materialsdto.GetAllMyMaterialsByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.GetAllMyMaterialsByRootShelfIdOperation,
		"/core/v1/materials/get-all-by-root-shelf-id",
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

func (c *MaterialController) CreateMyMaterial(ctx *gin.Context, requestDto *materialsdto.CreateMyMaterialRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.CreateMyMaterialRequestDto,
		materialsdto.CreateMyMaterialResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.CreateMyMaterialOperation,
		"/core/v1/materials/create",
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

func (c *MaterialController) UpdateMyMaterialById(ctx *gin.Context, requestDto *materialsdto.UpdateMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.UpdateMyMaterialByIdRequestDto,
		materialsdto.UpdateMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.UpdateMyMaterialByIdOperation,
		"/core/v1/materials/update",
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

func (c *MaterialController) SaveMyMaterialById(ctx *gin.Context, requestDto *materialsdto.SaveMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.SaveMyMaterialByIdRequestDto,
		materialsdto.SaveMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.SaveMyMaterialByIdOperation,
		"/core/v1/materials/save",
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

func (c *MaterialController) MoveMyMaterialById(ctx *gin.Context, requestDto *materialsdto.MoveMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.MoveMyMaterialByIdRequestDto,
		materialsdto.MoveMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.MoveMyMaterialByIdOperation,
		"/core/v1/materials/move",
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

func (c *MaterialController) MoveMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.MoveMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.MoveMyMaterialsByIdsRequestDto,
		materialsdto.MoveMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.MoveMyMaterialsByIdsOperation,
		"/core/v1/materials/move-many",
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

func (c *MaterialController) RestoreMyMaterialById(ctx *gin.Context, requestDto *materialsdto.RestoreMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.RestoreMyMaterialByIdRequestDto,
		materialsdto.RestoreMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.RestoreMyMaterialByIdOperation,
		"/core/v1/materials/restore",
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

func (c *MaterialController) RestoreMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.RestoreMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.RestoreMyMaterialsByIdsRequestDto,
		materialsdto.RestoreMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.RestoreMyMaterialsByIdsOperation,
		"/core/v1/materials/restore-many",
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

func (c *MaterialController) DeleteMyMaterialById(ctx *gin.Context, requestDto *materialsdto.DeleteMyMaterialByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.DeleteMyMaterialByIdRequestDto,
		materialsdto.DeleteMyMaterialByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.DeleteMyMaterialByIdOperation,
		"/core/v1/materials/delete",
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

func (c *MaterialController) DeleteMyMaterialsByIds(ctx *gin.Context, requestDto *materialsdto.DeleteMyMaterialsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		materialsdto.DeleteMyMaterialsByIdsRequestDto,
		materialsdto.DeleteMyMaterialsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		materialsdto.DeleteMyMaterialsByIdsOperation,
		"/core/v1/materials/delete-many",
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
