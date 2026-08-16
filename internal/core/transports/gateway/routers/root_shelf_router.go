package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/root-shelves"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	shelfservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/shelves"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type RootShelfRouterDependencies struct {
	Service          shelfservices.RootShelfServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureRootShelfRoutes(
	router *gin.RouterGroup,
	deps RootShelfRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewRootShelfEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	rootShelfRoutes := router.Group("/root-shelves")
	{
		rootShelfRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRootShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyRootShelfById,
		)
		rootShelfRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRootShelfOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateRootShelf,
		)
		rootShelfRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRootShelvesOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateRootShelves,
		)
		rootShelfRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRootShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRootShelfById,
		)
		rootShelfRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRootShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRootShelvesByIds,
		)
		rootShelfRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyRootShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyRootShelfById,
		)
		rootShelfRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyRootShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyRootShelvesByIds,
		)
		rootShelfRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRootShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRootShelfById,
		)
		rootShelfRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRootShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRootShelvesByIds,
		)
		rootShelfRoutes.POST(
			"/permissions/get",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRootShelfPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyRootShelfPermission,
		)
		rootShelfRoutes.POST(
			"/permissions/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyRootShelfPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateMyRootShelfPermission,
		)
		rootShelfRoutes.POST(
			"/permissions/upsert",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpsertMyRootShelfPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpsertMyRootShelfPermission,
		)
		rootShelfRoutes.POST(
			"/permissions/upsert-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpsertMyRootShelfPermissionsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpsertMyRootShelfPermissions,
		)
		rootShelfRoutes.POST(
			"/permissions/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRootShelfPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRootShelfPermission,
		)
		rootShelfRoutes.POST(
			"/ownership/transfer",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.TransferMyRootShelfOwnershipOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.TransferMyRootShelfOwnership,
		)
		rootShelfRoutes.POST(
			"/permissions/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRootShelfPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRootShelfPermission,
		)
		rootShelfRoutes.POST(
			"/permissions/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRootShelfPermissionsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRootShelfPermissions,
		)
		rootShelfRoutes.POST(
			"/memberships/leave",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LeaveMyRootShelfOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LeaveMyRootShelf,
		)
		rootShelfRoutes.POST(
			"/memberships/leave-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LeaveMyRootShelvesOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LeaveMyRootShelves,
		)
		rootShelfRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchRootShelvesOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchRootShelves,
		)
	}
}
