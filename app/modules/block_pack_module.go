package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	"github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type BlockPackModule struct {
	Binder     binders.BlockPackBinderInterface
	Controller controllers.BlockPackControllerInterface
}

func NewBlockPackModule() *BlockPackModule {
	subShelfRepository := repositories.NewSubShelfRepository(scopes.NewSubShelfScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())

	blockPackService := services.NewBlockPackService(
		models.NotezyDB,
		subShelfRepository,
		blockPackRepository,
	)

	blockPackBinder := binders.NewBlockPackBinder()

	blockPackController := controllers.NewBlockPackController(
		blockPackService,
	)

	return &BlockPackModule{
		Binder:     blockPackBinder,
		Controller: blockPackController,
	}
}
