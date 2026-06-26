package modules

import (
	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTaskModule struct {
	Binder     binders.RoutineTaskBinderInterface
	Controller controllers.RoutineTaskControllerInterface
}

func NewRoutineTaskModule() *RoutineTaskModule {
	routineTaskRepository := repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope())
	editableBlockAdapter := adapters.NewEditableBlockAdapter()
	routineTaskPayloadAdapter := adapters.NewRoutineTaskPayloadAdapter(editableBlockAdapter)

	routineTaskService := services.NewRoutineTaskService(
		models.NotezyDB,
		routineTaskRepository,
		routineTaskPayloadAdapter,
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
