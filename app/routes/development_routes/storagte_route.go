package developmentroutes

import (
	"net/http"
	"notezy-backend/app/storages"

	"github.com/gin-gonic/gin"
)

func configureStorageRoutes() {
	storageRoute := DevelopmentRouterGroup.Group("/storage")
	{
		// only on test envrionment
		storageRoute.GET(
			"/mock/files/:presignedURL",
			func(ctx *gin.Context) {
				// technically, we use the presigned url as the key in in memory storage
				// since it is only for testing purposes
				key := ctx.Param("presignedURL")
				rc, object, exception := storages.InMemoryStorage.GetObjectByKey(ctx, key, nil)
				if exception != nil {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "Filed not found."})
					return
				}
				defer rc.Close()
				ctx.Data(http.StatusOK, object.ContentType, object.Data)
			},
		)
	}
}
