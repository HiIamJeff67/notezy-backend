package binders

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/api-keys"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
)

type APIKeyBinderInterface interface {
	BindCreateMyAPIKey(controllers.Func[*apicontract.CreateMyAPIKeyRequestDto]) gin.HandlerFunc
	BindListMyAPIKeys(controllers.Func[*apicontract.ListMyAPIKeysRequestDto]) gin.HandlerFunc
	BindRevokeMyAPIKey(controllers.Func[*apicontract.RevokeMyAPIKeyRequestDto]) gin.HandlerFunc
}

type APIKeyBinder struct{}

func NewAPIKeyBinder() APIKeyBinderInterface { return &APIKeyBinder{} }

func (b *APIKeyBinder) BindCreateMyAPIKey(controllerFunc controllers.Func[*apicontract.CreateMyAPIKeyRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateMyAPIKeyRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("APIKey").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, request)
	}
}

func (b *APIKeyBinder) BindListMyAPIKeys(controllerFunc controllers.Func[*apicontract.ListMyAPIKeysRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.ListMyAPIKeysRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, request)
	}
}

func (b *APIKeyBinder) BindRevokeMyAPIKey(controllerFunc controllers.Func[*apicontract.RevokeMyAPIKeyRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.RevokeMyAPIKeyRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		request.Param.PublicId = ctx.Param("api-key-id")
		controllerFunc(ctx, request)
	}
}
