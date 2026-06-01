package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type StationModule struct {
	Binder     binders.StationBinderInterface
	Controller controllers.StationControllerInterface
}

func NewStationModule() *StationModule {
	stationRepository := repositories.NewStationRepository(scopes.NewStationScope())

	stationService := services.NewStationService(
		models.NotezyDB,
		stationRepository,
	)

	stationBinder := binders.NewStationBinder()

	stationController := controllers.NewStationController(
		stationService,
	)

	return &StationModule{
		Binder:     stationBinder,
		Controller: stationController,
	}
}
