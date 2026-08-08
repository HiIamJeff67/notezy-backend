package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
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
				apicontract.GetMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockPackAndItsParentByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPackAndItsParentById,
		)
		blockPackRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockPacksByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyBlockPacksByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyBlockPacksByRootShelfId,
		)
		blockPackRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateBlockPackOperation,
			),
			authMiddleware,
			endpoint.CreateBlockPack,
		)
		blockPackRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateBlockPacksOperation,
			),
			authMiddleware,
			endpoint.CreateBlockPacks,
		)
		blockPackRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPackByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPackByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPacksByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many-by-parent-sub-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPacksByParentSubShelfIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfIds,
		)
		blockPackRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyBlockPackByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyBlockPacksByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchBlockPacksOperation,
			),
			authMiddleware,
			endpoint.SearchBlockPacks,
		)
	}
}
