package contexts

import (
	"context"
	"slices"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func WithAllowedPermissions(
	ctx context.Context,
	allowedPermissions []enums.AccessControlPermission,
) context.Context {
	return context.WithValue(
		ctx,
		types.ContextFieldName_Allowed_Permissions,
		slices.Clone(allowedPermissions),
	)
}

func GetAllowedPermissions(
	ctx context.Context,
) ([]enums.AccessControlPermission, *exceptions.Exception) {
	value := ctx.Value(types.ContextFieldName_Allowed_Permissions)
	if value == nil {
		return nil, exceptions.Context.FailedToGetContextFieldOfSpecificName(
			types.ContextFieldName_Allowed_Permissions.String(),
		)
	}

	allowedPermissions, ok := value.([]enums.AccessControlPermission)
	if !ok {
		return nil, exceptions.Context.FailedToConvertContextFieldToSpecificType(
			"[]enums.AccessControlPermission",
		)
	}

	return slices.Clone(allowedPermissions), nil
}

func IntersectAllowedPermissions(
	ctx context.Context,
	permissions []enums.AccessControlPermission,
) []enums.AccessControlPermission {
	allowedPermissions, exception := GetAllowedPermissions(ctx)
	if exception != nil {
		return nil
	}

	intersection := make([]enums.AccessControlPermission, 0, len(permissions))
	for _, permission := range permissions {
		if slices.Contains(allowedPermissions, permission) {
			intersection = append(intersection, permission)
		}
	}

	return intersection
}
