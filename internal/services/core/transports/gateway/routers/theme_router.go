package routers

import (
	"github.com/gin-gonic/gin"

	themesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/themes"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureThemeRoutes(router *gin.RouterGroup, endpoint endpoints.ThemeEndpointInterface) {
	themeRoutes := router.Group("/themes")
	{
		themeRoutes.POST(
			"/graphql/search",
			middlewares.DelegationMiddleware(
				themesdto.SearchThemesOperation,
			),
			endpoint.SearchThemes,
		)
	}
}
