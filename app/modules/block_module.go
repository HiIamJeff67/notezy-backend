package modules

import (
	adapters "notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type BlockModule struct {
	Binder     binders.BlockBinderInterface
	Controller controllers.BlockControllerInterface
}

func NewBlockModule() *BlockModule {
	blockPackRepository := repositories.NewBlockPackRepository()
	blockGroupRepository := repositories.NewBlockGroupRepository()
	blockRepository := repositories.NewBlockRepository()
	editableBlockAdapter := adapters.NewEditableBlockAdapter()

	blockService := services.NewBlockService(
		models.NotezyDB,
		blockPackRepository,
		blockGroupRepository,
		blockRepository,
		editableBlockAdapter,
	)

	blockBinder := binders.NewBlockBinder()

	blockController := controllers.NewBlockController(
		blockService,
	)

	return &BlockModule{
		Binder:     blockBinder,
		Controller: blockController,
	}
}
