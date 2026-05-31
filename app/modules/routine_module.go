package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
)

type RoutineModule struct {
	Binder     binders.RoutineBinderInterface
	Controller controllers.RoutineControllerInterface
}

func NewRoutineModule() *RoutineModule {
	routineRepository := repositories.NewRoutineRepository(scopes.NewRoutineScope())

	routineService := services.NewRoutineService(
		models.NotezyDB,
		routineRepository,
	)

	routineBinder := binders.NewRoutineBinder()

	routineController := controllers.NewRoutineController(
		routineService,
	)

	return &RoutineModule{
		Binder:     routineBinder,
		Controller: routineController,
	}
}
