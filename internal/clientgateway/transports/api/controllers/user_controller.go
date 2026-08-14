package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/users"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type UserControllerInterface interface {
	GetUserData(ctx *gin.Context, requestDto *apicontract.GetUserDataRequestDto)
	GetMe(ctx *gin.Context, requestDto *apicontract.GetMeRequestDto)
	UpdateMe(ctx *gin.Context, requestDto *apicontract.UpdateMeRequestDto)
}
type UserController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewUserController(coreClient *coreadapters.CoreAdapter) UserControllerInterface {
	return &UserController{
		coreClient: coreClient,
	}
}

func (c *UserController) GetUserData(ctx *gin.Context, requestDto *apicontract.GetUserDataRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetUserDataRequestDto, apicontract.GetUserDataResponseDto](
		ctx, c.coreClient, requestDto, apicontract.GetUserDataOperation, "/core/v1/users/data",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *UserController) GetMe(ctx *gin.Context, requestDto *apicontract.GetMeRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetMeRequestDto, apicontract.GetMeResponseDto](
		ctx, c.coreClient, requestDto, apicontract.GetMeOperation, "/core/v1/users/me",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *UserController) UpdateMe(ctx *gin.Context, requestDto *apicontract.UpdateMeRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.UpdateMeRequestDto, apicontract.UpdateMeResponseDto](
		ctx, c.coreClient, requestDto, apicontract.UpdateMeOperation, "/core/v1/users/me/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}
