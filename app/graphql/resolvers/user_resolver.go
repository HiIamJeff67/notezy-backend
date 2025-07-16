package resolvers

import (
	"context"

	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserResolverInterface interface {
	UserInfo(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicUserInfo, error)
	Badge(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicBadge, error)
}

type UserResolver struct {
	*Resolver
	userService services.UserServiceInterface
}

func NewUserResolver(service services.UserServiceInterface) UserResolverInterface {
	return &UserResolver{
		userService: service,
	}
}

/* ============================== Resolver Methods ============================== */

// [PublicUser ---PublicUser.PublicId---> PublicUserInfo]
func (r *UserResolver) UserInfo(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicUserInfo, error) {
	future := r.dataloader.UserInfoLoader.Load(ctx, obj.PublicID)
	publicUserInfo, err := future()
	if err != nil {
		return nil, err
	}
	return publicUserInfo, nil
}

// [PublicUser ---PublicUser.PublicId---> PublicBadges]
func (r *UserResolver) Badge(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicBadge, error) {
	future := r.dataloader.BadgeLoader.Load(ctx, obj.PublicID)
	publicBadge, err := future()
	if err != nil {
		return nil, err
	}
	return publicBadge, nil
}
