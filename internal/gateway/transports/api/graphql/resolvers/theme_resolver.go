package resolvers

import (
	"context"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"
)

type ThemeResolverInterface interface{}

type ThemeResolver struct {
	*Resolver
}

func NewThemeResolver(resolver *Resolver) ThemeResolverInterface {
	return &ThemeResolver{
		Resolver: resolver,
	}
}

/* ============================== Resolver Methods ============================== */
// [MainSchema(as the filename) ---Indicator of MainSchema---> RelativeSchema(has the relationship between the MainSchema)]

// [PublicTheme ---PublicTheme.PublicId---> PublicUser]
func (r *ThemeResolver) Author(ctx context.Context, obj *gqlmodels.PublicTheme) (*gqlmodels.PublicUser, error) {
	return r.dataloader.UserDataLoader.LoadByThemePublicId(ctx, obj.PublicID)
}
