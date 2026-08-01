package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	subshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/sub-shelves"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type SubShelfControllerInterface interface {
	GetMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelfByIdRequestDto)
	GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto)
	GetAllMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto)
	GetMySubShelvesAndItemsByPrevSubShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto)
	CreateSubShelfByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.CreateSubShelfByRootShelfIdRequestDto)
	CreateSubShelvesByRootShelfIds(ctx *gin.Context, requestDto *subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto)
	UpdateMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.UpdateMySubShelfByIdRequestDto)
	UpdateMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.UpdateMySubShelvesByIdsRequestDto)
	MoveMySubShelfByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto)
	MoveMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto)
	MoveMySubShelvesByRootShelfIds(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto)
	RestoreMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.RestoreMySubShelfByIdRequestDto)
	RestoreMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.RestoreMySubShelvesByIdsRequestDto)
	DeleteMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.DeleteMySubShelfByIdRequestDto)
	DeleteMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.DeleteMySubShelvesByIdsRequestDto)
}

type SubShelfController struct {
	coreClient *coreadapters.CoreClient
}

func NewSubShelfController(coreClient *coreadapters.CoreClient) SubShelfControllerInterface {
	return &SubShelfController{
		coreClient: coreClient,
	}
}

func (c *SubShelfController) GetMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.GetMySubShelfByIdRequestDto,
		subshelvesdto.GetMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.GetMySubShelfByIdOperation,
		"/core/v1/sub-shelves/get-by-id",
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

func (c *SubShelfController) GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto,
		subshelvesdto.GetMySubShelvesByPrevSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.GetMySubShelvesByPrevSubShelfIdOperation,
		"/core/v1/sub-shelves/get-by-prev-sub-shelf-id",
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

func (c *SubShelfController) GetAllMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto,
		subshelvesdto.GetAllMySubShelvesByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.GetAllMySubShelvesByRootShelfIdOperation,
		"/core/v1/sub-shelves/get-all-by-root-shelf-id",
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

func (c *SubShelfController) GetMySubShelvesAndItemsByPrevSubShelfId(ctx *gin.Context, requestDto *subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto,
		subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdOperation,
		"/core/v1/sub-shelves/get-and-items-by-prev-sub-shelf-id",
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

func (c *SubShelfController) CreateSubShelfByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.CreateSubShelfByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.CreateSubShelfByRootShelfIdRequestDto,
		subshelvesdto.CreateSubShelfByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.CreateSubShelfByRootShelfIdOperation,
		"/core/v1/sub-shelves/create",
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

func (c *SubShelfController) CreateSubShelvesByRootShelfIds(ctx *gin.Context, requestDto *subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto,
		subshelvesdto.CreateSubShelvesByRootShelfIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.CreateSubShelvesByRootShelfIdsOperation,
		"/core/v1/sub-shelves/create-many",
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

func (c *SubShelfController) UpdateMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.UpdateMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.UpdateMySubShelfByIdRequestDto,
		subshelvesdto.UpdateMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.UpdateMySubShelfByIdOperation,
		"/core/v1/sub-shelves/update",
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

func (c *SubShelfController) UpdateMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.UpdateMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.UpdateMySubShelvesByIdsRequestDto,
		subshelvesdto.UpdateMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.UpdateMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/update-many",
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

func (c *SubShelfController) MoveMySubShelfByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto,
		subshelvesdto.MoveMySubShelfByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.MoveMySubShelfByRootShelfIdOperation,
		"/core/v1/sub-shelves/move",
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

func (c *SubShelfController) MoveMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto,
		subshelvesdto.MoveMySubShelvesByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.MoveMySubShelvesByRootShelfIdOperation,
		"/core/v1/sub-shelves/move-many",
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

func (c *SubShelfController) MoveMySubShelvesByRootShelfIds(ctx *gin.Context, requestDto *subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto,
		subshelvesdto.MoveMySubShelvesByRootShelfIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.MoveMySubShelvesByRootShelfIdsOperation,
		"/core/v1/sub-shelves/move-many-by-root-shelves",
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

func (c *SubShelfController) RestoreMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.RestoreMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.RestoreMySubShelfByIdRequestDto,
		subshelvesdto.RestoreMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.RestoreMySubShelfByIdOperation,
		"/core/v1/sub-shelves/restore",
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

func (c *SubShelfController) RestoreMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.RestoreMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.RestoreMySubShelvesByIdsRequestDto,
		subshelvesdto.RestoreMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.RestoreMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/restore-many",
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

func (c *SubShelfController) DeleteMySubShelfById(ctx *gin.Context, requestDto *subshelvesdto.DeleteMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.DeleteMySubShelfByIdRequestDto,
		subshelvesdto.DeleteMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.DeleteMySubShelfByIdOperation,
		"/core/v1/sub-shelves/delete",
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

func (c *SubShelfController) DeleteMySubShelvesByIds(ctx *gin.Context, requestDto *subshelvesdto.DeleteMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		subshelvesdto.DeleteMySubShelvesByIdsRequestDto,
		subshelvesdto.DeleteMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		subshelvesdto.DeleteMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/delete-many",
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
