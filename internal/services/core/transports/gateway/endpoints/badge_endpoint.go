package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	badgesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/badges"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type BadgeEndpointInterface interface {
	LoadUserBadges(ctx *gin.Context)
}

type BadgeEndpoint struct {
	badgeService services.BadgeServiceInterface
}

func NewBadgeEndpoint(
	badgeService services.BadgeServiceInterface,
) BadgeEndpointInterface {
	return &BadgeEndpoint{
		badgeService: badgeService,
	}
}

func (t *BadgeEndpoint) LoadUserBadges(ctx *gin.Context) {
	request := &core.Request[badgesdto.LoadUserBadgesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDtos, exception := t.badgeService.GetPublicBadgesByUserPublicIds(ctx.Request.Context(), request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[badgesdto.LoadUserBadgesResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: responseDtos,
	})
}
