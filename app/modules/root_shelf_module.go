package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type RootShelfModule struct {
	Binder     binders.RootShelfBinderInterface
	Controller controllers.RootShelfControllerInterface
}

func NewRootShelfModule() *RootShelfModule {
	rootShelfRepository := repositories.NewRootShelfRepository()

	rootShelfService := services.NewRootShelfService(
		models.NotezyDB,
		rootShelfRepository,
	)

	rootShelfBinder := binders.NewRootShelfBinder()

	rootShelfController := controllers.NewRootShelfController(
		rootShelfService,
	)

	return &RootShelfModule{
		Binder:     rootShelfBinder,
		Controller: rootShelfController,
	}
}
