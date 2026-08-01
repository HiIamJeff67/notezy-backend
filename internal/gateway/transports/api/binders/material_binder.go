package binders

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	materialsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/materials"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type MaterialBinderInterface interface {
	BindGetMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindGetMyMaterialAndItsParentById(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialAndItsParentByIdRequestDto]) gin.HandlerFunc
	BindGetMyMaterialsByParentSubShelfId(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMyMaterialsByRootShelfId(controllerFunc apitransport.ControllerFunc[*materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateMyMaterial(controllerFunc apitransport.ControllerFunc[*materialsdto.CreateMyMaterialRequestDto]) gin.HandlerFunc
	BindUpdateMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.UpdateMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindSaveMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.SaveMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindMoveMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.MoveMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindMoveMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.MoveMyMaterialsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.RestoreMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.RestoreMyMaterialsByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.DeleteMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.DeleteMyMaterialsByIdsRequestDto]) gin.HandlerFunc
}

type MaterialBinder struct{}

func NewMaterialBinder() MaterialBinderInterface {
	return &MaterialBinder{}
}

func (b *MaterialBinder) BindGetMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.GetMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("isDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetMyMaterialAndItsParentById(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialAndItsParentByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.GetMyMaterialAndItsParentByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("isDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetMyMaterialsByParentSubShelfId(controllerFunc apitransport.ControllerFunc[*materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("areDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("parentSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetAllMyMaterialsByRootShelfId(controllerFunc apitransport.ControllerFunc[*materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("areDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindCreateMyMaterial(controllerFunc apitransport.ControllerFunc[*materialsdto.CreateMyMaterialRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.CreateMyMaterialRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("parentSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindUpdateMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.UpdateMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.UpdateMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindSaveMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.SaveMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.SaveMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		fileHeader, err := ctx.FormFile("contentFile")
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		defer file.Close()

		contentFile, err := io.ReadAll(file)
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ContentFile = contentFile

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindMoveMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.MoveMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.MoveMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindMoveMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.MoveMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.MoveMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.RestoreMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.RestoreMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.RestoreMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.RestoreMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialById(controllerFunc apitransport.ControllerFunc[*materialsdto.DeleteMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.DeleteMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("materialId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialsByIds(controllerFunc apitransport.ControllerFunc[*materialsdto.DeleteMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto materialsdto.DeleteMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
