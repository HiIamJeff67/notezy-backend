package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	usersettingsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-settings"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type UserSettingEndpointInterface interface {
	GetMySetting(ctx *gin.Context)
	UpdateMySetting(ctx *gin.Context)
}

type UserSettingEndpoint struct {
	userSettingService services.UserSettingServiceInterface
}

func NewUserSettingEndpoint(
	userSettingService services.UserSettingServiceInterface,
) UserSettingEndpointInterface {
	return &UserSettingEndpoint{
		userSettingService: userSettingService,
	}
}

func (t *UserSettingEndpoint) GetMySetting(ctx *gin.Context) {
	request := &gatewaycontract.Request[usersettingsdto.GetMySettingRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[usersettingsdto.GetMySettingResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *UserSettingEndpoint) UpdateMySetting(ctx *gin.Context) {
	request := &gatewaycontract.Request[usersettingsdto.UpdateMySettingRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[usersettingsdto.UpdateMySettingResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
