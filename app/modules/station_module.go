package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	scopes "notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
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
