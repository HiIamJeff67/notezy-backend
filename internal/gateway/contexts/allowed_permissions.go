package contexts

import (
	"context"
	"net/http"
	"slices"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
)

func WithAllowedPermissions(
	ctx context.Context,
	allowedPermissions []enumcontract.AccessControlPermission,
) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		sharedcontexts.ContextFieldName_Allowed_Permissions,
		slices.Clone(allowedPermissions),
	)
}

func GetAllowedPermissions(
	ctx context.Context,
) ([]enumcontract.AccessControlPermission, *exceptions.Exception) {
	allowedPermissions, err := sharedcontexts.GetValue[[]enumcontract.AccessControlPermission](
		ctx,
		sharedcontexts.ContextFieldName_Allowed_Permissions,
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
