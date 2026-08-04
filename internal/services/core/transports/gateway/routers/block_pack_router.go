package routers

import (
	"github.com/gin-gonic/gin"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureBlockPackRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.BlockPackEndpointInterface,
) {
	blockPackRoutes := router.Group("/block-packs")
	{
		blockPackRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.GetMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.GetMyBlockPackAndItsParentByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPackAndItsParentById,
		)
		blockPackRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.GetMyBlockPacksByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.GetAllMyBlockPacksByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyBlockPacksByRootShelfId,
		)
		blockPackRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.CreateBlockPackOperation,
			),
			authMiddleware,
			endpoint.CreateBlockPack,
		)
		blockPackRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.CreateBlockPacksOperation,
			),
			authMiddleware,
			endpoint.CreateBlockPacks,
		)
		blockPackRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.UpdateMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.UpdateMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.MoveMyBlockPackByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPackByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.MoveMyBlockPacksByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many-by-parent-sub-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfIds,
		)
		blockPackRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.RestoreMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.RestoreMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.DeleteMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.DeleteMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				blockpacksdto.SearchBlockPacksOperation,
			),
			authMiddleware,
			endpoint.SearchBlockPacks,
		)
	}
}
