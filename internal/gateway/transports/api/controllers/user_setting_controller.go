package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	usersettingsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-settings"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type UserSettingControllerInterface interface {
	GetMySetting(ctx *gin.Context, requestDto *usersettingsdto.GetMySettingRequestDto)
	UpdateMySetting(ctx *gin.Context, requestDto *usersettingsdto.UpdateMySettingRequestDto)
}

type UserSettingController struct {
	coreClient *coreadapters.CoreClient
}

func NewUserSettingController(coreClient *coreadapters.CoreClient) UserSettingControllerInterface {
	return &UserSettingController{
		coreClient: coreClient,
	}
}

func (c *UserSettingController) GetMySetting(ctx *gin.Context, requestDto *usersettingsdto.GetMySettingRequestDto) {
	response, exception := coreadapters.CallSecurly[
		usersettingsdto.GetMySettingRequestDto,
		usersettingsdto.GetMySettingResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		usersettingsdto.GetMySettingOperation,
		"/core/v1/user-settings/get",
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

func (c *UserSettingController) UpdateMySetting(ctx *gin.Context, requestDto *usersettingsdto.UpdateMySettingRequestDto) {
	response, exception := coreadapters.CallSecurly[
		usersettingsdto.UpdateMySettingRequestDto,
		usersettingsdto.UpdateMySettingResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		usersettingsdto.UpdateMySettingOperation,
		"/core/v1/user-settings/update",
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
