package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
)

type RoutineTaskModule struct {
	Binder     binders.RoutineTaskBinderInterface
	Controller controllers.RoutineTaskControllerInterface
}

func NewRoutineTaskModule() *RoutineTaskModule {
	routineTaskRepository := repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope())

	routineTaskService := services.NewRoutineTaskService(
		models.NotezyDB,
		routineTaskRepository,
	)

	routineTaskBinder := binders.NewRoutineTaskBinder()

	routineTaskController := controllers.NewRoutineTaskController(
		routineTaskService,
	)

	return &RoutineTaskModule{
		Binder:     routineTaskBinder,
		Controller: routineTaskController,
	}
}
