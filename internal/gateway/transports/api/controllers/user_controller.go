package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	usersdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/users"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type UserControllerInterface interface {
	GetUserData(ctx *gin.Context, requestDto *usersdto.GetUserDataRequestDto)
	GetMe(ctx *gin.Context, requestDto *usersdto.GetMeRequestDto)
	UpdateMe(ctx *gin.Context, requestDto *usersdto.UpdateMeRequestDto)
}
type UserController struct {
	coreClient *coreadapters.CoreClient
}

func NewUserController(coreClient *coreadapters.CoreClient) UserControllerInterface {
	return &UserController{
		coreClient: coreClient,
	}
}

func (c *UserController) GetUserData(ctx *gin.Context, requestDto *usersdto.GetUserDataRequestDto) {
	response, exception := coreadapters.CallSecurly[usersdto.GetUserDataRequestDto, usersdto.GetUserDataResponseDto](
		ctx, c.coreClient, requestDto, usersdto.GetUserDataOperation, "/core/v1/users/data",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true, "data": response.Data, "exception": nil,
	})
}

func (c *UserController) GetMe(ctx *gin.Context, requestDto *usersdto.GetMeRequestDto) {
	response, exception := coreadapters.CallSecurly[usersdto.GetMeRequestDto, usersdto.GetMeResponseDto](
		ctx, c.coreClient, requestDto, usersdto.GetMeOperation, "/core/v1/users/me",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true, "data": response.Data, "exception": nil,
	})
}

func (c *UserController) UpdateMe(ctx *gin.Context, requestDto *usersdto.UpdateMeRequestDto) {
	response, exception := coreadapters.CallSecurly[usersdto.UpdateMeRequestDto, usersdto.UpdateMeResponseDto](
		ctx, c.coreClient, requestDto, usersdto.UpdateMeOperation, "/core/v1/users/me/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true, "data": response.Data, "exception": nil,
	})
}
