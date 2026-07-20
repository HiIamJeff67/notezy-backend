package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RealtimeControllerInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context, reqDto *dtos.GetMyBlockPackRealtimeParticipantsReqDto)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context, reqDto *dtos.CreateMyRealtimeConnectionTicketReqDto)
	CreateMyBlockPackChannelTicket(ctx *gin.Context, reqDto *dtos.CreateMyBlockPackChannelTicketReqDto)
}

type RealtimeController struct {
	realtimeService services.RealtimeServiceInterface
}

func NewRealtimeController(service services.RealtimeServiceInterface) RealtimeControllerInterface {
	return &RealtimeController{
		realtimeService: service,
	}
}

func (c *RealtimeController) GetMyBlockPackRealtimeParticipants(
	ctx *gin.Context, reqDto *dtos.GetMyBlockPackRealtimeParticipantsReqDto,
) {
	resDto, exception := c.realtimeService.GetMyBlockPackRealtimeParticipants(ctx.Request.Context(), reqDto)
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

func (c *RealtimeController) CreateMyRealtimeConnectionTicket(
	ctx *gin.Context,
	reqDto *dtos.CreateMyRealtimeConnectionTicketReqDto,
) {
	resDto, exception := c.realtimeService.CreateMyRealtimeConnectionTicket(ctx.Request.Context(), reqDto)
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

func (c *RealtimeController) CreateMyBlockPackChannelTicket(
	ctx *gin.Context,
	reqDto *dtos.CreateMyBlockPackChannelTicketReqDto,
) {
	resDto, exception := c.realtimeService.CreateMyBlockPackChannelTicket(ctx.Request.Context(), reqDto)
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
