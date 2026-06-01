package modules

import (
	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type BlockGroupModule struct {
	Binder     binders.BlockGroupBinderInterface
	Controller controllers.BlockGroupControllerInterface
}

func NewBlockGroupModule() *BlockGroupModule {
	blockGroupScope := scopes.NewBlockGroupScope()
	blockGroupRepository := repositories.NewBlockGroupRepository(blockGroupScope)
	blockRepository := repositories.NewBlockRepository(scopes.NewBlockScope())
	editableBlockAdapter := adapters.NewEditableBlockAdapter()

	blockGroupService := services.NewBlockGroupService(
		models.NotezyDB,
		blockGroupRepository,
		blockRepository,
		editableBlockAdapter,
	)

	blockGroupBinder := binders.NewBlockGroupBinder()

	blockGroupController := controllers.NewBlockGroupController(
		blockGroupService,
	)

	return &BlockGroupModule{
		Binder:     blockGroupBinder,
		Controller: blockGroupController,
	}
}
