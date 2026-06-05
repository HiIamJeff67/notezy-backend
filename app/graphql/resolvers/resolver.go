package resolvers

import (
	dataloaders "github.com/HiIamJeff67/notezy-backend/app/graphql/dataloaders"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataloader         dataloaders.Dataloaders
	userService        services.UserServiceInterface
	themeService       services.ThemeServiceInterface
	rootShelfService   services.RootShelfServiceInterface
	stationService     services.StationServiceInterface
	routineService     services.RoutineServiceInterface
	routineTagService  services.RoutineTagServiceInterface
	routineTaskService services.RoutineTaskServiceInterface
}

func NewResolver(
	dataloader dataloaders.Dataloaders,
	userService services.UserServiceInterface,
	themeService services.ThemeServiceInterface,
	rootShelfService services.RootShelfServiceInterface,
	stationService services.StationServiceInterface,
	routineService services.RoutineServiceInterface,
	routineTagService services.RoutineTagServiceInterface,
	routineTaskService services.RoutineTaskServiceInterface,
) *Resolver {
	return &Resolver{
		dataloader:         dataloader,
		userService:        userService,
		themeService:       themeService,
		rootShelfService:   rootShelfService,
		stationService:     stationService,
		routineService:     routineService,
		routineTagService:  routineTagService,
		routineTaskService: routineTaskService,
	}
}
