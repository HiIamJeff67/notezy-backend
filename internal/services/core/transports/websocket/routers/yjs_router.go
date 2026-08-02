package routers

import (
	"github.com/gin-gonic/gin"

	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/websocket/endpoints"
)

func ConfigureYjsRoutes(
	router *gin.RouterGroup,
	yjsPersistenceService services.YjsPersistenceServiceInterface,
) {
	endpoint := endpoints.NewYjsEndpoint(yjsPersistenceService)

	router.POST(
		"/yjs-document-compaction/load",
		middlewares.DelegationMiddleware(
			websocketdto.LoadCompactableYjsDocumentOperation,
		),
		endpoint.LoadCompactableYjsDocument,
	)
	router.POST(
		"/yjs-document-compaction/apply",
		middlewares.DelegationMiddleware(
			websocketdto.ApplyCompactedYjsDocumentOperation,
		),
		endpoint.ApplyCompactedYjsDocument,
	)
	router.POST(
		"/yjs-document/load",
		middlewares.DelegationMiddleware(
			websocketdto.LoadYjsDocumentOperation,
		),
		endpoint.LoadYjsDocument,
	)
	router.POST(
		"/yjs-update/append",
		middlewares.DelegationMiddleware(
			websocketdto.AppendYjsUpdateOperation,
		),
		endpoint.AppendYjsUpdate,
	)
}
