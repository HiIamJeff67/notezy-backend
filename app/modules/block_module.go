package modules

import (
	adapters "notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
)

type BlockModule struct {
	Binder     binders.BlockBinderInterface
	Controller controllers.BlockControllerInterface
}

func NewBlockModule() *BlockModule {
	blockGroupScope := scopes.NewBlockGroupScope()
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())
	blockGroupRepository := repositories.NewBlockGroupRepository(blockGroupScope)
	blockRepository := repositories.NewBlockRepository(scopes.NewBlockScope())
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
