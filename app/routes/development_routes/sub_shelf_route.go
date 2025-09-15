package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
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
		middlewares.RateLimitMiddleware(1),
	)
	{
		subShelfRoutes.GET(
			"/getMySubShelfById",
			subShelfBinder.BindGetMySubShelfById(
				subShelfController.GetMySubShelfById,
			),
		)
		subShelfRoutes.GET(
			"/getAllSubShelvesByRootShelfId",
			subShelfBinder.BindGetAllSubShelvesByRootShelfId(
				subShelfController.GetAllSubShelvesByRootShelfId,
			),
		)
		subShelfRoutes.POST(
			"/createSubShelfByRootShelfId",
			subShelfBinder.BindCreateSubShelfByRootShelfId(
				subShelfController.CreateSubShelfByRootShelfId,
			),
		)
		subShelfRoutes.PUT(
			"/renameMySubShelfById",
			subShelfBinder.BindRenameMySubShelfById(
				subShelfController.RenameMySubShelfById,
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
