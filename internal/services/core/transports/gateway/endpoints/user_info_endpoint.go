package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-infos"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type UserInfoEndpointInterface interface {
	GetMyInfo(ctx *gin.Context)
	UpdateMyInfo(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	LoadUserInfos(ctx *gin.Context)
}

type UserInfoEndpoint struct {
	userInfoService services.UserInfoServiceInterface
}

func NewUserInfoEndpoint(userInfoService services.UserInfoServiceInterface) UserInfoEndpointInterface {
	return &UserInfoEndpoint{userInfoService: userInfoService}
}

func (t *UserInfoEndpoint) GetMyInfo(ctx *gin.Context) {
	request := &core.Request[userinfosdto.GetMyInfoRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userInfoService.GetMyInfo(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[userinfosdto.GetMyInfoResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *UserInfoEndpoint) UpdateMyInfo(ctx *gin.Context) {
	request := &core.Request[userinfosdto.UpdateMyInfoRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.userInfoService.UpdateMyInfo(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[userinfosdto.UpdateMyInfoResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
