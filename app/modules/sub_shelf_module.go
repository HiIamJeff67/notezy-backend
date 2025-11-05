package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type SubShelfModule struct {
	Binder     binders.SubShelfBinderInterface
	Controller controllers.SubShelfControllerInterface
}

func NewSubShelfModule() *SubShelfModule {
	subShelfRepository := repositories.NewSubShelfRepository()

	subShelfService := services.NewSubShelfService(
		models.NotezyDB,
		subShelfRepository,
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
