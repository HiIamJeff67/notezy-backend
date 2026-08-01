package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	routinetagsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-tags"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type RoutineTagBinderInterface interface {
	BindGetMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutineTags(controllerFunc apitransport.ControllerFunc[*routinetagsdto.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc
	BindCreateRoutineTag(controllerFunc apitransport.ControllerFunc[*routinetagsdto.CreateRoutineTagRequestDto]) gin.HandlerFunc
	BindCreateRoutineTags(controllerFunc apitransport.ControllerFunc[*routinetagsdto.CreateRoutineTagsRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
}

type RoutineTagBinder struct{}

func NewRoutineTagBinder() RoutineTagBinderInterface {
	return &RoutineTagBinder{}
}

func (b *RoutineTagBinder) BindGetMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.GetMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &isDeleted
		}

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindGetAllMyRoutineTags(controllerFunc apitransport.ControllerFunc[*routinetagsdto.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.GetAllMyRoutineTagsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTag(controllerFunc apitransport.ControllerFunc[*routinetagsdto.CreateRoutineTagRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.CreateRoutineTagRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTags(controllerFunc apitransport.ControllerFunc[*routinetagsdto.CreateRoutineTagsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.CreateRoutineTagsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.UpdateMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
