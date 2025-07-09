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

// func (r *UserResolver) SearchUsers(ctx context.Context, input gqlmodels.SearchableUserInput) {
// 	// 調用 Service 層
// 	result, exception := r.userService.SearchUsers(ctx, input)
// 	if exception != nil {
// 		return nil
// 	}

// 	return result, nil
// }
