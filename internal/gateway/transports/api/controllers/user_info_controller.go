package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/exceptionwriter"
)

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context, request *userinfosdto.GetMyInfoRequestDto)
	UpdateMyInfo(ctx *gin.Context, request *userinfosdto.UpdateMyInfoRequestDto)
}

type UserInfoController struct {
	coreClient *coreadapters.CoreClient
}

func NewUserInfoController(
	coreClient *coreadapters.CoreClient,
) UserInfoControllerInterface {
	return &UserInfoController{
		coreClient: coreClient,
	}
}

func (c *UserInfoController) GetMyInfo(ctx *gin.Context, request *userinfosdto.GetMyInfoRequestDto) {
	response, exception := coreadapters.CallSecurly[
		userinfosdto.GetMyInfoRequestDto,
		userinfosdto.GetMyInfoResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		userinfosdto.GetMyInfoOperation,
		"/core/v1/user-infos/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *UserInfoController) UpdateMyInfo(ctx *gin.Context, request *userinfosdto.UpdateMyInfoRequestDto) {
	response, exception := coreadapters.CallSecurly[
		userinfosdto.UpdateMyInfoRequestDto,
		userinfosdto.UpdateMyInfoResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		userinfosdto.UpdateMyInfoOperation,
		"/core/v1/user-infos/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}
