package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	storages "github.com/HiIamJeff67/notezy-backend/app/storages"
)

type SubShelfModule struct {
	Binder     binders.SubShelfBinderInterface
	Controller controllers.SubShelfControllerInterface
}

func NewSubShelfModule() *SubShelfModule {
	subShelfScope := scopes.NewSubShelfScope()
	subShelfRepository := repositories.NewSubShelfRepository(subShelfScope)
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())

	subShelfService := services.NewSubShelfService(
		models.NotezyDB,
		storages.InMemoryStorage,
		subShelfScope,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
	)

	subShelfBinder := binders.NewSubShelfBinder()

	subShelfController := controllers.NewSubShelfController(
		subShelfService,
	)

	return &SubShelfModule{
		Binder:     subShelfBinder,
		Controller: subShelfController,
	}
}
