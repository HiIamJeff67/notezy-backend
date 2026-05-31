package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
)

type RoutineTagModule struct {
	Binder     binders.RoutineTagBinderInterface
	Controller controllers.RoutineTagControllerInterface
}

func NewRoutineTagModule() *RoutineTagModule {
	routineTagRepository := repositories.NewRoutineTagRepository(scopes.NewRoutineTagScope())

	routineTagService := services.NewRoutineTagService(
		models.NotezyDB,
		routineTagRepository,
	)

	routineTagBinder := binders.NewRoutineTagBinder()

	routineTagController := controllers.NewRoutineTagController(
		routineTagService,
	)

	return &RoutineTagModule{
		Binder:     routineTagBinder,
		Controller: routineTagController,
	}
}
