package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfServiceInterface interface {
	GetMyRootShelfById(ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception)
	CreateRootShelf(ctx context.Context, reqDto *dtos.CreateRootShelfReqDto) (*dtos.CreateRootShelfResDto, *exceptions.Exception)
	CreateRootShelves(ctx context.Context, reqDto *dtos.CreateRootShelvesReqDto) (*dtos.CreateRootShelvesResDto, *exceptions.Exception)
	UpdateMyRootShelfById(ctx context.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto) (*dtos.UpdateMyRootShelfByIdResDto, *exceptions.Exception)
	UpdateMyRootShelvesByIds(ctx context.Context, reqDto *dtos.UpdateMyRootShelvesByIdsReqDto) (*dtos.UpdateMyRootShelvesByIdsResDto, *exceptions.Exception)
	RestoreMyRootShelfById(ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception)
	RestoreMyRootShelvesByIds(ctx context.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto) (*dtos.RestoreMyRootShelvesByIdsResDto, *exceptions.Exception)
	DeleteMyRootShelfById(ctx context.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto) (*dtos.DeleteMyRootShelfByIdResDto, *exceptions.Exception)
	DeleteMyRootShelvesByIds(ctx context.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto) (*dtos.DeleteMyRootShelvesByIdsResDto, *exceptions.Exception)

	GetMyRootShelfPermission(ctx context.Context, reqDto *dtos.GetMyRootShelfPermissionReqDto) (*dtos.GetMyRootShelfPermissionResDto, *exceptions.Exception)
	CreateMyRootShelfPermission(ctx context.Context, reqDto *dtos.CreateMyRootShelfPermissionReqDto) (*dtos.CreateMyRootShelfPermissionResDto, *exceptions.Exception)
	UpsertMyRootShelfPermission(ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto) (*dtos.UpsertMyRootShelfPermissionResDto, *exceptions.Exception)
	UpsertMyRootShelfPermissions(ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionsReqDto) (*dtos.UpsertMyRootShelfPermissionsResDto, *exceptions.Exception)
	UpdateMyRootShelfPermission(ctx context.Context, reqDto *dtos.UpdateMyRootShelfPermissionReqDto) (*dtos.UpdateMyRootShelfPermissionResDto, *exceptions.Exception)
	TransferMyRootShelfOwnership(ctx context.Context, reqDto *dtos.TransferMyRootShelfOwnershipReqDto) (*dtos.TransferMyRootShelfOwnershipResDto, *exceptions.Exception)
	DeleteMyRootShelfPermission(ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto) *exceptions.Exception
	DeleteMyRootShelfPermissions(ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionsReqDto) *exceptions.Exception
	LeaveMyRootShelf(ctx context.Context, reqDto *dtos.LeaveMyRootShelfReqDto) *exceptions.Exception
	LeaveMyRootShelves(ctx context.Context, reqDto *dtos.LeaveMyRootShelvesReqDto) *exceptions.Exception

	SearchPrivateRootShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception)
}

type RootShelfService struct {
	db                       *gorm.DB
	rootShelfScope           scopes.RootShelfScopeInterface
	rootShelfRepository      repositories.RootShelfRepositoryInterface
	usersToShelvesRepository repositories.UsersToShelvesRepositoryInterface
	blockPackRepository      repositories.BlockPackRepositoryInterface
	realtimeLeaseStore       *caches.RealtimeLeaseStore
}

func NewRootShelfService(
	db *gorm.DB,
	rootShelfScope scopes.RootShelfScopeInterface,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
	usersToShelvesRepository repositories.UsersToShelvesRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	realtimeLeaseStore *caches.RealtimeLeaseStore,
) RootShelfServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	if realtimeLeaseStore == nil {
		realtimeLeaseStore = caches.NewRealtimeLeaseStore(caches.RedisClientMap)
	}
	return &RootShelfService{
		db:                       db,
		rootShelfScope:           rootShelfScope,
		rootShelfRepository:      rootShelfRepository,
		usersToShelvesRepository: usersToShelvesRepository,
		blockPackRepository:      blockPackRepository,
		realtimeLeaseStore:       realtimeLeaseStore,
	}
}

/* ============================== Auxiliary Functions ============================== */

func (s *RootShelfService) saveMyRootShelfPermission(
	ctx context.Context,
	actorUserId uuid.UUID,
	rootShelfId uuid.UUID,
	targetUserPublicId uuid.UUID,
	permission enums.AccessControlPermission,
	requireExisting *bool,
) (*dtos.UpsertMyRootShelfPermissionResDto, *exceptions.Exception) {
	if permission == enums.AccessControlPermission_Owner {
		return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership through an access control")
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Shelf.FailedToBeginTransaction(
			"Failed to begin RootShelf permission transaction",
		).WithOrigin(tx.Error)
	}

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		rootShelfId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	var targetUser schemas.User
	if result := tx.Where("public_id = ?", targetUserPublicId).First(&targetUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}

	targetPermission, targetException := s.usersToShelvesRepository.GetOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if targetException != nil && !errors.Is(targetException.Origin, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, targetException
	}
	if requireExisting != nil && *requireExisting != (targetPermission != nil) {
		tx.Rollback()
		if *requireExisting {
			return nil, targetException
		}
		return nil, exceptions.Shelf.NoChanges()
	}
	if targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("modify the RootShelf owner")
	}
	if actorPermission != enums.AccessControlPermission_Owner && (permission == enums.AccessControlPermission_Admin || targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Admin) {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("manage Admin permissions")
	}

	var relation *schemas.UsersToShelves
	if targetPermission == nil {
		relation, exception = s.usersToShelvesRepository.CreateOne(
			rootShelf.Id,
			targetUser.Id,
			permission,
			options.WithTransactionDB(tx),
		)
	} else {
		relation, exception = s.usersToShelvesRepository.UpdateOne(
			rootShelf.Id,
			targetUser.Id,
			permission,
			options.WithTransactionDB(tx),
		)
	}
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{rootShelf.Id},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &dtos.UpsertMyRootShelfPermissionResDto{
		UserPublicId: targetUser.PublicId,
		Permission:   relation.Permission,
		UpdatedAt:    relation.UpdatedAt,
		CreatedAt:    relation.CreatedAt,
	}, nil
}

/* ============================== Service Methods for RootShelf ============================== */

func (s *RootShelfService) GetMyRootShelfById(
	ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto,
) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
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
		options.WithAllowedPermissions(allowedPermissions),
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
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
		options.WithAllowedPermissions(allowedPermissions),
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
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
	exception = s.rootShelfRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRootShelvesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) RestoreMyRootShelfById(
	ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto,
) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredRootShelf, exception := s.rootShelfRepository.RestoreSoftDeletedOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredRootShelves, exception := s.rootShelfRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, permission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{rootShelf.Id},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}

	if permission == enums.AccessControlPermission_Owner {
		result := tx.
			Model(&schemas.RootShelf{}).
			Where("id = ?", rootShelf.Id).
			Update("deleted_at", time.Now())
		if result.Error != nil {
			tx.Rollback()
			return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, exceptions.Shelf.NoChanges()
		}
	} else {
		exception = s.usersToShelvesRepository.DeleteOne(
			rootShelf.Id,
			reqDto.ContextFields.UserId,
			options.WithTransactionDB(tx),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	if permission == enums.AccessControlPermission_Owner {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
	} else {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(reqDto.ContextFields.UserId, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		reqDto.Body.RootShelfIds,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}

	exception = s.rootShelfRepository.SoftDeleteManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &dtos.DeleteMyRootShelvesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) GetMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.GetMyRootShelfPermissionReqDto,
) (*dtos.GetMyRootShelfPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	if _, _, exception = s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	); exception != nil {
		return nil, exception
	}

	var targetUser schemas.User
	if result := db.Where("public_id = ?", reqDto.Param.UserPublicId).First(&targetUser); result.Error != nil {
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	relation, exception := s.usersToShelvesRepository.GetOne(
		reqDto.Param.RootShelfId,
		targetUser.Id,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyRootShelfPermissionResDto{UserPublicId: targetUser.PublicId, Permission: relation.Permission, UpdatedAt: relation.UpdatedAt, CreatedAt: relation.CreatedAt}, nil
}

func (s *RootShelfService) CreateMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.CreateMyRootShelfPermissionReqDto,
) (*dtos.CreateMyRootShelfPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	requireExisting := false
	return s.saveMyRootShelfPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.RootShelfId, reqDto.Param.UserPublicId, reqDto.Body.Permission, &requireExisting)
}

func (s *RootShelfService) UpsertMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto,
) (*dtos.UpsertMyRootShelfPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	return s.saveMyRootShelfPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.RootShelfId, reqDto.Param.UserPublicId, reqDto.Body.Permission, nil)
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
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

	existingPermissions, exception := s.usersToShelvesRepository.GetMany(
		rootShelf.Id,
		userIds,
		options.WithTransactionDB(tx),
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

	updatedPermissions, exception := s.usersToShelvesRepository.UpsertMany(
		rootShelf.Id,
		userIds,
		permissions,
		options.WithTransactionDB(tx),
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

func (s *RootShelfService) UpdateMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.UpdateMyRootShelfPermissionReqDto,
) (*dtos.UpdateMyRootShelfPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	requireExisting := true
	return s.saveMyRootShelfPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.RootShelfId, reqDto.Param.UserPublicId, reqDto.Body.Permission, &requireExisting)
}

func (s *RootShelfService) TransferMyRootShelfOwnership(
	ctx context.Context,
	reqDto *dtos.TransferMyRootShelfOwnershipReqDto,
) (*dtos.TransferMyRootShelfOwnershipResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Shelf.FailedToBeginTransaction(
			"Failed to begin RootShelf ownership transfer transaction",
		).WithOrigin(tx.Error)
	}
	rootShelf, permission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if permission != enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership")
	}

	var actorUser schemas.User
	if result := tx.Select("id, public_id").Where("id = ?", reqDto.ContextFields.UserId).First(&actorUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	var targetUser schemas.User
	if result := tx.Select("id, public_id").Where("public_id = ?", reqDto.Body.TargetUserPublicId).First(&targetUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if targetUser.Id == reqDto.ContextFields.UserId {
		tx.Rollback()
		return nil, exceptions.Shelf.NoChanges()
	}

	targetMembership, exception := s.usersToShelvesRepository.GetOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if targetMembership.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Shelf.NoChanges()
	}

	var accounts []schemas.UserAccount
	result := tx.
		Clauses(clause.Locking{Strength: options.LockingStrengthUpdate}).
		Where("user_id IN ?", []uuid.UUID{reqDto.ContextFields.UserId, targetUser.Id}).
		Order("user_id").
		Find(&accounts)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
	}
	if len(accounts) != 2 {
		tx.Rollback()
		return nil, exceptions.User.NotFound()
	}

	var maximumSubscribers int32
	result = tx.
		Model(&schemas.User{}).
		Select(`"PlanLimitationTable".max_realtime_room_subscriber_count`).
		Joins(`INNER JOIN "PlanLimitationTable" ON "PlanLimitationTable".key = "UserTable".plan`).
		Where(`"UserTable".id = ?`, targetUser.Id).
		Scan(&maximumSubscribers)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		tx.Rollback()
		return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership to a plan without realtime room capacity")
	}

	var blockPackIds []uuid.UUID
	result = tx.
		Model(&schemas.BlockPack{}).
		Select(`"BlockPackTable".id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".id = "BlockPackTable".parent_sub_shelf_id`).
		Where(`"SubShelfTable".root_shelf_id = ?`, rootShelf.Id).
		Where(`"BlockPackTable".deleted_at IS NULL`).
		Where(`"SubShelfTable".deleted_at IS NULL`).
		Find(&blockPackIds)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(result.Error)
	}

	subscriberCounts, err := s.realtimeLeaseStore.GetBlockPackSubscriberCounts(blockPackIds)
	if err != nil {
		tx.Rollback()
		return nil, exceptions.Cache.FailedToUpdate("realtime block pack subscriber counts").WithOrigin(err)
	}
	for _, subscriberCount := range subscriberCounts {
		if subscriberCount > int64(maximumSubscribers) {
			tx.Rollback()
			return nil, exceptions.Shelf.NoPermission("transfer RootShelf ownership to a plan with insufficient realtime room capacity")
		}
	}

	if _, exception = s.usersToShelvesRepository.UpdateOne(
		rootShelf.Id,
		reqDto.ContextFields.UserId,
		enums.AccessControlPermission_Admin,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newOwnerMembership, exception := s.usersToShelvesRepository.UpdateOne(
		rootShelf.Id,
		targetUser.Id,
		enums.AccessControlPermission_Owner,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	result = tx.Model(&schemas.RootShelf{}).
		Where("id = ?", rootShelf.Id).
		Update("owner_id", targetUser.Id)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, exceptions.Shelf.NotFound()
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.TransferMyRootShelfOwnershipResDto{
		RootShelfId:               rootShelf.Id,
		PreviousOwnerUserPublicId: actorUser.PublicId,
		NewOwnerUserPublicId:      targetUser.PublicId,
		UpdatedAt:                 newOwnerMembership.UpdatedAt,
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelfPermission(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelfPermissionReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
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

	targetPermission, exception := s.usersToShelvesRepository.GetOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
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

	exception = s.usersToShelvesRepository.DeleteOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{rootShelf.Id},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
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

	targetPermissions, exception := s.usersToShelvesRepository.GetMany(
		rootShelf.Id,
		userIds,
		options.WithTransactionDB(tx),
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

	exception = s.usersToShelvesRepository.DeleteMany(
		rootShelf.Id,
		userIds,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{rootShelf.Id},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	for _, targetUser := range targetUsers {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
	}

	return nil
}

func (s *RootShelfService) LeaveMyRootShelf(
	ctx context.Context, reqDto *dtos.LeaveMyRootShelfReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Shelf.FailedToBeginTransaction("Failed to begin RootShelf leave transaction").WithOrigin(tx.Error)
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{reqDto.Param.RootShelfId},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}
	rootShelf, permission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		nil,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return exceptions.Shelf.NoPermission("transfer RootShelf ownership before leaving")
	}
	if exception = s.usersToShelvesRepository.DeleteOne(
		rootShelf.Id,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(reqDto.ContextFields.UserId, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return nil
}

func (s *RootShelfService) LeaveMyRootShelves(
	ctx context.Context, reqDto *dtos.LeaveMyRootShelvesReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	rootShelfIdSet := make(map[uuid.UUID]struct{}, len(reqDto.Body.RootShelves))
	rootShelfIds := make([]uuid.UUID, len(reqDto.Body.RootShelves))
	for index, rootShelfReqDto := range reqDto.Body.RootShelves {
		if _, exists := rootShelfIdSet[rootShelfReqDto.RootShelfId]; exists {
			return exceptions.Shelf.InvalidDto("rootShelves cannot contain duplicate rootShelfIds")
		}
		rootShelfIdSet[rootShelfReqDto.RootShelfId] = struct{}{}
		rootShelfIds[index] = rootShelfReqDto.RootShelfId
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Shelf.FailedToBeginTransaction("Failed to begin RootShelf leave transaction").WithOrigin(tx.Error)
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		rootShelfIds,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	blockPackIds := make([]uuid.UUID, len(blockPacks))
	for index, blockPack := range blockPacks {
		blockPackIds[index] = blockPack.Id
	}
	relations, exception := s.usersToShelvesRepository.GetManyByRootShelfIdsAndUserId(
		rootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if len(relations) != len(rootShelfIds) {
		tx.Rollback()
		return exceptions.Shelf.NotFound()
	}
	for _, relation := range relations {
		if relation.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return exceptions.Shelf.NoPermission("transfer RootShelf ownership before leaving")
		}
	}

	if exception = s.usersToShelvesRepository.DeleteManyByRootShelfIdsAndUserId(
		rootShelfIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(reqDto.ContextFields.UserId, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Negative
	if gqlInput.IsDeletedAt != nil && *gqlInput.IsDeletedAt {
		onlyDeleted = types.Ternary_Positive
	}

	query := db.Model(&schemas.RootShelf{}).
		Select(`"RootShelfTable".*, uts.permission AS permission`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON "RootShelfTable".id = uts.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.rootShelfScope.FilterOnlyDeleted(onlyDeleted))

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
			schemas.RootShelfRelation_Items,
		},
	)).Find(&shelves).Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(err)
	}

	userIds := make([]uuid.UUID, 0)
	userIdsSeen := make(map[uuid.UUID]struct{})
	for _, shelf := range shelves {
		if _, exists := userIdsSeen[shelf.OwnerId]; !exists {
			userIds = append(userIds, shelf.OwnerId)
			userIdsSeen[shelf.OwnerId] = struct{}{}
		}

		for _, usersToShelf := range shelf.UsersToShelves {
			if _, exists := userIdsSeen[usersToShelf.UserId]; exists {
				continue
			}

			userIds = append(userIds, usersToShelf.UserId)
			userIdsSeen[usersToShelf.UserId] = struct{}{}
		}
	}

	users := make([]schemas.User, 0, len(userIds))
	if len(userIds) > 0 {
		if err := db.
			Where("id IN ?", userIds).
			Find(&users).Error; err != nil {
			return nil, exceptions.User.NotFound().WithOrigin(err)
		}
	}

	publicUsersById := make(map[uuid.UUID]*gqlmodels.PublicUser, len(users))
	for _, user := range users {
		publicUsersById[user.Id] = user.ToPublicUser()
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

		privateRootShelf := shelf.RootShelf.ToPrivateRootShelf(shelf.Permission)
		owner, exists := publicUsersById[shelf.OwnerId]
		if !exists {
			return nil, exceptions.User.NotFound()
		}

		privateRootShelf.Owner = owner
		for _, usersToShelf := range shelf.UsersToShelves {
			if usersToShelf.UserId == shelf.OwnerId {
				continue
			}

			sharer, exists := publicUsersById[usersToShelf.UserId]
			if !exists {
				return nil, exceptions.User.NotFound()
			}

			privateRootShelf.Sharers = append(privateRootShelf.Sharers, sharer)
		}

		searchEdges[index] = &gqlmodels.SearchRootShelfEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                privateRootShelf,
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
