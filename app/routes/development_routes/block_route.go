package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentBlockRoutes() {
	blockModule := modules.NewBlockModule()

	blockRoutes := DevelopmentRouterGroup.Group("/block")
	blockRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		blockRoutes.GET(
			"/getMyBlockById",
			blockModule.Binder.BindGetMyBlockById(
				blockModule.Controller.GetMyBlockById,
			),
		)
		blockRoutes.GET(
			"/getMyBlocksByIds",
			blockModule.Binder.BindGetMyBlocksByIds(
				blockModule.Controller.GetMyBlocksByIds,
			),
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockGroupId",
			blockModule.Binder.BindGetMyBlocksByBlockGroupId(
				blockModule.Controller.GetMyBlocksByBlockGroupId,
			),
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockGroupIds",
			blockModule.Binder.BindGetMyBlocksByBlockGroupIds(
				blockModule.Controller.GetMyBlocksByBlockGroupIds,
			),
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockPackId",
			blockModule.Binder.BindGetMyBlocksByBlockPackId(
				blockModule.Controller.GetMyBlocksByBlockPackId,
			),
		)
		blockRoutes.GET(
			"/getAllMyBlocks",
			blockModule.Binder.BindGetAllMyBlocks(
				blockModule.Controller.GetAllMyBlocks,
			),
		)
		blockRoutes.POST(
			"/insertBlock",
			blockModule.Binder.BindInsertBlock(
				blockModule.Controller.InsertBlock,
			),
		)
		blockRoutes.POST(
			"/insertBlocks",
			blockModule.Binder.BindInsertBlocks(
				blockModule.Controller.InsertBlocks,
			),
		)
		blockRoutes.PUT(
			"/updateMyBlockById",
			blockModule.Binder.BindUpdateMyBlockById(
				blockModule.Controller.UpdateMyBlockById,
			),
		)
		blockRoutes.PUT(
			"/updateMyBlocksByIds",
			blockModule.Binder.BindUpdateMyBlocksByIds(
				blockModule.Controller.UpdateMyBlocksByIds,
			),
		)
		blockRoutes.PATCH(
			"/restoreMyBlockById",
			blockModule.Binder.BindRestoreMyBlockById(
				blockModule.Controller.RestoreMyBlockById,
			),
		)
		blockRoutes.PATCH(
			"/restoreMyBlocksByIds",
			blockModule.Binder.BindRestoreMyBlocksByIds(
				blockModule.Controller.RestoreMyBlocksByIds,
			),
		)
		blockRoutes.DELETE(
			"/deleteMyBlockById",
			blockModule.Binder.BindDeleteMyBlockById(
				blockModule.Controller.DeleteMyBlockById,
			),
		)
		blockRoutes.DELETE(
			"/deleteMyBlocksByIds",
			blockModule.Binder.BindDeleteMyBlocksByIds(
				blockModule.Controller.DeleteMyBlocksByIds,
			),
		)
	}
}
