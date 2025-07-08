package resolvers

import (
	"context"
	gqlmodels "notezy-backend/app/graphql/models"
)

type QueryResolverInterface interface{}

type QueryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolverInterface {
	return &QueryResolver{r}
}

func (r *QueryResolver) SearchUsers(ctx context.Context, input gqlmodels.SearchableUserInput) {

}
