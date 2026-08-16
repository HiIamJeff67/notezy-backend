package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/user-infos"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context, request *apicontract.GetMyInfoRequestDto)
	UpdateMyInfo(ctx *gin.Context, request *apicontract.UpdateMyInfoRequestDto)
}

type UserInfoController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewUserInfoController(
	coreAdapter *coreadapters.CoreAdapter,
) UserInfoControllerInterface {
	return &UserInfoController{
		coreAdapter: coreAdapter,
	}
}

func (c *UserInfoController) GetMyInfo(ctx *gin.Context, request *apicontract.GetMyInfoRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyInfoRequestDto,
		apicontract.GetMyInfoResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.GetMyInfoOperation,
		"/core/v1/user-infos/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *UserInfoController) UpdateMyInfo(ctx *gin.Context, request *apicontract.UpdateMyInfoRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyInfoRequestDto,
		apicontract.UpdateMyInfoResponseDto,
	](
		ctx,
		c.coreAdapter,
		request,
		apicontract.UpdateMyInfoOperation,
		"/core/v1/user-infos/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
