package binders

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type StationBinderInterface interface {
	BindGetMyStationById(controllerFunc types.ControllerFunc[*dtos.GetMyStationByIdReqDto]) gin.HandlerFunc
	BindGetAllMyStations(controllerFunc types.ControllerFunc[*dtos.GetAllMyStationsReqDto]) gin.HandlerFunc
	BindCreateStation(controllerFunc types.ControllerFunc[*dtos.CreateStationReqDto]) gin.HandlerFunc
	BindCreateStations(controllerFunc types.ControllerFunc[*dtos.CreateStationsReqDto]) gin.HandlerFunc
	BindUpdateMyStationById(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationByIdReqDto]) gin.HandlerFunc
	BindUpdateMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationsByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyStationById(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationByIdReqDto]) gin.HandlerFunc
	BindRestoreMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationsByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationByIdReqDto]) gin.HandlerFunc
	BindDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationsByIdsReqDto]) gin.HandlerFunc
	BindHardDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationsByIdsReqDto]) gin.HandlerFunc
	BindVisualizeMyTotalCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyTotalCountReqDto]) gin.HandlerFunc
}

type StationBinder struct{}

func NewStationBinder() StationBinderInterface {
	return &StationBinder{}
}

func (b *StationBinder) BindGetMyStationById(controllerFunc types.ControllerFunc[*dtos.GetMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.IsDeleted = &isDeleted
		}

		stationIdString := ctx.Query("stationId")
		if stationIdString == "" {
			exceptions.Station.InvalidInput().WithOrigin(fmt.Errorf("stationId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		stationId, err := uuid.Parse(stationIdString)
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindGetAllMyStations(controllerFunc types.ControllerFunc[*dtos.GetAllMyStationsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyStationsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindCreateStation(controllerFunc types.ControllerFunc[*dtos.CreateStationReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateStationReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindCreateStations(controllerFunc types.ControllerFunc[*dtos.CreateStationsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateStationsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpdateMyStationById(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpdateMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindRestoreMyStationById(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindRestoreMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindHardDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindHardDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindVisualizeMyTotalCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyTotalCountReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyTotalCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.Station.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		controllerFunc(ctx, &reqDto)
	}
}
