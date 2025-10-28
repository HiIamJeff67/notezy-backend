package binders

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/models/schemas/enums"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type MaterialBinderInterface interface {
	BindGetMyMaterialById(controllerFunc types.ControllerFunc[*dtos.GetMyMaterialByIdReqDto]) gin.HandlerFunc
	BindGetMyMaterialAndItsParentById(controllerFunc types.ControllerFunc[*dtos.GetMyMaterialAndItsParentByIdReqDto]) gin.HandlerFunc
	BindGetAllMyMaterialsByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyMaterialsByParentSubShelfIdReqDto]) gin.HandlerFunc
	BindGetAllMyMaterialsByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyMaterialsByRootShelfIdReqDto]) gin.HandlerFunc
	BindCreateTextbookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateTextbookMaterialReqDto]) gin.HandlerFunc
	BindCreateNotebookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateNotebookMaterialReqDto]) gin.HandlerFunc
	BindUpdateMyMaterialById(controllerFunc types.ControllerFunc[*dtos.UpdateMyMaterialByIdReqDto]) gin.HandlerFunc
	BindSaveMyNotebookMaterialById(controllerFunc types.ControllerFunc[*dtos.SaveMyMaterialByIdReqDto]) gin.HandlerFunc
	BindMoveMyMaterialById(controllerFunc types.ControllerFunc[*dtos.MoveMyMaterialByIdReqDto]) gin.HandlerFunc
	BindMoveMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyMaterialsByIdsReqDto]) gin.HandlerFunc
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

		// for the uuid in the parameter, we MUST extract it and parse it manually
		materialIdString := ctx.Query("materialId")
		if materialIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("materialId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		materialId, err := uuid.Parse(materialIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.MaterialId = materialId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindGetMyMaterialAndItsParentById(controllerFunc types.ControllerFunc[*dtos.GetMyMaterialAndItsParentByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyMaterialAndItsParentByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		// for the uuid in the parameter, we MUST extract it and parse it manually
		materialIdString := ctx.Query("materialId")
		if materialIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("materialId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		materialId, err := uuid.Parse(materialIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.MaterialId = materialId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindGetAllMyMaterialsByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyMaterialsByParentSubShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyMaterialsByParentSubShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		// for the uuid in the parameter, we MUST extract it and parse it manually
		parentSubShelfIdString := ctx.Query("parentSubShelfId")
		if parentSubShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("parentSubShelfId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		parentSubShelfId, err := uuid.Parse(parentSubShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.ParentSubShelfId = parentSubShelfId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Material.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindGetAllMyMaterialsByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyMaterialsByRootShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyMaterialsByRootShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		// for the uuid in the parameter, we MUST extract it and parse it manually
		rootShelfIdString := ctx.Query("rootShelfId")
		if rootShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("rootShelfId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		rootShelfId, err := uuid.Parse(rootShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *MaterialBinder) BindCreateTextbookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateTextbookMaterialReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateTextbookMaterialReqDto

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

func (b *MaterialBinder) BindCreateNotebookMaterial(controllerFunc types.ControllerFunc[*dtos.CreateNotebookMaterialReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateNotebookMaterialReqDto

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

func (b *MaterialBinder) BindUpdateMyMaterialById(controllerFunc types.ControllerFunc[*dtos.UpdateMyMaterialByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyMaterialByIdReqDto

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

func (b *MaterialBinder) BindSaveMyTextbookMaterialById(controllerFunc types.ControllerFunc[*dtos.SaveMyMaterialByIdReqDto]) gin.HandlerFunc {
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

			fileHeader := fileHeaders[0]

			// check the file size
			if fileHeader.Size > constants.MaxNotebookFileSize {
				exceptions.Material.FileTooLarge(fileHeader.Size, constants.MaxNotebookFileSize).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// try to open the file
			fileInterface, err := fileHeader.Open()
			if err != nil {
				exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// peek the first few bytes of the given file
			bufferedReader := bufio.NewReader(fileInterface)
			peekedBytes, err := bufferedReader.Peek(int(constants.PeekFileSize))
			if err != nil && err != io.EOF {
				fileInterface.Close()
				exceptions.Material.CannotPeekFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// detect the content type from the peeked bytes
			detectedContentType := http.DetectContentType(peekedBytes)
			// note that the content type may be detected in text/plain sometimes
			if !strings.HasPrefix(detectedContentType, enums.MaterialContentType_PlainText.String()) &&
				!strings.HasPrefix(detectedContentType, enums.MaterialContentType_JSON.String()) {
				fileInterface.Close()
				exceptions.Material.InvalidType(detectedContentType).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// try to reset(restore) the reading pointer to the beginning of the file
			if seeker, ok := fileInterface.(io.Seeker); ok { // try using seeker to do so
				_, err := seeker.Seek(0, io.SeekStart)
				if err != nil {
					fileInterface.Close()
					exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
					return
				}
			} else { // if it cannot be seeked, then re-open the file
				fileInterface.Close()
				fileInterface, err = fileHeader.Open()
				if err != nil {
					exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
					return
				}
			}

			reqDto.Body.ContentFile = fileInterface // bind the file interface here
			reqDto.ContextFields.Size = &fileHeaders[0].Size

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

func (b *MaterialBinder) BindSaveMyNotebookMaterialById(controllerFunc types.ControllerFunc[*dtos.SaveMyMaterialByIdReqDto]) gin.HandlerFunc {
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

			fileHeader := fileHeaders[0]

			// check the file size
			if fileHeader.Size > constants.MaxNotebookFileSize {
				exceptions.Material.FileTooLarge(fileHeader.Size, constants.MaxNotebookFileSize).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// try to open the file
			fileInterface, err := fileHeader.Open()
			if err != nil {
				exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// peek the first few bytes of the given file
			bufferedReader := bufio.NewReader(fileInterface)
			peekedBytes, err := bufferedReader.Peek(int(constants.PeekFileSize))
			if err != nil && err != io.EOF {
				fileInterface.Close()
				exceptions.Material.CannotPeekFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// detect the content type from the peeked bytes
			detectedContentType := http.DetectContentType(peekedBytes)
			// note that the content type may be detected in text/plain sometimes
			if !strings.HasPrefix(detectedContentType, enums.MaterialContentType_PlainText.String()) &&
				!strings.HasPrefix(detectedContentType, enums.MaterialContentType_JSON.String()) {
				fileInterface.Close()
				exceptions.Material.InvalidType(detectedContentType).Log().SafelyResponseWithJSON(ctx)
				return
			}

			// try to reset(restore) the reading pointer to the beginning of the file
			if seeker, ok := fileInterface.(io.Seeker); ok { // try using seeker to do so
				_, err := seeker.Seek(0, io.SeekStart)
				if err != nil {
					fileInterface.Close()
					exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
					return
				}
			} else { // if it cannot be seeked, then re-open the file
				fileInterface.Close()
				fileInterface, err = fileHeader.Open()
				if err != nil {
					exceptions.Material.CannotOpenFiles().WithError(err).Log().SafelyResponseWithJSON(ctx)
					return
				}
			}

			reqDto.Body.ContentFile = fileInterface // bind the file interface here
			reqDto.ContextFields.Size = &fileHeaders[0].Size

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

func (b *MaterialBinder) BindMoveMyMaterialsByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyMaterialsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyMaterialsByIdsReqDto

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
