package binders

import (
	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type AuthBinderInterface interface {
	BindRegister(controllerFunc types.ControllerFunc[*dtos.RegisterReqDto]) gin.HandlerFunc
	BindLogin(controllerFunc types.ControllerFunc[*dtos.LoginReqDto]) gin.HandlerFunc
	BindLogout(controllerFunc types.ControllerFunc[*dtos.LogoutReqDto]) gin.HandlerFunc
	BindSendAuthCode(controllerFunc types.ControllerFunc[*dtos.SendAuthCodeReqDto]) gin.HandlerFunc
	BindValidateEmail(controllerFunc types.ControllerFunc[*dtos.ValidateEmailReqDto]) gin.HandlerFunc
	BindResetEmail(controllerFunc types.ControllerFunc[*dtos.ResetEmailReqDto]) gin.HandlerFunc
	BindForgetPassword(controllerFunc types.ControllerFunc[*dtos.ForgetPasswordReqDto]) gin.HandlerFunc
	BindDeleteMe(controllerFunc types.ControllerFunc[*dtos.DeleteMeReqDto]) gin.HandlerFunc
}

type AuthBinder struct{}

func NewAuthBinder() AuthBinderInterface {
	return &AuthBinder{}
}

/* ============================== Binder ============================== */

func (b *AuthBinder) BindRegister(controllerFunc types.ControllerFunc[*dtos.RegisterReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RegisterReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindLogin(controllerFunc types.ControllerFunc[*dtos.LoginReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LoginReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindLogout(controllerFunc types.ControllerFunc[*dtos.LogoutReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LogoutReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindSendAuthCode(controllerFunc types.ControllerFunc[*dtos.SendAuthCodeReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SendAuthCodeReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindValidateEmail(controllerFunc types.ControllerFunc[*dtos.ValidateEmailReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ValidateEmailReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindResetEmail(controllerFunc types.ControllerFunc[*dtos.ResetEmailReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ResetEmailReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindForgetPassword(controllerFunc types.ControllerFunc[*dtos.ForgetPasswordReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ForgetPasswordReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *AuthBinder) BindDeleteMe(controllerFunc types.ControllerFunc[*dtos.DeleteMeReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMeReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Auth.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
