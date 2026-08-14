package routers

import (
	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata"
	authservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/auth"
	blockservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/blocks"
	materialservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/material"
	otherservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/other"
	realtimeservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/realtime"
	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
	shelfservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/shelves"
	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func NewRouter(
	authMiddleware gin.HandlerFunc,
	apiKeyMiddleware gin.HandlerFunc,
	authService authservices.AuthServiceInterface,
	rootShelfService shelfservices.RootShelfServiceInterface,
	stationService routineservices.StationServiceInterface,
	userSettingService userservices.UserSettingServiceInterface,
	userInfoService userservices.UserInfoServiceInterface,
	userAccountService userservices.UserAccountServiceInterface,
	userService userservices.UserServiceInterface,
	blockService blockservices.BlockServiceInterface,
	realtimeService realtimeservices.RealtimeServiceInterface,
	routineTagService routineservices.RoutineTagServiceInterface,
	routineTaskRecordService routineservices.RoutineTaskRecordServiceInterface,
	subShelfService shelfservices.SubShelfServiceInterface,
	blockPackService blockservices.BlockPackServiceInterface,
	materialService materialservices.MaterialServiceInterface,
	routineService routineservices.RoutineServiceInterface,
	routineTaskService routineservices.RoutineTaskServiceInterface,
	themeService otherservices.ThemeServiceInterface,
	itemService shelfservices.ItemServiceInterface,
	badgeService otherservices.BadgeServiceInterface,
	userDataCacheClient *userdata.UserDataCacheClient,
) *gin.Engine {
	router := gin.New()
	router.Use(middlewares.TokenResponseMiddleware())

	coreRouterGroup := router.Group("/core/" + gatewaycontract.Version)
	anonymousCoreRouterGroup := coreRouterGroup.Group("")
	secureCoreRouterGroup := coreRouterGroup.Group("")
	clientAuthMiddleware := authMiddleware
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	authEndpoint := endpoints.NewAuthEndpoint(authService)
	configureAnonymousAuthRoutes(anonymousCoreRouterGroup, authEndpoint)
	configureAuthenticatedAuthRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		authEndpoint,
		userDataCacheClient,
	)

	rootShelfEndpoint := endpoints.NewRootShelfEndpoint(rootShelfService)
	configureRootShelfRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		rootShelfEndpoint,
	)
	stationEndpoint := endpoints.NewStationEndpoint(stationService)
	configureStationRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		stationEndpoint,
	)
	userSettingEndpoint := endpoints.NewUserSettingEndpoint(userSettingService)
	configureUserSettingRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		userSettingEndpoint,
	)
	userInfoEndpoint := endpoints.NewUserInfoEndpoint(userInfoService)
	configureUserInfoRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		userInfoEndpoint,
	)
	userAccountEndpoint := endpoints.NewUserAccountEndpoint(userAccountService)
	configureUserAccountRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		userAccountEndpoint,
	)
	userEndpoint := endpoints.NewUserEndpoint(userService)
	configureUserRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		userEndpoint,
	)
	blockEndpoint := endpoints.NewBlockEndpoint(blockService)
	configureBlockRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		blockEndpoint,
	)
	realtimeEndpoint := endpoints.NewRealtimeEndpoint(realtimeService)
	configureRealtimeRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		realtimeEndpoint,
	)
	routineTagEndpoint := endpoints.NewRoutineTagEndpoint(routineTagService)
	configureRoutineTagRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		routineTagEndpoint,
	)
	routineTaskRecordEndpoint := endpoints.NewRoutineTaskRecordEndpoint(routineTaskRecordService)
	configureRoutineTaskRecordRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		routineTaskRecordEndpoint,
	)
	subShelfEndpoint := endpoints.NewSubShelfEndpoint(subShelfService)
	configureSubShelfRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		subShelfEndpoint,
	)
	blockPackEndpoint := endpoints.NewBlockPackEndpoint(blockPackService)
	configureBlockPackRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		blockPackEndpoint,
	)
	materialEndpoint := endpoints.NewMaterialEndpoint(materialService)
	configureMaterialRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		materialEndpoint,
	)
	routineEndpoint := endpoints.NewRoutineEndpoint(routineService)
	configureRoutineRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		routineEndpoint,
	)
	routineTaskEndpoint := endpoints.NewRoutineTaskEndpoint(routineTaskService)
	configureRoutineTaskRoutes(
		secureCoreRouterGroup,
		apiCompatibleAuthMiddleware,
		routineTaskEndpoint,
	)
	themeEndpoint := endpoints.NewThemeEndpoint(themeService)
	configureThemeRoutes(anonymousCoreRouterGroup, themeEndpoint)
	itemEndpoint := endpoints.NewItemEndpoint(itemService)
	configureItemRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		itemEndpoint,
	)
	badgeEndpoint := endpoints.NewBadgeEndpoint(badgeService)
	configureBadgeRoutes(
		secureCoreRouterGroup,
		clientAuthMiddleware,
		badgeEndpoint,
	)

	return router
}
