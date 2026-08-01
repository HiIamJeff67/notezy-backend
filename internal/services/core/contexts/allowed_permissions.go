package contexts

import (
	"context"
	"net/http"
	"slices"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/internal/shared/contexts"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func WithAllowedPermissions(
	ctx context.Context,
	allowedPermissions []enums.AccessControlPermission,
) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		types.ContextFieldName_Allowed_Permissions,
		slices.Clone(allowedPermissions),
	)
}

func WithActorUserId(ctx context.Context, actorUserId uuid.UUID) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		types.ContextFieldName_User_Id,
		actorUserId,
	)
}

func WithActorUserName(ctx context.Context, actorUserName string) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		types.ContextFieldName_User_Name,
		actorUserName,
	)
}

func WithActorUserPublicId(ctx context.Context, actorUserPublicId uuid.UUID) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		types.ContextFieldName_User_PublicId,
		actorUserPublicId,
	)
}

func GetAllowedPermissions(
	ctx context.Context,
) ([]enums.AccessControlPermission, *exceptions.Exception) {
	allowedPermissions, err := sharedcontexts.GetValue[[]enums.AccessControlPermission](
		ctx,
		types.ContextFieldName_Allowed_Permissions,
	)
	if err != nil {
		return nil, exceptions.New(
			"DelegationClaimsInvalid",
			"API",
			"ReadAllowedPermissions",
			"The verified delegation context does not contain valid allowed permissions",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return slices.Clone(allowedPermissions), nil
}

func GetActorUserId(ctx context.Context) (uuid.UUID, *exceptions.Exception) {
	actorUserId, err := sharedcontexts.GetValue[uuid.UUID](
		ctx,
		types.ContextFieldName_User_Id,
	)
	if err != nil || actorUserId == uuid.Nil {
		return uuid.Nil, exceptions.New(
			"DelegationClaimsInvalid",
			"API",
			"ReadActorUserId",
			"The verified delegation context does not contain a valid actor user ID",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return actorUserId, nil
}

func GetActorUserName(ctx context.Context) (string, *exceptions.Exception) {
	actorUserName, err := sharedcontexts.GetValue[string](
		ctx,
		types.ContextFieldName_User_Name,
	)
	if err != nil || actorUserName == "" {
		return "", exceptions.New(
			"DelegationClaimsInvalid",
			"API",
			"ReadActorUserName",
			"The verified delegation context does not contain a valid actor user name",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return actorUserName, nil
}

func GetActorUserPublicId(ctx context.Context) (uuid.UUID, *exceptions.Exception) {
	actorUserPublicId, err := sharedcontexts.GetValue[uuid.UUID](
		ctx,
		types.ContextFieldName_User_PublicId,
	)
	if err != nil || actorUserPublicId == uuid.Nil {
		return uuid.Nil, exceptions.New(
			"DelegationClaimsInvalid",
			"API",
			"ReadActorUserPublicId",
			"The verified delegation context does not contain a valid actor public ID",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return actorUserPublicId, nil
}
