package modules

import (
	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type BlockModule struct {
	Binder     binders.BlockBinderInterface
	Controller controllers.BlockControllerInterface
}

func NewBlockModule() *BlockModule {
	blockScope := scopes.NewBlockScope()
	blockPackScope := scopes.NewBlockPackScope()
	subShelfScope := scopes.NewSubShelfScope()
	blockPackRepository := repositories.NewBlockPackRepository(blockPackScope)
	blockRepository := repositories.NewBlockRepository(blockScope)
	editableBlockAdapter := adapters.NewEditableBlockAdapter()

	blockService := services.NewBlockService(
		models.NotezyDB,
		blockScope,
		blockPackScope,
		subShelfScope,
		blockPackRepository,
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
