package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-accounts"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type UserAccountControllerInterface interface {
	GetMyAccount(ctx *gin.Context, requestDto *apicontract.GetMyAccountRequestDto)
	UpdateMyAccount(ctx *gin.Context, requestDto *apicontract.UpdateMyAccountRequestDto)
	BindGoogleAccount(ctx *gin.Context, requestDto *apicontract.BindGoogleAccountRequestDto)
	UnbindGoogleAccount(ctx *gin.Context, requestDto *apicontract.UnbindGoogleAccountRequestDto)
}

type UserAccountController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewUserAccountController(coreClient *coreadapters.CoreAdapter) UserAccountControllerInterface {
	return &UserAccountController{
		coreClient: coreClient,
	}
}

func (c *UserAccountController) GetMyAccount(
	ctx *gin.Context,
	requestDto *apicontract.GetMyAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyAccountRequestDto,
		apicontract.GetMyAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyAccountOperation,
		"/core/v1/user-accounts/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *UserAccountController) UpdateMyAccount(
	ctx *gin.Context,
	requestDto *apicontract.UpdateMyAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyAccountRequestDto,
		apicontract.UpdateMyAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.UpdateMyAccountOperation,
		"/core/v1/user-accounts/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *UserAccountController) BindGoogleAccount(
	ctx *gin.Context,
	requestDto *apicontract.BindGoogleAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.BindGoogleAccountRequestDto,
		apicontract.BindGoogleAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.BindGoogleAccountOperation,
		"/core/v1/user-accounts/google/bind",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *UserAccountController) UnbindGoogleAccount(
	ctx *gin.Context,
	requestDto *apicontract.UnbindGoogleAccountRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UnbindGoogleAccountRequestDto,
		apicontract.UnbindGoogleAccountResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.UnbindGoogleAccountOperation,
		"/core/v1/user-accounts/google/unbind",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
