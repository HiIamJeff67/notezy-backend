package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	routinetagsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tags"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type RoutineTagEndpointInterface interface {
	GetMyRoutineTagById(ctx *gin.Context)
	GetAllMyRoutineTags(ctx *gin.Context)
	CreateRoutineTag(ctx *gin.Context)
	CreateRoutineTags(ctx *gin.Context)
	UpdateMyRoutineTagById(ctx *gin.Context)
	UpdateMyRoutineTagsByIds(ctx *gin.Context)
	HardDeleteMyRoutineTagById(ctx *gin.Context)
	HardDeleteMyRoutineTagsByIds(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchRoutineTags(ctx *gin.Context)
}

type RoutineTagEndpoint struct {
	routineTagService services.RoutineTagServiceInterface
}

func NewRoutineTagEndpoint(
	routineTagService services.RoutineTagServiceInterface,
) RoutineTagEndpointInterface {
	return &RoutineTagEndpoint{
		routineTagService: routineTagService,
	}
}

func (t *RoutineTagEndpoint) GetMyRoutineTagById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.GetMyRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.GetMyRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.GetMyRoutineTagByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) GetAllMyRoutineTags(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.GetAllMyRoutineTagsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.GetAllMyRoutineTags(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.GetAllMyRoutineTagsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) CreateRoutineTag(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.CreateRoutineTagRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.CreateRoutineTag(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.CreateRoutineTagResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) CreateRoutineTags(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.CreateRoutineTagsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.CreateRoutineTags(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.CreateRoutineTagsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) UpdateMyRoutineTagById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.UpdateMyRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.UpdateMyRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.UpdateMyRoutineTagByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) UpdateMyRoutineTagsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.UpdateMyRoutineTagsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.UpdateMyRoutineTagsByIdsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) HardDeleteMyRoutineTagById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.HardDeleteMyRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.HardDeleteMyRoutineTagByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RoutineTagEndpoint) HardDeleteMyRoutineTagsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTagService.HardDeleteMyRoutineTagsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetagsdto.HardDeleteMyRoutineTagsByIdsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
