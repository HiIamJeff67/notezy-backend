package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RootShelfModule struct {
	Binder     binders.RootShelfBinderInterface
	Controller controllers.RootShelfControllerInterface
}

func NewRootShelfModule() *RootShelfModule {
	rootShelfScope := scopes.NewRootShelfScope()
	rootShelfRepository := repositories.NewRootShelfRepository(rootShelfScope)

	rootShelfService := services.NewRootShelfService(
		models.NotezyDB,
		rootShelfScope,
		rootShelfRepository,
	)

	rootShelfBinder := binders.NewRootShelfBinder()

	rootShelfController := controllers.NewRootShelfController(
		rootShelfService,
	)

	return &RootShelfModule{
		Binder:     rootShelfBinder,
		Controller: rootShelfController,
	}
}
