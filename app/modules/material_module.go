package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	"notezy-backend/app/models/scopes"
	services "notezy-backend/app/services"
	storages "notezy-backend/app/storages"
)

type MaterialModule struct {
	Binder     binders.MaterialBinderInterface
	Controller controllers.MaterialControllerInterface
}

func NewMaterialModule() *MaterialModule {
	subShelfRepository := repositories.NewSubShelfRepository(scopes.NewSubShelfScope())
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())

	materialService := services.NewMaterialService(
		models.NotezyDB,
		storages.InMemoryStorage,
		subShelfRepository,
		materialRepository,
	)

	materialBinder := binders.NewMaterialBinder()

	materialController := controllers.NewMaterialController(
		materialService,
	)

	return &MaterialModule{
		Binder:     materialBinder,
		Controller: materialController,
	}
}
