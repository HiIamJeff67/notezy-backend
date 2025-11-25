package modules

import (
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

	blockPackService := services.NewBlockPackService(
		models.NotezyDB,
		subShelfRepository,
		blockPackRepository,
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
