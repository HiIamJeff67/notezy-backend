package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context, request *apicontract.GetMyInfoRequestDto)
	UpdateMyInfo(ctx *gin.Context, request *apicontract.UpdateMyInfoRequestDto)
}

type UserInfoController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewUserInfoController(
	coreClient *coreadapters.CoreAdapter,
) UserInfoControllerInterface {
	return &UserInfoController{
		coreClient: coreClient,
	}
}

func (c *UserInfoController) GetMyInfo(ctx *gin.Context, request *apicontract.GetMyInfoRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyInfoRequestDto,
		apicontract.GetMyInfoResponseDto,
	](
		ctx,
		c.coreClient,
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
		c.coreClient,
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
