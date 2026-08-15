package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/badges"

	otherservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/other"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type BadgeRouterDependencies struct {
	Service        otherservices.BadgeServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureBadgeRoutes(
	router *gin.RouterGroup,
	deps BadgeRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewBadgeEndpoint(deps.Service)
	badgeRoutes := router.Group("/badges")
	{
		badgeRoutes.POST(
			"/graphql/load",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LoadUserBadgesOperation,
			),
			authMiddleware,
			endpoint.LoadUserBadges,
		)
	}
}
