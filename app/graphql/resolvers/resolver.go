package resolvers

import (
	"notezy-backend/app/graphql/dataloaders"
	"notezy-backend/app/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataloader  dataloaders.Dataloaders
	userService services.UserServiceInterface
}

func NewResolver(
	dataloader dataloaders.Dataloaders,
	userService services.UserServiceInterface,
) *Resolver {
	return &Resolver{
		dataloader:  dataloader,
		userService: userService,
	}
}
