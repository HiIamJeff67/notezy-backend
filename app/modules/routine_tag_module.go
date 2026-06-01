package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
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
