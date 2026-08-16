package controllers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/api-keys"
	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type APIKeyControllerInterface interface {
	CreateMyAPIKey(*gin.Context, *apicontract.CreateMyAPIKeyRequestDto)
	ListMyAPIKeys(*gin.Context, *apicontract.ListMyAPIKeysRequestDto)
	RevokeMyAPIKey(*gin.Context, *apicontract.RevokeMyAPIKeyRequestDto)
}

type APIKeyController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewAPIKeyController(coreAdapter *coreadapters.CoreAdapter) APIKeyControllerInterface {
	return &APIKeyController{coreAdapter: coreAdapter}
}

func (c *APIKeyController) CreateMyAPIKey(ctx *gin.Context, request *apicontract.CreateMyAPIKeyRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.CreateMyAPIKeyRequestDto, apicontract.CreateMyAPIKeyResponseDto](ctx, c.coreAdapter, request, apicontract.CreateMyAPIKeyOperation, "/core/v1/api-keys/create")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeCreatedClientResponse(ctx, response.Data)
}

func (c *APIKeyController) ListMyAPIKeys(ctx *gin.Context, request *apicontract.ListMyAPIKeysRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.ListMyAPIKeysRequestDto, apicontract.ListMyAPIKeysResponseDto](ctx, c.coreAdapter, request, apicontract.ListMyAPIKeysOperation, "/core/v1/api-keys/list")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *APIKeyController) RevokeMyAPIKey(ctx *gin.Context, request *apicontract.RevokeMyAPIKeyRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.RevokeMyAPIKeyRequestDto, apicontract.RevokeMyAPIKeyResponseDto](ctx, c.coreAdapter, request, apicontract.RevokeMyAPIKeyOperation, "/core/v1/api-keys/revoke")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}
