package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/themes"

	otherservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/other"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type ThemeRouterDependencies struct {
	Service otherservices.ThemeServiceInterface
}

func configureThemeRoutes(router *gin.RouterGroup, deps ThemeRouterDependencies) {
	endpoint := endpoints.NewThemeEndpoint(deps.Service)
	themeRoutes := router.Group("/themes")
	{
		themeRoutes.POST(
			"/graphql/search",
			middlewares.DelegationMiddleware(
				apicontract.SearchThemesOperation,
			),
			endpoint.SearchThemes,
		)
	}
}
