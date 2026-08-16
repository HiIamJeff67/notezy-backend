package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/stations"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	routineservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/routines"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type StationRouterDependencies struct {
	Service          routineservices.StationServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureStationRoutes(
	router *gin.RouterGroup,
	deps StationRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewStationEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	stationRoutes := router.Group("/stations")
	{
		stationRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyStationByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyStationById,
		)
		stationRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyStationsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyStations,
		)
		stationRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateStationOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateStation,
		)
		stationRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateStationsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateStations,
		)
		stationRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyStationByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyStationById,
		)
		stationRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyStationsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyStationsByIds,
		)
		stationRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyStationByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyStationById,
		)
		stationRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyStationsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyStationsByIds,
		)
		stationRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyStationByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyStationById,
		)
		stationRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyStationsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyStationsByIds,
		)
		stationRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyStationByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.HardDeleteMyStationById,
		)
		stationRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyStationsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.HardDeleteMyStationsByIds,
		)
		stationRoutes.POST(
			"/visualizations/total-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyTotalCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyTotalCount,
		)
		stationRoutes.POST(
			"/permissions/get",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyStationPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyStationPermission,
		)
		stationRoutes.POST(
			"/permissions/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyStationPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateMyStationPermission,
		)
		stationRoutes.POST(
			"/permissions/upsert",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpsertMyStationPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpsertMyStationPermission,
		)
		stationRoutes.POST(
			"/permissions/upsert-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpsertMyStationPermissionsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpsertMyStationPermissions,
		)
		stationRoutes.POST(
			"/permissions/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyStationPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyStationPermission,
		)
		stationRoutes.POST(
			"/ownership/transfer",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.TransferMyStationOwnershipOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.TransferMyStationOwnership,
		)
		stationRoutes.POST(
			"/permissions/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyStationPermissionOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyStationPermission,
		)
		stationRoutes.POST(
			"/permissions/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyStationPermissionsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyStationPermissions,
		)
		stationRoutes.POST(
			"/memberships/leave",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LeaveMyStationOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LeaveMyStation,
		)
		stationRoutes.POST(
			"/memberships/leave-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LeaveMyStationsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LeaveMyStations,
		)
		stationRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchStationsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchStations,
		)
	}
}
