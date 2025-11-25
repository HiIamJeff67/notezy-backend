package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentBlockPackRoutes() {
	blockPackModule := modules.NewBlockPackModule()

	blockPackRoutes := DevelopmentRouterGroup.Group("/blockPack")
	blockPackRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		blockPackRoutes.GET(
			"/getMyBlockPackById",
			blockPackModule.Binder.BindGetMyBlockPackById(
				blockPackModule.Controller.GetMyBlockPackById,
			),
		)
		blockPackRoutes.GET(
			"/getAllMyBlockPacksByParentSubShelfId",
			blockPackModule.Binder.BindGetAllMyBlockPacksByParentSubShelfId(
				blockPackModule.Controller.GetAllMyBlockPacksByParentSubShelfId,
			),
		)
		blockPackRoutes.GET(
			"/getAllMyBlockPacksByRootShelfId",
			blockPackModule.Binder.BindGetAllMyBlockPacksByRootShelfId(
				blockPackModule.Controller.GetAllMyBlockPacksByRootShelfId,
			),
		)
		blockPackRoutes.POST(
			"/createBlockPack",
			blockPackModule.Binder.BindCreateBlockPack(
				blockPackModule.Controller.CreateBlockPack,
			),
		)
		blockPackRoutes.PUT(
			"/updateMyBlockPackById",
			blockPackModule.Binder.BindUpdateMyBlockPackById(
				blockPackModule.Controller.UpdateMyBlockPackById,
			),
		)
		blockPackRoutes.PUT(
			"/moveMyBlockPackById",
			blockPackModule.Binder.BindMoveMyBlockPackById(
				blockPackModule.Controller.MoveMyBlockPackById,
			),
		)
		blockPackRoutes.PUT(
			"/moveMyBlockPacksByIds",
			blockPackModule.Binder.BindMoveMyBlockPacksByIds(
				blockPackModule.Controller.MoveMyBlockPacksByIds,
			),
		)
		blockPackRoutes.PATCH(
			"/restoreMyBlockPackById",
			blockPackModule.Binder.BindRestoreMyBlockPackById(
				blockPackModule.Controller.RestoreMyBlockPackById,
			),
		)
		blockPackRoutes.PATCH(
			"/restoreMyBlockPacksByIds",
			blockPackModule.Binder.BindRestoreMyBlockPacksByIds(
				blockPackModule.Controller.RestoreMyBlockPacksByIds,
			),
		)
		blockPackRoutes.DELETE(
			"/deleteMyBlockPackById",
			blockPackModule.Binder.BindDeleteMyBlockPackById(
				blockPackModule.Controller.DeleteMyBlockPackById,
			),
		)
		blockPackRoutes.DELETE(
			"/deleteMyBlockPacksByIds",
			blockPackModule.Binder.BindDeleteMyBlockPacksByIds(
				blockPackModule.Controller.DeleteMyBlockPacksByIds,
			),
		)
	}
}
