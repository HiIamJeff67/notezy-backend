package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
	storages "notezy-backend/app/storages"
)

type SubShelfModule struct {
	Binder     binders.SubShelfBinderInterface
	Controller controllers.SubShelfControllerInterface
}

func NewSubShelfModule() *SubShelfModule {
	subShelfRepository := repositories.NewSubShelfRepository()
	rootShelfRepository := repositories.NewRootShelfRepository()
	materialRepository := repositories.NewMaterialRepository()
	blockPackRepository := repositories.NewBlockPackRepository()

	subShelfService := services.NewSubShelfService(
		models.NotezyDB,
		storages.InMemoryStorage,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
	)

	subShelfBinder := binders.NewSubShelfBinder()

	subShelfController := controllers.NewSubShelfController(
		subShelfService,
	)

	return &SubShelfModule{
		Binder:     subShelfBinder,
		Controller: subShelfController,
	}
}
