package modules

import (
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
	blockGroupRepository := repositories.NewBlockGroupRepository()
	blockRepository := repositories.NewBlockRepository()

	blockService := services.NewBlockService(
		models.NotezyDB,
		blockGroupRepository,
		blockRepository,
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
