package resolvers

import (
	"context"

	gqlmodels "notezy-backend/app/graphql/models"
)

/* ============================== Interface & Instance ============================== */

type UserResolverInterface interface {
	UserInfo(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicUserInfo, error)
	Badge(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicBadge, error)
}

type UserResolver struct {
	*Resolver
}

func NewUserResolver() UserResolverInterface {
	return &UserResolver{}
}

/* ============================== Resolver Methods ============================== */
// [MainSchema(as the filename) ---Indicator of MainSchema---> RelativeSchema(has the relationship between the MainSchema)]

// [PublicUser ---PublicUser.PublicId---> PublicUserInfo]
func (r *UserResolver) UserInfo(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicUserInfo, error) {
	return r.dataloader.UserInfoDataLoader.LoadByUserPublicId(ctx, obj.PublicID)
}

// [PublicUser ---PublicUser.PublicId---> PublicBadges]
func (r *UserResolver) Badge(ctx context.Context, obj *gqlmodels.PublicUser) (*gqlmodels.PublicBadge, error) {
	return r.dataloader.BadgeDataLoader.LoadByUserPublicId(ctx, obj.PublicID)
}
