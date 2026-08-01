package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	useraccountsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-accounts"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type UserAccountControllerInterface interface {
	GetMyAccount(ctx *gin.Context, requestDto *useraccountsdto.GetMyAccountRequestDto)
	UpdateMyAccount(ctx *gin.Context, requestDto *useraccountsdto.UpdateMyAccountRequestDto)
	BindGoogleAccount(ctx *gin.Context, requestDto *useraccountsdto.BindGoogleAccountRequestDto)
	UnbindGoogleAccount(ctx *gin.Context, requestDto *useraccountsdto.UnbindGoogleAccountRequestDto)
}

type UserAccountController struct {
	coreClient *coreadapters.CoreClient
}

func NewUserAccountController(coreClient *coreadapters.CoreClient) UserAccountControllerInterface {
	return &UserAccountController{
		coreClient: coreClient,
	}
}

func (c *UserAccountController) GetMyAccount(
	ctx *gin.Context,
	requestDto *useraccountsdto.GetMyAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		useraccountsdto.GetMyAccountRequestDto,
		useraccountsdto.GetMyAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		useraccountsdto.GetMyAccountOperation,
		"/core/v1/user-accounts/get",
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

func (c *UserAccountController) UpdateMyAccount(
	ctx *gin.Context,
	requestDto *useraccountsdto.UpdateMyAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		useraccountsdto.UpdateMyAccountRequestDto,
		useraccountsdto.UpdateMyAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		useraccountsdto.UpdateMyAccountOperation,
		"/core/v1/user-accounts/update",
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

func (c *UserAccountController) BindGoogleAccount(
	ctx *gin.Context,
	requestDto *useraccountsdto.BindGoogleAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		useraccountsdto.BindGoogleAccountRequestDto,
		useraccountsdto.BindGoogleAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		useraccountsdto.BindGoogleAccountOperation,
		"/core/v1/user-accounts/google/bind",
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

func (c *UserAccountController) UnbindGoogleAccount(
	ctx *gin.Context,
	requestDto *useraccountsdto.UnbindGoogleAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		useraccountsdto.UnbindGoogleAccountRequestDto,
		useraccountsdto.UnbindGoogleAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		useraccountsdto.UnbindGoogleAccountOperation,
		"/core/v1/user-accounts/google/unbind",
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
