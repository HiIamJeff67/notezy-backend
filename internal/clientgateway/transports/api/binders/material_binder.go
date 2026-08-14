package binders

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/materials"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
)

type MaterialBinderInterface interface {
	BindGetMyMaterialById(controllerFunc controllers.Func[*apicontract.GetMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindGetMyMaterialAndItsParentById(controllerFunc controllers.Func[*apicontract.GetMyMaterialAndItsParentByIdRequestDto]) gin.HandlerFunc
	BindGetMyMaterialsByParentSubShelfId(controllerFunc controllers.Func[*apicontract.GetMyMaterialsByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMyMaterialsByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMyMaterialsByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateMyMaterial(controllerFunc controllers.Func[*apicontract.CreateMyMaterialRequestDto]) gin.HandlerFunc
	BindUpdateMyMaterialById(controllerFunc controllers.Func[*apicontract.UpdateMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindSaveMyMaterialById(controllerFunc controllers.Func[*apicontract.SaveMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindMoveMyMaterialById(controllerFunc controllers.Func[*apicontract.MoveMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindMoveMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.MoveMyMaterialsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyMaterialById(controllerFunc controllers.Func[*apicontract.RestoreMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.RestoreMyMaterialsByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyMaterialById(controllerFunc controllers.Func[*apicontract.DeleteMyMaterialByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.DeleteMyMaterialsByIdsRequestDto]) gin.HandlerFunc
}

type MaterialBinder struct{}

func NewMaterialBinder() MaterialBinderInterface {
	return &MaterialBinder{}
}

func (b *MaterialBinder) BindGetMyMaterialById(controllerFunc controllers.Func[*apicontract.GetMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("isDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetMyMaterialAndItsParentById(controllerFunc controllers.Func[*apicontract.GetMyMaterialAndItsParentByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyMaterialAndItsParentByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("isDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetMyMaterialsByParentSubShelfId(controllerFunc controllers.Func[*apicontract.GetMyMaterialsByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyMaterialsByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("areDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("parent-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindGetAllMyMaterialsByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMyMaterialsByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetAllMyMaterialsByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		valueString := ctx.Query("areDeleted")
		if valueString != "" {
			value, err := strconv.ParseBool(valueString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindCreateMyMaterial(controllerFunc controllers.Func[*apicontract.CreateMyMaterialRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.CreateMyMaterialRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("parent-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindUpdateMyMaterialById(controllerFunc controllers.Func[*apicontract.UpdateMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.UpdateMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindSaveMyMaterialById(controllerFunc controllers.Func[*apicontract.SaveMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.SaveMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		fileHeader, err := ctx.FormFile("contentFile")
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		defer file.Close()

		contentFile, err := io.ReadAll(file)
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ContentFile = contentFile

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindMoveMyMaterialById(controllerFunc controllers.Func[*apicontract.MoveMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindMoveMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.MoveMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialById(controllerFunc controllers.Func[*apicontract.RestoreMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindRestoreMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.RestoreMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialById(controllerFunc controllers.Func[*apicontract.DeleteMyMaterialByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMyMaterialByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("material-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Material").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.MaterialId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *MaterialBinder) BindDeleteMyMaterialsByIds(controllerFunc controllers.Func[*apicontract.DeleteMyMaterialsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMyMaterialsByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Material").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
