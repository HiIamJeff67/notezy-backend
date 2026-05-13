package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
	storages "notezy-backend/app/storages"
)

type SubShelfModule struct {
	Binder     binders.SubShelfBinderInterface
	Controller controllers.SubShelfControllerInterface
}

func NewSubShelfModule() *SubShelfModule {
	subShelfRepository := repositories.NewSubShelfRepository(scopes.NewSubShelfScope())
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())

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
