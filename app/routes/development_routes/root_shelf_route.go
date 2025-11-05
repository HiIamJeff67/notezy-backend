package developmentroutes

import (
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentRootShelfRoutes() {
	rootShelfModule := modules.NewRootShelfModule()

	rootShelfRoutes := DevelopmentRouterGroup.Group("/rootShelf")
	rootShelfRoutes.Use(
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		rootShelfRoutes.GET(
			"/getMyRootShelfById",
			rootShelfModule.Binder.BindGetMyRootShelfById(
				rootShelfModule.Controller.GetMyRootShelfById,
			),
		)
		rootShelfRoutes.GET(
			"/searchRecentRootShelves",
			rootShelfModule.Binder.BindSearchRecentRootShelves(
				rootShelfModule.Controller.SearchRecentRootShelves,
			),
		)
		rootShelfRoutes.POST(
			"/createRootShelf",
			rootShelfModule.Binder.BindCreateRootShelf(
				rootShelfModule.Controller.CreateRootShelf,
			),
		)
		rootShelfRoutes.PUT(
			"/updateMyRootShelfById",
			rootShelfModule.Binder.BindUpdateMyRootShelfById(
				rootShelfModule.Controller.UpdateMyRootShelfById,
			),
		)
		rootShelfRoutes.PATCH(
			"/restoreMyRootShelfById",
			rootShelfModule.Binder.BindRestoreMyRootShelfById(
				rootShelfModule.Controller.RestoreMyRootShelfById,
			),
		)
		rootShelfRoutes.PATCH(
			"/restoreMyRootShelvesByIds",
			rootShelfModule.Binder.BindRestoreMyRootShelvesByIds(
				rootShelfModule.Controller.RestoreMyRootShelvesByIds,
			),
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelfById",
			rootShelfModule.Binder.BindDeleteMyRootShelfById(
				rootShelfModule.Controller.DeleteMyRootShelfById,
			),
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelvesByIds",
			rootShelfModule.Binder.BindDeleteMyRootShelvesByIds(
				rootShelfModule.Controller.DeleteMyRootShelvesByIds,
			),
		)
	}
}
