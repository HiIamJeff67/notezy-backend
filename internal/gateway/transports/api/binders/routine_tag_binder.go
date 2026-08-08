package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tags"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type RoutineTagBinderInterface interface {
	BindGetMyRoutineTagById(controllerFunc controllers.Func[*apicontract.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutineTags(controllerFunc controllers.Func[*apicontract.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc
	BindCreateRoutineTag(controllerFunc controllers.Func[*apicontract.CreateRoutineTagRequestDto]) gin.HandlerFunc
	BindCreateRoutineTags(controllerFunc controllers.Func[*apicontract.CreateRoutineTagsRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagById(controllerFunc controllers.Func[*apicontract.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagsByIds(controllerFunc controllers.Func[*apicontract.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagById(controllerFunc controllers.Func[*apicontract.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagsByIds(controllerFunc controllers.Func[*apicontract.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc
}

type RoutineTagBinder struct{}

func NewRoutineTagBinder() RoutineTagBinderInterface {
	return &RoutineTagBinder{}
}

func (b *RoutineTagBinder) BindGetMyRoutineTagById(controllerFunc controllers.Func[*apicontract.GetMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &isDeleted
		}

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindGetAllMyRoutineTags(controllerFunc controllers.Func[*apicontract.GetAllMyRoutineTagsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetAllMyRoutineTagsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTag(controllerFunc controllers.Func[*apicontract.CreateRoutineTagRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.CreateRoutineTagRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTags(controllerFunc controllers.Func[*apicontract.CreateRoutineTagsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.CreateRoutineTagsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagById(controllerFunc controllers.Func[*apicontract.UpdateMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagsByIds(controllerFunc controllers.Func[*apicontract.UpdateMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMyRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagById(controllerFunc controllers.Func[*apicontract.HardDeleteMyRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.HardDeleteMyRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		routineTagId, err := uuid.Parse(ctx.Param("routineTagId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTag").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagsByIds(controllerFunc controllers.Func[*apicontract.HardDeleteMyRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.HardDeleteMyRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTag").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
