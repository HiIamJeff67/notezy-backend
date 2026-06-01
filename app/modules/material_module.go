package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	"github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	storages "github.com/HiIamJeff67/notezy-backend/app/storages"
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
