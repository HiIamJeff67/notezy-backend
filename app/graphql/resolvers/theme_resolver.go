package resolvers

import (
	"context"

	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type ThemeResolverInterface interface{}

type ThemeResolver struct {
	*Resolver
	themeService services.ThemeServiceInterface
}

func NewThemeReolsver(service services.ThemeServiceInterface) ThemeResolverInterface {
	return &ThemeResolver{
		themeService: service,
	}
}

/* ============================== Resolver Methods ============================== */
// [MainSchema(as the filename) ---Indicator of MainSchema---> RelativeSchema(has the relationship between the MainSchema)]

// [PublicTheme ---PublicTheme.PublicId---> PublicUser]
func (r *ThemeResolver) Auther(ctx context.Context, obj *gqlmodels.PublicTheme) (*gqlmodels.PublicUser, error) {
	return r.dataloader.UserDataLoader.LoadByThemePublicId(ctx, obj.PublicID)
}
