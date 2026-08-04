package routers

import (
	"github.com/gin-gonic/gin"

	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
)

func ConfigureRoutes(router *gin.Engine, realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient) {
	gatewayRoutes := router.Group("/gateway/v1")
	{
		ConfigureBlockPackRoutes(gatewayRoutes, realtimeLeaseCache)
	}
}
