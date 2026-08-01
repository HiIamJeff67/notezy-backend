package endpoints

import (
	"net/http"
	"time"

	usersdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/users"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	"github.com/gin-gonic/gin"
)

type UserEndpointInterface interface {
	GetUserData(*gin.Context)
	GetMe(*gin.Context)
	UpdateMe(*gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchUsers(*gin.Context)
	LoadThemeAuthors(*gin.Context)
}

type UserEndpoint struct {
	userService services.UserServiceInterface
}

func NewUserEndpoint(userService services.UserServiceInterface) UserEndpointInterface {
	return &UserEndpoint{userService: userService}
}

func (t *UserEndpoint) GetUserData(ctx *gin.Context) {
	request := &core.Request[usersdto.GetUserDataRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.GetUserData(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[usersdto.GetUserDataResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}

func (t *UserEndpoint) GetMe(ctx *gin.Context) {
	request := &core.Request[usersdto.GetMeRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.GetMe(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[usersdto.GetMeResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}

func (t *UserEndpoint) UpdateMe(ctx *gin.Context) {
	request := &core.Request[usersdto.UpdateMeRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.UpdateMe(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[usersdto.UpdateMeResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}
