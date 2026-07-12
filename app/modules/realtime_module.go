package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RealtimeModule struct {
	Binder     binders.RealtimeBinderInterface
	Controller controllers.RealtimeControllerInterface
}

func NewRealtimeModule() *RealtimeModule {
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())
	realtimeService := services.NewRealtimeService(models.NotezyDB, blockPackRepository)
	realtimeBinder := binders.NewRealtimeBinder()
	realtimeController := controllers.NewRealtimeController(realtimeService)

	return &RealtimeModule{
		Binder:     realtimeBinder,
		Controller: realtimeController,
	}
}
