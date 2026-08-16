package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/badges"
	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"

	otherservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/other"
)

type BadgeEndpointInterface interface {
	LoadUserBadges(ctx *gin.Context)
}

type BadgeEndpoint struct {
	badgeService otherservices.BadgeServiceInterface
}

func NewBadgeEndpoint(
	badgeService otherservices.BadgeServiceInterface,
) BadgeEndpointInterface {
	return &BadgeEndpoint{
		badgeService: badgeService,
	}
}

func (t *BadgeEndpoint) LoadUserBadges(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LoadUserBadgesRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDtos, exception := t.badgeService.GetPublicBadgesByUserPublicIds(ctx.Request.Context(), request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LoadUserBadgesResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: responseDtos,
	})
}
