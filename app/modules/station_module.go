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
	stationScope := scopes.NewStationScope()
	stationRepository := repositories.NewStationRepository(stationScope)
	usersToStationsRepository := repositories.NewUsersToStationsRepository()

	stationService := services.NewStationService(
		models.NotezyDB,
		stationScope,
		stationRepository,
		usersToStationsRepository,
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
