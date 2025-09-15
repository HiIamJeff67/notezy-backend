package binders

import (
	"io"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type MaterialBinderInterface interface {
	BindGetMyMaterialById(controllerFunc types.ControllerFunc[*dtos.GetMyMaterialByIdReqDto]) gin.HandlerFunc
	BindSearchMyMaterialsByShelfId(controllerFunc types.ControllerFunc[*dtos.SearchMyMaterialsByShelfIdReqDto]) gin.HandlerFunc
	BindCreateTextbookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateMaterialReqDto]) gin.HandlerFunc
	BindSaveMyMaterialById(controllerFunc types.ControllerFunc[*dtos.SaveMyMaterialByIdReqDto]) gin.HandlerFunc
	BindMoveMyMaterialById(controllerFunc types.ControllerFunc[*dtos.MoveMyMaterialByIdReqDto]) gin.HandlerFunc
	BindRestoreMyMaterialById(controllerFunc types.ControllerFunc[*dtos.RestoreMyMaterialByIdReqDto]) gin.HandlerFunc
	BindRestoreMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyMaterialsByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyMaterialById(controllerFunc types.ControllerFunc[*dtos.DeleteMyMaterialByIdReqDto]) gin.HandlerFunc
	BindDeleteMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyMaterialsByIdsReqDto]) gin.HandlerFunc
}

type MaterialBinder struct{}

func NewMaterialBinder() MaterialBinderInterface {
	return &MaterialBinder{}
}

/* ============================== Binder ============================== */

func (b *MaterialBinder) BindGetMyMaterialById(controllerFunc types.ControllerFunc[*dtos.GetMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyMaterialByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exception.Log()
			exceptions.User.InvalidInput().WithError(err).ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindSearchMyMaterialsByShelfId(controllerFunc types.ControllerFunc[*dtos.SearchMyMaterialsByShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SearchMyMaterialsByShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exception.Log()
			exceptions.User.InvalidInput().WithError(err).ResponseWithJSON(ctx)
			return
		}

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindCreateTextbookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateMaterialReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMaterialReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		userPublicId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_PublicId)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserPublicId = *userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindSaveMyMaterialById(controllerFunc types.ControllerFunc[*dtos.SaveMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SaveMyMaterialByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		userPublicId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_PublicId)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserPublicId = *userPublicId

		// extract the fileHeader from the context field, and make sure it's only one fileHeader
		fileHeaders, exception := contexts.GetAndConvertContextToMultipartFileHeaders(ctx)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		var numberOfFiles = int64(len(fileHeaders))

		if numberOfFiles > 0 {
			if numberOfFiles > 1 {
				exceptions.Material.TooManyFiles(numberOfFiles).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// try to open the file
			fileInterface, err := fileHeaders[0].Open()
			if err != nil {
				exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
				return
			}
			reqDto.Body.ContentFile = fileInterface // bind the file interface here
			reqDto.Body.Size = &fileHeaders[0].Size

			// make sure the file is closed at the end
			defer func(f io.Reader) {
				if closer, ok := f.(io.Closer); ok {
					closer.Close()
				}
			}(fileInterface)
		}

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindMoveMyMaterialById(controllerFunc types.ControllerFunc[*dtos.MoveMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyMaterialByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialById(controllerFunc types.ControllerFunc[*dtos.RestoreMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyMaterialByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyMaterialsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyMaterialsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialById(controllerFunc types.ControllerFunc[*dtos.DeleteMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyMaterialByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyMaterialsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyMaterialsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Material.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
