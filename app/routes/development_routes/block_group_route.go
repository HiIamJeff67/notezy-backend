package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentBlockGroupRoutes() {
	blockGroupModule := modules.NewBlockGroupModule()

	blockGroupRoutes := DevelopmentRouterGroup.Group("/blockGroup")
	blockGroupRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		blockGroupRoutes.GET(
			"/getMyBlockGroupById",
			blockGroupModule.Binder.BindGetMyBlockGroupById(
				blockGroupModule.Controller.GetMyBlockGroupById,
			),
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupAndItsBlocksById",
			blockGroupModule.Binder.BindGetMyBlockGroupAndItsBlocksById(
				blockGroupModule.Controller.GetMyBlockGroupAndItsBlocksById,
			),
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupsAndTheirBlocksByBlockPackId",
			blockGroupModule.Binder.BindGetMyBlockGroupsAndTheirBlocksByBlockPackId(
				blockGroupModule.Controller.GetMyBlockGroupsAndTheirBlocksByBlockPackId,
			),
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupsByPrevBlockGroupId",
			blockGroupModule.Binder.BindGetMyBlockGroupsByPrevBlockGroupId(
				blockGroupModule.Controller.GetMyBlockGroupsByPrevBlockGroupId,
			),
		)
		blockGroupRoutes.GET(
			"/getAllMyBlockGroupsByBlockPackId",
			blockGroupModule.Binder.BindGetAllMyBlockGroupsByBlockPackId(
				blockGroupModule.Controller.GetAllMyBlockGroupsByBlockPackId,
			),
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupByBlockPackId",
			blockGroupModule.Binder.BindInsertBlockGroupByBlockPackId(
				blockGroupModule.Controller.InsertBlockGroupByBlockPackId,
			),
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupAndItsBlocksByBlockPackId",
			blockGroupModule.Binder.BindInsertBlockGroupAndItsBlocksByBlockPackId(
				blockGroupModule.Controller.InsertBlockGroupAndItsBlocksByBlockPackId,
			),
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupsAndTheirBlocksByBlockPackId",
			blockGroupModule.Binder.BindInsertBlockGroupsAndTheirBlocksByBlockPackId(
				blockGroupModule.Controller.InsertBlockGroupsAndTheirBlocksByBlockPackId,
			),
		)
		blockGroupRoutes.POST(
			"/insertSequentialBlockGroupsAndTheirBlocksByBlockPackId",
			blockGroupModule.Binder.BindInsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(
				blockGroupModule.Controller.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId,
			),
		)
		blockGroupRoutes.PUT(
			"/moveMyBlockGroupsByIds",
			blockGroupModule.Binder.BindMoveMyBlockGroupsByIds(
				blockGroupModule.Controller.MoveMyBlockGroupsByIds,
			),
		)
		blockGroupRoutes.PATCH(
			"/restoreMyBlockGroupById",
			blockGroupModule.Binder.BindRestoreMyBlockGroupById(
				blockGroupModule.Controller.RestoreMyBlockGroupById,
			),
		)
		blockGroupRoutes.PATCH(
			"/restoreMyBlockGroupsByIds",
			blockGroupModule.Binder.BindRestoreMyBlockGroupsByIds(
				blockGroupModule.Controller.RestoreMyBlockGroupsByIds,
			),
		)
		blockGroupRoutes.DELETE(
			"/deleteMyBlockGroupById",
			blockGroupModule.Binder.BindDeleteMyBlockGroupById(
				blockGroupModule.Controller.DeleteMyBlockGroupById,
			),
		)
		blockGroupRoutes.DELETE(
			"/deleteMyBlockGroupsByIds",
			blockGroupModule.Binder.BindDeleteMyBlockGroupsByIds(
				blockGroupModule.Controller.DeleteMyBlockGroupsByIds,
			),
		)
	}
}
