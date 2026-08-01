package routers

import (
	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
)

func NewRouter(
	authMiddleware gin.HandlerFunc,
	authService services.AuthServiceInterface,
	rootShelfService services.RootShelfServiceInterface,
	stationService services.StationServiceInterface,
	userSettingService services.UserSettingServiceInterface,
	userInfoService services.UserInfoServiceInterface,
	userAccountService services.UserAccountServiceInterface,
	userService services.UserServiceInterface,
	blockService services.BlockServiceInterface,
	realtimeService services.RealtimeServiceInterface,
	routineTagService services.RoutineTagServiceInterface,
	routineTaskRecordService services.RoutineTaskRecordServiceInterface,
	subShelfService services.SubShelfServiceInterface,
	blockPackService services.BlockPackServiceInterface,
	materialService services.MaterialServiceInterface,
	routineService services.RoutineServiceInterface,
	routineTaskService services.RoutineTaskServiceInterface,
	themeService services.ThemeServiceInterface,
	itemService services.ItemServiceInterface,
	badgeService services.BadgeServiceInterface,
) *gin.Engine {
	router := gin.New()

	coreRouterGroup := router.Group("/core/" + core.Version)
	anonymousCoreRouterGroup := coreRouterGroup.Group("")
	secureCoreRouterGroup := coreRouterGroup.Group("")

	authEndpoint := endpoints.NewAuthEndpoint(authService)
	configureAnonymousAuthRoutes(anonymousCoreRouterGroup, authEndpoint)
	configureAuthenticatedAuthRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		authEndpoint,
	)

	rootShelfEndpoint := endpoints.NewRootShelfEndpoint(rootShelfService)
	configureRootShelfRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		rootShelfEndpoint,
	)
	stationEndpoint := endpoints.NewStationEndpoint(stationService)
	configureStationRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		stationEndpoint,
	)
	userSettingEndpoint := endpoints.NewUserSettingEndpoint(userSettingService)
	configureUserSettingRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		userSettingEndpoint,
	)
	userInfoEndpoint := endpoints.NewUserInfoEndpoint(userInfoService)
	configureUserInfoRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		userInfoEndpoint,
	)
	userAccountEndpoint := endpoints.NewUserAccountEndpoint(userAccountService)
	configureUserAccountRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		userAccountEndpoint,
	)
	userEndpoint := endpoints.NewUserEndpoint(userService)
	configureUserRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		userEndpoint,
	)
	blockEndpoint := endpoints.NewBlockEndpoint(blockService)
	configureBlockRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		blockEndpoint,
	)
	realtimeEndpoint := endpoints.NewRealtimeEndpoint(realtimeService)
	configureRealtimeRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		realtimeEndpoint,
	)
	routineTagEndpoint := endpoints.NewRoutineTagEndpoint(routineTagService)
	configureRoutineTagRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		routineTagEndpoint,
	)
	routineTaskRecordEndpoint := endpoints.NewRoutineTaskRecordEndpoint(routineTaskRecordService)
	configureRoutineTaskRecordRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		routineTaskRecordEndpoint,
	)
	subShelfEndpoint := endpoints.NewSubShelfEndpoint(subShelfService)
	configureSubShelfRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		subShelfEndpoint,
	)
	blockPackEndpoint := endpoints.NewBlockPackEndpoint(blockPackService)
	configureBlockPackRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		blockPackEndpoint,
	)
	materialEndpoint := endpoints.NewMaterialEndpoint(materialService)
	configureMaterialRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		materialEndpoint,
	)
	routineEndpoint := endpoints.NewRoutineEndpoint(routineService)
	configureRoutineRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		routineEndpoint,
	)
	routineTaskEndpoint := endpoints.NewRoutineTaskEndpoint(routineTaskService)
	configureRoutineTaskRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		routineTaskEndpoint,
	)
	themeEndpoint := endpoints.NewThemeEndpoint(themeService)
	configureThemeRoutes(anonymousCoreRouterGroup, themeEndpoint)
	itemEndpoint := endpoints.NewItemEndpoint(itemService)
	configureItemRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		itemEndpoint,
	)
	badgeEndpoint := endpoints.NewBadgeEndpoint(badgeService)
	configureBadgeRoutes(
		secureCoreRouterGroup,
		authMiddleware,
		badgeEndpoint,
	)

	return router
}
