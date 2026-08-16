package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/materials"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	materialservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/material"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type MaterialRouterDependencies struct {
	Service          materialservices.MaterialServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureMaterialRoutes(
	router *gin.RouterGroup,
	deps MaterialRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewMaterialEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	materialRoutes := router.Group("/materials")
	{
		materialRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyMaterialById,
		)
		materialRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyMaterialAndItsParentByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyMaterialAndItsParentById,
		)
		materialRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyMaterialsByParentSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyMaterialsByParentSubShelfId,
		)
		materialRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyMaterialsByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyMaterialsByRootShelfId,
		)
		materialRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyMaterialOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateMyMaterial,
		)
		materialRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyMaterialById,
		)
		materialRoutes.POST(
			"/save",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SaveMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SaveMyMaterialById,
		)
		materialRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMyMaterialById,
		)
		materialRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyMaterialsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyMaterialById,
		)
		materialRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyMaterialsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyMaterialByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyMaterialById,
		)
		materialRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyMaterialsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchMaterialsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchMaterials,
		)
	}
}
