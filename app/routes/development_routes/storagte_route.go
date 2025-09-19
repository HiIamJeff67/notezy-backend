package developmentroutes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	logs "notezy-backend/app/logs"
	storages "notezy-backend/app/storages"
)

func configureStorageRoutes() {
	storageRoute := DevelopmentRouterGroup.Group("/storage")
	{
		// only on test environment
		storageRoute.GET(
			"/mock/files/:presignedURL",
			func(ctx *gin.Context) {
				// technically, we use the presigned url as the key in in memory storage
				// since it is only for testing purposes
				key := ctx.Param("presignedURL")
				rc, object, exception := storages.InMemoryStorage.GetObjectByKey(ctx, key, nil)
				if exception != nil {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found."})
					return
				}
				defer rc.Close()
				logs.Info("Successfully get the files!")
				logs.Info("Details: ", object)
				ctx.Data(http.StatusOK, object.ContentType, object.Data)
			},
		)
		// only on test environment
		storageRoute.GET(
			"/listAllInTerminal",
			func(ctx *gin.Context) {
				storages.InMemoryStorage.ListAllInTerminal()
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		)
	}
}
