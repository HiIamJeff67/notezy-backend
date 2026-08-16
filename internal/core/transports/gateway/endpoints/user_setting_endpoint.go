package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/user-settings"
	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"

	userservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/user"
)

type UserSettingEndpointInterface interface {
	GetMySetting(ctx *gin.Context)
	UpdateMySetting(ctx *gin.Context)
}

type UserSettingEndpoint struct {
	userSettingService userservices.UserSettingServiceInterface
}

func NewUserSettingEndpoint(
	userSettingService userservices.UserSettingServiceInterface,
) UserSettingEndpointInterface {
	return &UserSettingEndpoint{
		userSettingService: userSettingService,
	}
}

func (t *UserSettingEndpoint) GetMySetting(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetMySettingRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userSettingService.GetMySetting(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetMySettingResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *UserSettingEndpoint) UpdateMySetting(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.UpdateMySettingRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userSettingService.UpdateMySetting(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.UpdateMySettingResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
