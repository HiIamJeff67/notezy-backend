package developmentroutes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	storages "github.com/HiIamJeff67/notezy-backend/internal/shared/storage"
)

func configureStorageRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	storageRoute := router.Group("/storage")
	storageRoute.Use(
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(5*time.Second),
	)
	{
		// only on test environment
		storageRoute.GET(
			"/mock/files/:presignedURL",
			func(ctx *gin.Context) {
				// technically, we use the presigned url as the key in in memory storage
				// since it is only for testing purposes
				key := ctx.Param("presignedURL")
				rc, object, err := storages.InMemoryStorage.GetObjectByKey(ctx, key, nil)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found."})
					return
				}
				defer rc.Close()
				logs.NotezyLogger.Info(ctx.Request.Context(), "Successfully get the files!")
				ctx.Data(http.StatusOK, object.ContentType, object.Data)
			},
		)
		// only on test environment
		storageRoute.GET(
			"/all",
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"keys": storages.InMemoryStorage.ListKeys(),
				})
			},
		)
	}
}
