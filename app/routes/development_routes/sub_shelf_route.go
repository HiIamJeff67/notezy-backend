package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentSubShelfRoutes() {
	subShelfBinder := binders.NewSubShelfBinder()
	subShelfController := controllers.NewSubShelfController(
		services.NewSubShelfService(
			models.NotezyDB,
		),
	)

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
			subShelfBinder.BindGetMySubShelfById(
				subShelfController.GetMySubShelfById,
			),
		)
		subShelfRoutes.GET(
			"/getMySubShelvesByPrevSubShelfId",
			subShelfBinder.BindGetMySubShelvesByPrevSubShelfId(
				subShelfController.GetMySubShelvesByPrevSubShelfId,
			),
		)
		subShelfRoutes.GET(
			"/getAllMySubShelvesByRootShelfId",
			subShelfBinder.BindGetAllMySubShelvesByRootShelfId(
				subShelfController.GetAllMySubShelvesByRootShelfId,
			),
		)
		subShelfRoutes.POST(
			"/createSubShelfByRootShelfId",
			subShelfBinder.BindCreateSubShelfByRootShelfId(
				subShelfController.CreateSubShelfByRootShelfId,
			),
		)
		subShelfRoutes.PUT(
			"/updateMySubShelfById",
			subShelfBinder.BindUpdateMySubShelfById(
				subShelfController.UpdateMySubShelfById,
			),
		)
		subShelfRoutes.PUT(
			"/moveMySubShelf",
			subShelfBinder.BindMoveMySubShelf(
				subShelfController.MoveMySubShelf,
			),
		)
		subShelfRoutes.PUT(
			"/moveMySubShelves",
			subShelfBinder.BindMoveMySubShelves(
				subShelfController.MoveMySubShelves,
			),
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelfById",
			subShelfBinder.BindRestoreMySubShelfById(
				subShelfController.RestoreMySubShelfById,
			),
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelvesByIds",
			subShelfBinder.BindRestoreMySubShelvesByIds(
				subShelfController.RestoreMySubShelvesByIds,
			),
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelfById",
			subShelfBinder.BindDeleteMySubShelfById(
				subShelfController.DeleteMySubShelfById,
			),
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelvesByIds",
			subShelfBinder.BindDeleteMySubShelvesByIds(
				subShelfController.DeleteMySubShelvesByIds,
			),
		)
	}
}
