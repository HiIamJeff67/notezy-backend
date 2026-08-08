package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/users"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
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
	userService userservices.UserServiceInterface
}

func NewUserEndpoint(userService userservices.UserServiceInterface) UserEndpointInterface {
	return &UserEndpoint{userService: userService}
}

func (t *UserEndpoint) GetUserData(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetUserDataRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.GetUserData(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetUserDataResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}

func (t *UserEndpoint) GetMe(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetMeRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.GetMe(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetMeResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}

func (t *UserEndpoint) UpdateMe(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.UpdateMeRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, exception := t.userService.UpdateMe(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.UpdateMeResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *response,
	})
}
