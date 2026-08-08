package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/badges"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureBadgeRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.BadgeEndpointInterface,
) {
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
