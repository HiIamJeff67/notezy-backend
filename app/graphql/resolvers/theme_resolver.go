package resolvers

import (
	"context"
	gqlmodels "notezy-backend/app/graphql/models"
	"notezy-backend/app/services"
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

// [PublicTheme ---PublicTheme.PublicId---> PublicUser]
func (r *ThemeResolver) Auther(ctx context.Context, obj *gqlmodels.PublicTheme) {
	// future := r.dataloader.
}
