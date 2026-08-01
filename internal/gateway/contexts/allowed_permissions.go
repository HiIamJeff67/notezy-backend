package contexts

import (
	"context"
	"net/http"
	"slices"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/internal/shared/contexts"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func WithAllowedPermissions(
	ctx context.Context,
	allowedPermissions []sharedtypes.AccessControlPermission,
) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		types.ContextFieldName_Allowed_Permissions,
		slices.Clone(allowedPermissions),
	)
}

func GetAllowedPermissions(
	ctx context.Context,
) ([]sharedtypes.AccessControlPermission, *exceptions.Exception) {
	allowedPermissions, err := sharedcontexts.GetValue[[]sharedtypes.AccessControlPermission](
		ctx,
		types.ContextFieldName_Allowed_Permissions,
	)
	if err != nil {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			"ReadAllowedPermissions",
			"The request context does not contain valid allowed permissions",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return slices.Clone(allowedPermissions), nil
}
