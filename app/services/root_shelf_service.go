package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfServiceInterface interface {
	GetMyRootShelfById(ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception)
	SearchRecentRootShelves(ctx context.Context, reqDto *dtos.SearchRecentRootShelvesReqDto) (*dtos.SearchRecentRootShelvesResDto, *exceptions.Exception)
	CreateRootShelf(ctx context.Context, reqDto *dtos.CreateRootShelfReqDto) (*dtos.CreateRootShelfResDto, *exceptions.Exception)
	CreateRootShelves(ctx context.Context, reqDto *dtos.CreateRootShelvesReqDto) (*dtos.CreateRootShelvesResDto, *exceptions.Exception)
	UpdateMyRootShelfById(ctx context.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto) (*dtos.UpdateMyRootShelfByIdResDto, *exceptions.Exception)
	UpdateMyRootShelvesByIds(ctx context.Context, reqDto *dtos.UpdateMyRootShelvesByIdsReqDto) (*dtos.UpdateMyRootShelvesByIdsResDto, *exceptions.Exception)
	UpsertMyRootShelfPermission(ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto) (*dtos.UpsertMyRootShelfPermissionResDto, *exceptions.Exception)
	UpsertMyRootShelfPermissions(ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionsReqDto) (*dtos.UpsertMyRootShelfPermissionsResDto, *exceptions.Exception)
	RestoreMyRootShelfById(ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception)
	RestoreMyRootShelvesByIds(ctx context.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto) (*dtos.RestoreMyRootShelvesByIdsResDto, *exceptions.Exception)
	DeleteMyRootShelfById(ctx context.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto) (*dtos.DeleteMyRootShelfByIdResDto, *exceptions.Exception)
	DeleteMyRootShelvesByIds(ctx context.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto) (*dtos.DeleteMyRootShelvesByIdsResDto, *exceptions.Exception)
	DeleteMyRootShelfPermission(ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto) *exceptions.Exception
	DeleteMyRootShelfPermissions(ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionsReqDto) *exceptions.Exception

	SearchPrivateRootShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception)
}

type RootShelfService struct {
	db                  *gorm.DB
	rootShelfScope      scopes.RootShelfScopeInterface
	rootShelfRepository repositories.RootShelfRepositoryInterface
}

func NewRootShelfService(
	db *gorm.DB,
	rootShelfScope scopes.RootShelfScopeInterface,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
) RootShelfServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RootShelfService{
		db:                  db,
		rootShelfScope:      rootShelfScope,
		rootShelfRepository: rootShelfRepository,
	}
}

func (s *RootShelfService) GetMyRootShelfById(
	ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto,
) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	shelf, permission, exception := s.rootShelfRepository.GetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyRootShelfByIdResDto{
		Id:             shelf.Id,
		Name:           shelf.Name,
		Permission:     permission,
		SubShelfCount:  shelf.SubShelfCount,
		ItemCount:      shelf.ItemCount,
		LastAnalyzedAt: shelf.LastAnalyzedAt,
		DeletedAt:      shelf.DeletedAt,
		UpdatedAt:      shelf.UpdatedAt,
		CreatedAt:      shelf.CreatedAt,
	}, nil
}

func (s *RootShelfService) SearchRecentRootShelves(
	ctx context.Context, reqDto *dtos.SearchRecentRootShelvesReqDto,
) (*dtos.SearchRecentRootShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	resDto := dtos.SearchRecentRootShelvesResDto{}

	query := db.Model(&schemas.RootShelf{}).
		Where(`owner_id = ? AND "RootShelfTable".deleted_at IS NULL`,
			reqDto.ContextFields.UserId,
		)
	if len(strings.ReplaceAll(reqDto.Param.Query, " ", "")) > 0 {
		query = query.Where("name ILIKE ?", "%"+reqDto.Param.Query+"%")
	}

	result := query.Order("updated_at DESC").
		Limit(int(reqDto.Param.Limit)).
		Offset(int(reqDto.Param.Offset)).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(err)
	}

	for index := range resDto {
		resDto[index].Permission = enums.AccessControlPermission_Owner
	}

	return &resDto, nil
}

func (s *RootShelfService) CreateRootShelf(
	ctx context.Context, reqDto *dtos.CreateRootShelfReqDto,
) (*dtos.CreateRootShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	now := time.Now()
	newRootShelfId, exception := s.rootShelfRepository.CreateOne(
		reqDto.ContextFields.UserId,
		inputs.CreateRootShelfInput{
			Id:             reqDto.Body.Id,
			Name:           reqDto.Body.Name,
			LastAnalyzedAt: &now,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRootShelfResDto{
		Id:             *newRootShelfId,
		LastAnalyzedAt: now,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *RootShelfService) CreateRootShelves(
	ctx context.Context, reqDto *dtos.CreateRootShelvesReqDto,
) (*dtos.CreateRootShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	now := time.Now()
	input := make([]inputs.CreateRootShelfInput, len(reqDto.Body.CreatedRootShelves))
	for index, createdRootShelf := range reqDto.Body.CreatedRootShelves {
		input[index] = inputs.CreateRootShelfInput{
			Id:             createdRootShelf.Id,
			Name:           createdRootShelf.Name,
			LastAnalyzedAt: &now,
		}
	}
	newRootShelfIds, exception := s.rootShelfRepository.CreateMany(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRootShelvesResDto{
		Ids:            newRootShelfIds,
		LastAnalyzedAt: now,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelfById(
	ctx context.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto,
) (*dtos.UpdateMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	rootShelf, exception := s.rootShelfRepository.UpdateOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRootShelfInput{
			Values: inputs.UpdateRootShelfInput{
				Name: reqDto.Body.Values.Name,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRootShelfByIdResDto{
		UpdatedAt: rootShelf.UpdatedAt,
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelvesByIds(
	ctx context.Context, reqDto *dtos.UpdateMyRootShelvesByIdsReqDto,
) (*dtos.UpdateMyRootShelvesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.UpdateRootShelfByIdInput, len(reqDto.Body.UpdatedRootShelves))
	for index, updatedRootShelf := range reqDto.Body.UpdatedRootShelves {
		input[index] = inputs.UpdateRootShelfByIdInput{
			Id: updatedRootShelf.RootShelfId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateRootShelfInput]{
				Values: inputs.UpdateRootShelfInput{
					Name: updatedRootShelf.Values.Name,
				},
				SetNull: updatedRootShelf.SetNull,
			},
		}
	}
	exception := s.rootShelfRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRootShelvesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) UpsertMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto,
) (*dtos.UpsertMyRootShelfPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	if reqDto.Body.Permission == enums.AccessControlPermission_Owner {
		return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership through an access control")
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		},
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	var targetUser schemas.User
	result := tx.
		Model(&schemas.User{}).
		Where("public_id = ?", reqDto.Param.UserPublicId).
		First(&targetUser)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}

	targetPermission, targetException := s.rootShelfRepository.GetPermissionByRootShelfIdAndUserId(
		rootShelf.Id,
		targetUser.Id,
		options.WithDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if targetException != nil && !errors.Is(targetException.Origin, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, targetException
	}
	if targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("modify the RootShelf owner")
	}
	if actorPermission != enums.AccessControlPermission_Owner &&
		(reqDto.Body.Permission == enums.AccessControlPermission_Admin ||
			targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Admin) {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("grant Admin access")
	}

	permission, exception := s.rootShelfRepository.UpsertPermissionByUserId(
		rootShelf.Id,
		targetUser.Id,
		reqDto.Body.Permission,
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.UpsertMyRootShelfPermissionResDto{
		UserPublicId: targetUser.PublicId,
		Permission:   permission.Permission,
		UpdatedAt:    permission.UpdatedAt,
		CreatedAt:    permission.CreatedAt,
	}, nil
}

func (s *RootShelfService) UpsertMyRootShelfPermissions(
	ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionsReqDto,
) (*dtos.UpsertMyRootShelfPermissionsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	userPublicIds := make([]uuid.UUID, len(reqDto.Body.Permissions))
	permissionByPublicId := make(map[uuid.UUID]enums.AccessControlPermission, len(reqDto.Body.Permissions))
	for index, input := range reqDto.Body.Permissions {
		if input.Permission == enums.AccessControlPermission_Owner {
			return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership through permissions")
		}
		if _, exists := permissionByPublicId[input.UserPublicId]; exists {
			return nil, exceptions.Shelf.InvalidDto("permissions cannot contain duplicate userPublicIds")
		}

		userPublicIds[index] = input.UserPublicId
		permissionByPublicId[input.UserPublicId] = input.Permission
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		},
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", userPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if len(targetUsers) != len(userPublicIds) {
		tx.Rollback()
		return nil, exceptions.User.NotFound()
	}

	userByPublicId := make(map[uuid.UUID]schemas.User, len(targetUsers))
	userById := make(map[uuid.UUID]schemas.User, len(targetUsers))
	for _, user := range targetUsers {
		userByPublicId[user.PublicId] = user
		userById[user.Id] = user
	}

	userIds := make([]uuid.UUID, len(userPublicIds))
	for index, userPublicId := range userPublicIds {
		userIds[index] = userByPublicId[userPublicId].Id
	}

	existingPermissions, exception := s.rootShelfRepository.GetPermissionsByRootShelfIdAndUserIds(
		rootShelf.Id,
		userIds,
		options.WithDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	existingPermissionByUserId := make(map[uuid.UUID]enums.AccessControlPermission, len(existingPermissions))
	for _, existingPermission := range existingPermissions {
		existingPermissionByUserId[existingPermission.UserId] = existingPermission.Permission
	}

	permissions := make([]enums.AccessControlPermission, len(userIds))
	for index, userId := range userIds {
		user := userById[userId]
		permission := permissionByPublicId[user.PublicId]
		if existingPermissionByUserId[userId] == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return nil, exceptions.Shelf.NoPermission("modify the RootShelf owner")
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			(permission == enums.AccessControlPermission_Admin ||
				existingPermissionByUserId[userId] == enums.AccessControlPermission_Admin) {
			tx.Rollback()
			return nil, exceptions.Shelf.NoPermission("manage Admin permissions")
		}

		permissions[index] = permission
	}

	updatedPermissions, exception := s.rootShelfRepository.UpsertPermissionsByUserIds(
		rootShelf.Id,
		userIds,
		permissions,
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	updatedPermissionByUserId := make(map[uuid.UUID]schemas.UsersToShelves, len(updatedPermissions))
	for _, updatedPermission := range updatedPermissions {
		updatedPermissionByUserId[updatedPermission.UserId] = updatedPermission
	}

	resDto := make([]dtos.UpsertMyRootShelfPermissionResDto, len(userIds))
	for index, userId := range userIds {
		user := userById[userId]
		updatedPermission := updatedPermissionByUserId[userId]
		resDto[index] = dtos.UpsertMyRootShelfPermissionResDto{
			UserPublicId: user.PublicId,
			Permission:   updatedPermission.Permission,
			UpdatedAt:    updatedPermission.UpdatedAt,
			CreatedAt:    updatedPermission.CreatedAt,
		}
	}

	return &dtos.UpsertMyRootShelfPermissionsResDto{Permissions: resDto}, nil
}

func (s *RootShelfService) RestoreMyRootShelfById(
	ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto,
) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredRootShelf, exception := s.rootShelfRepository.RestoreSoftDeletedOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyRootShelfByIdResDto{
		Id:             restoredRootShelf.Id,
		Name:           restoredRootShelf.Name,
		SubShelfCount:  restoredRootShelf.SubShelfCount,
		ItemCount:      restoredRootShelf.ItemCount,
		LastAnalyzedAt: restoredRootShelf.LastAnalyzedAt,
		DeletedAt:      restoredRootShelf.DeletedAt,
		UpdatedAt:      restoredRootShelf.UpdatedAt,
		CreatedAt:      restoredRootShelf.CreatedAt,
	}, nil
}

func (s *RootShelfService) RestoreMyRootShelvesByIds(
	ctx context.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto,
) (*dtos.RestoreMyRootShelvesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredRootShelves, exception := s.rootShelfRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := dtos.RestoreMyRootShelvesByIdsResDto{}
	for _, restoredRootShelf := range restoredRootShelves {
		resDto = append(resDto, dtos.RestoreMyRootShelfByIdResDto{
			Id:             restoredRootShelf.Id,
			Name:           restoredRootShelf.Name,
			SubShelfCount:  restoredRootShelf.SubShelfCount,
			ItemCount:      restoredRootShelf.ItemCount,
			LastAnalyzedAt: restoredRootShelf.LastAnalyzedAt,
			DeletedAt:      restoredRootShelf.DeletedAt,
			UpdatedAt:      restoredRootShelf.UpdatedAt,
			CreatedAt:      restoredRootShelf.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *RootShelfService) DeleteMyRootShelfById(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto,
) (*dtos.DeleteMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.SoftDeleteOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRootShelfByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelvesByIds(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto,
) (*dtos.DeleteMyRootShelvesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.SoftDeleteManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRootShelvesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		},
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	var targetUser schemas.User
	result := tx.
		Model(&schemas.User{}).
		Where("public_id = ?", reqDto.Param.UserPublicId).
		First(&targetUser)
	if result.Error != nil {
		tx.Rollback()
		return exceptions.User.NotFound().WithOrigin(result.Error)
	}

	targetPermission, exception := s.rootShelfRepository.GetPermissionByRootShelfIdAndUserId(
		rootShelf.Id,
		targetUser.Id,
		options.WithDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return exceptions.Shelf.NoPermission("remove the RootShelf owner")
	}
	if actorPermission != enums.AccessControlPermission_Owner &&
		targetPermission.Permission == enums.AccessControlPermission_Admin {
		tx.Rollback()
		return exceptions.Shelf.NoPermission("revoke Admin access")
	}

	exception = s.rootShelfRepository.DeletePermissionByRootShelfIdAndUserId(
		rootShelf.Id,
		targetUser.Id,
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	if err := tx.Commit().Error; err != nil {
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	return nil
}

func (s *RootShelfService) DeleteMyRootShelfPermissions(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionsReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	userPublicIdSet := make(map[uuid.UUID]struct{}, len(reqDto.Body.UserPublicIds))
	for _, userPublicId := range reqDto.Body.UserPublicIds {
		if _, exists := userPublicIdSet[userPublicId]; exists {
			return exceptions.Shelf.InvalidDto("userPublicIds cannot contain duplicates")
		}

		userPublicIdSet[userPublicId] = struct{}{}
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		},
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", reqDto.Body.UserPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if len(targetUsers) != len(reqDto.Body.UserPublicIds) {
		tx.Rollback()
		return exceptions.User.NotFound()
	}

	userIdByPublicId := make(map[uuid.UUID]uuid.UUID, len(targetUsers))
	for _, targetUser := range targetUsers {
		userIdByPublicId[targetUser.PublicId] = targetUser.Id
	}

	userIds := make([]uuid.UUID, len(reqDto.Body.UserPublicIds))
	for index, userPublicId := range reqDto.Body.UserPublicIds {
		userIds[index] = userIdByPublicId[userPublicId]
	}

	targetPermissions, exception := s.rootShelfRepository.GetPermissionsByRootShelfIdAndUserIds(
		rootShelf.Id,
		userIds,
		options.WithDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if len(targetPermissions) != len(userIds) {
		tx.Rollback()
		return exceptions.Shelf.NotFound()
	}

	for _, targetPermission := range targetPermissions {
		if targetPermission.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return exceptions.Shelf.NoPermission("remove the RootShelf owner")
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			targetPermission.Permission == enums.AccessControlPermission_Admin {
			tx.Rollback()
			return exceptions.Shelf.NoPermission("revoke Admin access")
		}
	}

	exception = s.rootShelfRepository.DeletePermissionsByUserIds(
		rootShelf.Id,
		userIds,
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	if err := tx.Commit().Error; err != nil {
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	return nil
}

/* ============================== Service Methods for GraphQL RootShelf ============================== */

func (s *RootShelfService) SearchPrivateRootShelves(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput,
) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception) {
	type PrivateRootShelf struct {
		schemas.RootShelf
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	query := db.Model(&schemas.RootShelf{}).
		Select(`"RootShelfTable".*, uts.permission AS permission`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON "RootShelfTable".id = uts.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.rootShelfScope.FilterOnlyDeleted(types.Ternary_Negative))

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRootShelfCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRootShelfSortByName:
			query = query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRootShelfSortByLastUpdate:
			query = query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRootShelfSortByCreatedAt:
			query = query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query = query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var shelves []PrivateRootShelf
	if err := query.Scopes(s.rootShelfScope.IncludePreloads(
		[]schemas.RootShelfRelation{
			schemas.RootShelfRelation_UsersToShelves,
			schemas.RootShelfRelation_UsersToShelves_User,
			schemas.RootShelfRelation_Items,
		},
	)).Find(&shelves).Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(err)
	}

	hasNextPage := len(shelves) > limit
	searchEdges := make([]*gqlmodels.SearchRootShelfEdge, len(shelves))

	for index, shelf := range shelves {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRootShelfCursorFields]{
			Fields: gqlmodels.SearchRootShelfCursorFields{
				ID: shelf.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRootShelfEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                shelf.RootShelf.ToPrivateRootShelf(shelf.Permission),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchRootShelfConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
