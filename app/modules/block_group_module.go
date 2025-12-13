package modules

import (
	adapters "notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type BlockGroupModule struct {
	Binder     binders.BlockGroupBinderInterface
	Controller controllers.BlockGroupControllerInterface
}

func NewBlockGroupModule() *BlockGroupModule {
	blockGroupRepository := repositories.NewBlockGroupRepository()
	blockRepository := repositories.NewBlockRepository()
	editableBlockAdapter := adapters.NewEditableBlockAdapter()

	blockGroupService := services.NewBlockGroupService(
		models.NotezyDB,
		blockGroupRepository,
		blockRepository,
		editableBlockAdapter,
	)

	blockGroupBinder := binders.NewBlockGroupBinder()

	blockGroupController := controllers.NewBlockGroupController(
		blockGroupService,
	)

	return &BlockGroupModule{
		Binder:     blockGroupBinder,
		Controller: blockGroupController,
	}
}
