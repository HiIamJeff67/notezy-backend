package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTaskRecordModule struct {
	Binder     binders.RoutineTaskRecordBinderInterface
	Controller controllers.RoutineTaskRecordControllerInterface
}

func NewRoutineTaskRecordModule() *RoutineTaskRecordModule {
	routineTaskRecordRepository := repositories.NewRoutineTaskRecordRepository(scopes.NewRoutineTaskRecordScope())
	routineTaskRecordService := services.NewRoutineTaskRecordService(
		models.NotezyDB,
		routineTaskRecordRepository,
	)

	return &RoutineTaskRecordModule{
		Binder:     binders.NewRoutineTaskRecordBinder(),
		Controller: controllers.NewRoutineTaskRecordController(routineTaskRecordService),
	}
}
