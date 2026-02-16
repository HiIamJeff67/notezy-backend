package modules

import (
	"notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type BlockPackModule struct {
	Binder     binders.BlockPackBinderInterface
	Controller controllers.BlockPackControllerInterface
}

func NewBlockPackModule() *BlockPackModule {
	subShelfRepository := repositories.NewSubShelfRepository()
	blockPackRepository := repositories.NewBlockPackRepository()
	blockGroupRepository := repositories.NewBlockGroupRepository()
	blockRepository := repositories.NewBlockRepository()
	editableBlockAdapter := adapters.NewEditableBlockAdapter()

	blockPackService := services.NewBlockPackService(
		models.NotezyDB,
		subShelfRepository,
		blockPackRepository,
		blockGroupRepository,
		blockRepository,
		editableBlockAdapter,
	)

	blockPackBinder := binders.NewBlockPackBinder()

	blockPackController := controllers.NewBlockPackController(
		blockPackService,
	)

	return &BlockPackModule{
		Binder:     blockPackBinder,
		Controller: blockPackController,
	}
}
