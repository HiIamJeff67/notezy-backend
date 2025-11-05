package developmentroutes

import (
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentSubShelfRoutes() {
	subShelfModule := modules.NewSubShelfModule()

	subShelfRoutes := DevelopmentRouterGroup.Group("/subShelf")
	subShelfRoutes.Use(
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		subShelfRoutes.GET(
			"/getMySubShelfById",
			subShelfModule.Binder.BindGetMySubShelfById(
				subShelfModule.Controller.GetMySubShelfById,
			),
		)
		subShelfRoutes.GET(
			"/getMySubShelvesByPrevSubShelfId",
			subShelfModule.Binder.BindGetMySubShelvesByPrevSubShelfId(
				subShelfModule.Controller.GetMySubShelvesByPrevSubShelfId,
			),
		)
		subShelfRoutes.GET(
			"/getAllMySubShelvesByRootShelfId",
			subShelfModule.Binder.BindGetAllMySubShelvesByRootShelfId(
				subShelfModule.Controller.GetAllMySubShelvesByRootShelfId,
			),
		)
		subShelfRoutes.POST(
			"/createSubShelfByRootShelfId",
			subShelfModule.Binder.BindCreateSubShelfByRootShelfId(
				subShelfModule.Controller.CreateSubShelfByRootShelfId,
			),
		)
		subShelfRoutes.PUT(
			"/updateMySubShelfById",
			subShelfModule.Binder.BindUpdateMySubShelfById(
				subShelfModule.Controller.UpdateMySubShelfById,
			),
		)
		subShelfRoutes.PUT(
			"/moveMySubShelf",
			subShelfModule.Binder.BindMoveMySubShelf(
				subShelfModule.Controller.MoveMySubShelf,
			),
		)
		subShelfRoutes.PUT(
			"/moveMySubShelves",
			subShelfModule.Binder.BindMoveMySubShelves(
				subShelfModule.Controller.MoveMySubShelves,
			),
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelfById",
			subShelfModule.Binder.BindRestoreMySubShelfById(
				subShelfModule.Controller.RestoreMySubShelfById,
			),
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelvesByIds",
			subShelfModule.Binder.BindRestoreMySubShelvesByIds(
				subShelfModule.Controller.RestoreMySubShelvesByIds,
			),
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelfById",
			subShelfModule.Binder.BindDeleteMySubShelfById(
				subShelfModule.Controller.DeleteMySubShelfById,
			),
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelvesByIds",
			subShelfModule.Binder.BindDeleteMySubShelvesByIds(
				subShelfModule.Controller.DeleteMySubShelvesByIds,
			),
		)
	}
}
