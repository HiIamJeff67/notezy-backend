package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	routinetagsdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-tags"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type RoutineTagBinderInterface interface {
	BindGetMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutineTags(controllerFunc controllers.Func[*routinetagsdto.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc
	BindCreateRoutineTag(controllerFunc controllers.Func[*routinetagsdto.CreateRoutineTagRequestDto]) gin.HandlerFunc
	BindCreateRoutineTags(controllerFunc controllers.Func[*routinetagsdto.CreateRoutineTagsRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagsByIds(controllerFunc controllers.Func[*routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagsByIds(controllerFunc controllers.Func[*routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
}

type RoutineTagBinder struct{}

func NewRoutineTagBinder() RoutineTagBinderInterface {
	return &RoutineTagBinder{}
}

func (b *RoutineTagBinder) BindGetMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindGetAllMyRoutineTags(controllerFunc controllers.Func[*routinetagsdto.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindCreateRoutineTag(controllerFunc controllers.Func[*routinetagsdto.CreateRoutineTagRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindCreateRoutineTags(controllerFunc controllers.Func[*routinetagsdto.CreateRoutineTagsRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindUpdateMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindUpdateMyRoutineTagsByIds(controllerFunc controllers.Func[*routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagById(controllerFunc controllers.Func[*routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagsByIds(controllerFunc controllers.Func[*routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
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
