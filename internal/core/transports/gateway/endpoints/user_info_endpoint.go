package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
)

type UserInfoEndpointInterface interface {
	GetMyInfo(ctx *gin.Context)
	UpdateMyInfo(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	LoadUserInfos(ctx *gin.Context)
}

type UserInfoEndpoint struct {
	userInfoService userservices.UserInfoServiceInterface
}

func NewUserInfoEndpoint(userInfoService userservices.UserInfoServiceInterface) UserInfoEndpointInterface {
	return &UserInfoEndpoint{userInfoService: userInfoService}
}

func (t *UserInfoEndpoint) GetMyInfo(ctx *gin.Context) {
	request := &gatewaycontract.Request[userinfosdto.GetMyInfoRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userInfoService.GetMyInfo(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[userinfosdto.GetMyInfoResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *UserInfoEndpoint) UpdateMyInfo(ctx *gin.Context) {
	request := &gatewaycontract.Request[userinfosdto.UpdateMyInfoRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userInfoService.UpdateMyInfo(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[userinfosdto.UpdateMyInfoResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
