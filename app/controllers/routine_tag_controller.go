package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTagControllerInterface interface {
	GetMyRoutineTagById(ctx *gin.Context, reqDto *dtos.GetMyRoutineTagByIdReqDto)
	CreateRoutineTag(ctx *gin.Context, reqDto *dtos.CreateRoutineTagReqDto)
	CreateRoutineTags(ctx *gin.Context, reqDto *dtos.CreateRoutineTagsReqDto)
	UpdateMyRoutineTagById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTagByIdReqDto)
	UpdateMyRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto)
	HardDeleteMyRoutineTagById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto)
	HardDeleteMyRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto)
}

type RoutineTagController struct {
	routineTagService services.RoutineTagServiceInterface
}

func NewRoutineTagController(service services.RoutineTagServiceInterface) RoutineTagControllerInterface {
	return &RoutineTagController{
		routineTagService: service,
	}
}

func (c *RoutineTagController) GetMyRoutineTagById(ctx *gin.Context, reqDto *dtos.GetMyRoutineTagByIdReqDto) {
	resDto, exception := c.routineTagService.GetMyRoutineTagById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) CreateRoutineTag(ctx *gin.Context, reqDto *dtos.CreateRoutineTagReqDto) {
	resDto, exception := c.routineTagService.CreateRoutineTag(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) CreateRoutineTags(ctx *gin.Context, reqDto *dtos.CreateRoutineTagsReqDto) {
	resDto, exception := c.routineTagService.CreateRoutineTags(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) UpdateMyRoutineTagById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTagByIdReqDto) {
	resDto, exception := c.routineTagService.UpdateMyRoutineTagById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) UpdateMyRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto) {
	resDto, exception := c.routineTagService.UpdateMyRoutineTagsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) HardDeleteMyRoutineTagById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto) {
	resDto, exception := c.routineTagService.HardDeleteMyRoutineTagById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *RoutineTagController) HardDeleteMyRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto) {
	resDto, exception := c.routineTagService.HardDeleteMyRoutineTagsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
