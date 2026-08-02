package routers

import (
	"github.com/gin-gonic/gin"

	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

func ConfigureRoutes(
	router *gin.Engine,
	realtimeService services.RealtimeServiceInterface,
	yjsPersistenceService services.YjsPersistenceServiceInterface,
	blockService services.BlockServiceInterface,
) {
	websocketRoutes := router.Group("/websocket")
	{
		ConfigureBlockPackRoutes(websocketRoutes, realtimeService)
		ConfigureYjsRoutes(websocketRoutes, yjsPersistenceService)
		ConfigureBlockProjectionRoutes(websocketRoutes, blockService)
	}
}
