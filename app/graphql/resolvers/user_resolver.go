package resolvers

import (
	"context"
	gqlmodels "notezy-backend/app/graphql/models"
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

// resolver [PublicUser -> PublicUserInfo]
func (r *UserResolver) UserInfo(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicUserInfo, error) {
	future := r.dataloader.UserInfoDataloader.Load(ctx, obj.PublicID)
	userInfo, err := future()
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
