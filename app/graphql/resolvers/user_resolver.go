package resolvers

import (
	services "notezy-backend/app/services"
)

type UserResolverInterface interface{}

type UserResolver struct {
	*Resolver
	userService services.UserServiceInterface
}

func NewUserResolver(service services.UserServiceInterface) UserResolverInterface {
	return &UserResolver{
		userService: service,
	}
}

/* ============================== Resolvers ============================== */
