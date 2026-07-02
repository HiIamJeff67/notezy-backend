package resolvers

import (
	dataloaders "github.com/HiIamJeff67/notezy-backend/app/graphql/dataloaders"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataloader               dataloaders.Dataloaders
	userService              services.UserServiceInterface
	themeService             services.ThemeServiceInterface
	itemService              services.ItemServiceInterface
	blockService             services.BlockServiceInterface
	rootShelfService         services.RootShelfServiceInterface
	subShelfService          services.SubShelfServiceInterface
	stationService           services.StationServiceInterface
	routineService           services.RoutineServiceInterface
	routineTagService        services.RoutineTagServiceInterface
	routineTaskService       services.RoutineTaskServiceInterface
	routineTaskRecordService services.RoutineTaskRecordServiceInterface
}

func NewResolver(
	dataloader dataloaders.Dataloaders,
	userService services.UserServiceInterface,
	themeService services.ThemeServiceInterface,
	itemService services.ItemServiceInterface,
	blockService services.BlockServiceInterface,
	rootShelfService services.RootShelfServiceInterface,
	subShelfService services.SubShelfServiceInterface,
	stationService services.StationServiceInterface,
	routineService services.RoutineServiceInterface,
	routineTagService services.RoutineTagServiceInterface,
	routineTaskService services.RoutineTaskServiceInterface,
	routineTaskRecordService services.RoutineTaskRecordServiceInterface,
) *Resolver {
	return &Resolver{
		dataloader:               dataloader,
		userService:              userService,
		themeService:             themeService,
		itemService:              itemService,
		blockService:             blockService,
		rootShelfService:         rootShelfService,
		subShelfService:          subShelfService,
		stationService:           stationService,
		routineService:           routineService,
		routineTagService:        routineTagService,
		routineTaskService:       routineTaskService,
		routineTaskRecordService: routineTaskRecordService,
	}
}
