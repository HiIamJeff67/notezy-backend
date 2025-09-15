package resolvers

import (
	dataloaders "notezy-backend/app/graphql/dataloaders"
	services "notezy-backend/app/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataloader       dataloaders.Dataloaders
	userService      services.UserServiceInterface
	themeService     services.ThemeServiceInterface
	rootShelfService services.RootShelfServiceInterface
}

func NewResolver(
	dataloader dataloaders.Dataloaders,
	userService services.UserServiceInterface,
	themeService services.ThemeServiceInterface,
	rootShelfService services.RootShelfServiceInterface,
) *Resolver {
	return &Resolver{
		dataloader:       dataloader,
		userService:      userService,
		themeService:     themeService,
		rootShelfService: rootShelfService,
	}
}
