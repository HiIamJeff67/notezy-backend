package resolvers

import "notezy-backend/app/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// add a field for dataloader here
	userService services.UserServiceInterface
}

func NewResolver(
	userService services.UserServiceInterface,
) *Resolver {
	return &Resolver{
		userService: userService,
	}
}
