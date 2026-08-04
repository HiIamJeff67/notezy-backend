package routers

import (
	"github.com/gin-gonic/gin"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtimegateway/v1"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/gateway/middlewares"
)

func ConfigureBlockPackRoutes(
	router *gin.RouterGroup,
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient,
) {
	endpoint := endpoints.NewBlockPackEndpoint(realtimeLeaseCache)

	router.POST(
		"/block-pack-participants/get",
		middlewares.DelegationMiddleware(
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
		),
		endpoint.GetParticipants,
	)
}
