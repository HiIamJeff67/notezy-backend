package resolvers

type UserInfoResolverInterface interface{}

type UserInfoResolver struct {
	*Resolver
}

func NewUserInfoResolver() UserInfoResolverInterface {
	return &UserInfoResolver{}
}

/* ============================== Resolvers ============================== */
