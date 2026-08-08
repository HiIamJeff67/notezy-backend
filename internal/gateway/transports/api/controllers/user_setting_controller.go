package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-settings"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type UserSettingControllerInterface interface {
	GetMySetting(ctx *gin.Context, requestDto *apicontract.GetMySettingRequestDto)
	UpdateMySetting(ctx *gin.Context, requestDto *apicontract.UpdateMySettingRequestDto)
}

type UserSettingController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewUserSettingController(coreClient *coreadapters.CoreAdapter) UserSettingControllerInterface {
	return &UserSettingController{
		coreClient: coreClient,
	}
}

func (c *UserSettingController) GetMySetting(ctx *gin.Context, requestDto *apicontract.GetMySettingRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMySettingRequestDto,
		apicontract.GetMySettingResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMySettingOperation,
		"/core/v1/user-settings/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *UserSettingController) UpdateMySetting(ctx *gin.Context, requestDto *apicontract.UpdateMySettingRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMySettingRequestDto,
		apicontract.UpdateMySettingResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.UpdateMySettingOperation,
		"/core/v1/user-settings/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
