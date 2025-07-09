package resolvers

type QueryResolverInterface interface{}

type QueryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolverInterface {
	return &QueryResolver{r}
}
