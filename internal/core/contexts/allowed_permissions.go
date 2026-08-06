package contexts

import (
	"context"
	"net/http"
	"slices"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	enums "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas/enums"
)

func WithAllowedPermissions(
	ctx context.Context,
	allowedPermissions []enums.AccessControlPermission,
) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		sharedcontexts.ContextFieldName_Allowed_Permissions,
		slices.Clone(allowedPermissions),
	)
}

func WithActorUserId(ctx context.Context, actorUserId uuid.UUID) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		sharedcontexts.ContextFieldName_User_Id,
		actorUserId,
	)
}

func WithActorUserName(ctx context.Context, actorUserName string) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		sharedcontexts.ContextFieldName_User_Name,
		actorUserName,
	)
}

func WithActorUserPublicId(ctx context.Context, actorUserPublicId uuid.UUID) context.Context {
	return sharedcontexts.WithValue(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
		actorUserPublicId,
	)
}

func GetAllowedPermissions(
	ctx context.Context,
) ([]enums.AccessControlPermission, *exceptions.Exception) {
	allowedPermissions, err := sharedcontexts.GetValue[[]enums.AccessControlPermission](
		ctx,
		sharedcontexts.ContextFieldName_Allowed_Permissions,
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
		sharedcontexts.ContextFieldName_User_Id,
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
		sharedcontexts.ContextFieldName_User_Name,
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
		sharedcontexts.ContextFieldName_User_PublicId,
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
