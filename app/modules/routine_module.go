package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineModule struct {
	Binder     binders.RoutineBinderInterface
	Controller controllers.RoutineControllerInterface
}

func NewRoutineModule() *RoutineModule {
	stationRepository := repositories.NewStationRepository(scopes.NewStationScope())
	routineRepository := repositories.NewRoutineRepository(scopes.NewRoutineScope())
	routineTagRepository := repositories.NewRoutineTagRepository(scopes.NewRoutineTagScope())
	routineTaskRepository := repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope())
	itemRepository := repositories.NewItemRepository(scopes.NewItemScope())

	routineService := services.NewRoutineService(
		models.NotezyDB,
		stationRepository,
		routineRepository,
		routineTagRepository,
		routineTaskRepository,
		itemRepository,
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
