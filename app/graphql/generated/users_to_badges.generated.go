// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	gqlmodels "notezy-backend/app/graphql/models"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _PublicUsersToBadges_userId(ctx context.Context, field graphql.CollectedField, obj *gqlmodels.PublicUsersToBadges) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_PublicUsersToBadges_userId(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (any, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.UserID, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(uuid.UUID)
	fc.Result = res
	return ec.marshalNUUID2githubᚗcomᚋgoogleᚋuuidᚐUUID(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_PublicUsersToBadges_userId(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "PublicUsersToBadges",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type UUID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _PublicUsersToBadges_badgeId(ctx context.Context, field graphql.CollectedField, obj *gqlmodels.PublicUsersToBadges) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_PublicUsersToBadges_badgeId(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (any, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.BadgeID, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(uuid.UUID)
	fc.Result = res
	return ec.marshalNUUID2githubᚗcomᚋgoogleᚋuuidᚐUUID(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_PublicUsersToBadges_badgeId(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "PublicUsersToBadges",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type UUID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _PublicUsersToBadges_createdAt(ctx context.Context, field graphql.CollectedField, obj *gqlmodels.PublicUsersToBadges) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_PublicUsersToBadges_createdAt(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (any, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.CreatedAt, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(time.Time)
	fc.Result = res
	return ec.marshalNTime2timeᚐTime(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_PublicUsersToBadges_createdAt(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "PublicUsersToBadges",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Time does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _PublicUsersToBadges_user(ctx context.Context, field graphql.CollectedField, obj *gqlmodels.PublicUsersToBadges) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_PublicUsersToBadges_user(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (any, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.User, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*gqlmodels.PublicUser)
	fc.Result = res
	return ec.marshalNPublicUser2ᚖnotezyᚑbackendᚋappᚋgraphqlᚋmodelsᚐPublicUser(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_PublicUsersToBadges_user(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "PublicUsersToBadges",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext_PublicUser_name(ctx, field)
			case "displayName":
				return ec.fieldContext_PublicUser_displayName(ctx, field)
			case "email":
				return ec.fieldContext_PublicUser_email(ctx, field)
			case "role":
				return ec.fieldContext_PublicUser_role(ctx, field)
			case "plan":
				return ec.fieldContext_PublicUser_plan(ctx, field)
			case "status":
				return ec.fieldContext_PublicUser_status(ctx, field)
			case "createdAt":
				return ec.fieldContext_PublicUser_createdAt(ctx, field)
			case "updatedAt":
				return ec.fieldContext_PublicUser_updatedAt(ctx, field)
			case "userInfo":
				return ec.fieldContext_PublicUser_userInfo(ctx, field)
			case "badges":
				return ec.fieldContext_PublicUser_badges(ctx, field)
			case "themes":
				return ec.fieldContext_PublicUser_themes(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type PublicUser", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _PublicUsersToBadges_badge(ctx context.Context, field graphql.CollectedField, obj *gqlmodels.PublicUsersToBadges) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_PublicUsersToBadges_badge(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (any, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Badge, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*gqlmodels.PublicBadge)
	fc.Result = res
	return ec.marshalNPublicBadge2ᚖnotezyᚑbackendᚋappᚋgraphqlᚋmodelsᚐPublicBadge(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_PublicUsersToBadges_badge(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "PublicUsersToBadges",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_PublicBadge_id(ctx, field)
			case "title":
				return ec.fieldContext_PublicBadge_title(ctx, field)
			case "description":
				return ec.fieldContext_PublicBadge_description(ctx, field)
			case "type":
				return ec.fieldContext_PublicBadge_type(ctx, field)
			case "imageURL":
				return ec.fieldContext_PublicBadge_imageURL(ctx, field)
			case "createdAt":
				return ec.fieldContext_PublicBadge_createdAt(ctx, field)
			case "users":
				return ec.fieldContext_PublicBadge_users(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type PublicBadge", field.Name)
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var publicUsersToBadgesImplementors = []string{"PublicUsersToBadges"}

func (ec *executionContext) _PublicUsersToBadges(ctx context.Context, sel ast.SelectionSet, obj *gqlmodels.PublicUsersToBadges) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, publicUsersToBadgesImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("PublicUsersToBadges")
		case "userId":
			out.Values[i] = ec._PublicUsersToBadges_userId(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "badgeId":
			out.Values[i] = ec._PublicUsersToBadges_badgeId(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "createdAt":
			out.Values[i] = ec._PublicUsersToBadges_createdAt(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "user":
			out.Values[i] = ec._PublicUsersToBadges_user(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "badge":
			out.Values[i] = ec._PublicUsersToBadges_badge(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

// endregion ***************************** type.gotpl *****************************
