package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/root-shelves"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/searchcursor"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/shared/realtimelease"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RootShelfServiceInterface interface {
	GetMyRootShelfById(ctx context.Context, requestDto *rootshelvesdto.GetMyRootShelfByIdRequestDto) (*rootshelvesdto.GetMyRootShelfByIdResponseDto, *exceptions.Exception)
	CreateRootShelf(ctx context.Context, requestDto *rootshelvesdto.CreateRootShelfRequestDto) (*rootshelvesdto.CreateRootShelfResponseDto, *exceptions.Exception)
	CreateRootShelves(ctx context.Context, requestDto *rootshelvesdto.CreateRootShelvesRequestDto) (*rootshelvesdto.CreateRootShelvesResponseDto, *exceptions.Exception)
	UpdateMyRootShelfById(ctx context.Context, requestDto *rootshelvesdto.UpdateMyRootShelfByIdRequestDto) (*rootshelvesdto.UpdateMyRootShelfByIdResponseDto, *exceptions.Exception)
	UpdateMyRootShelvesByIds(ctx context.Context, requestDto *rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto) (*rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto, *exceptions.Exception)
	RestoreMyRootShelfById(ctx context.Context, requestDto *rootshelvesdto.RestoreMyRootShelfByIdRequestDto) (*rootshelvesdto.RestoreMyRootShelfByIdResponseDto, *exceptions.Exception)
	RestoreMyRootShelvesByIds(ctx context.Context, requestDto *rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto) (*rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto, *exceptions.Exception)
	DeleteMyRootShelfById(ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelfByIdRequestDto) (*rootshelvesdto.DeleteMyRootShelfByIdResponseDto, *exceptions.Exception)
	DeleteMyRootShelvesByIds(ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto) (*rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto, *exceptions.Exception)

	GetMyRootShelfPermission(ctx context.Context, requestDto *rootshelvesdto.GetMyRootShelfPermissionRequestDto) (*rootshelvesdto.GetMyRootShelfPermissionResponseDto, *exceptions.Exception)
	CreateMyRootShelfPermission(ctx context.Context, requestDto *rootshelvesdto.CreateMyRootShelfPermissionRequestDto) (*rootshelvesdto.CreateMyRootShelfPermissionResponseDto, *exceptions.Exception)
	UpsertMyRootShelfPermission(ctx context.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionRequestDto) (*rootshelvesdto.UpsertMyRootShelfPermissionResponseDto, *exceptions.Exception)
	UpsertMyRootShelfPermissions(ctx context.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto) (*rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto, *exceptions.Exception)
	UpdateMyRootShelfPermission(ctx context.Context, requestDto *rootshelvesdto.UpdateMyRootShelfPermissionRequestDto) (*rootshelvesdto.UpdateMyRootShelfPermissionResponseDto, *exceptions.Exception)
	TransferMyRootShelfOwnership(ctx context.Context, requestDto *rootshelvesdto.TransferMyRootShelfOwnershipRequestDto) (*rootshelvesdto.TransferMyRootShelfOwnershipResponseDto, *exceptions.Exception)
	DeleteMyRootShelfPermission(ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionRequestDto) (*rootshelvesdto.DeleteMyRootShelfPermissionResponseDto, *exceptions.Exception)
	DeleteMyRootShelfPermissions(ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto) (*rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto, *exceptions.Exception)
	LeaveMyRootShelf(ctx context.Context, requestDto *rootshelvesdto.LeaveMyRootShelfRequestDto) *exceptions.Exception
	LeaveMyRootShelves(ctx context.Context, requestDto *rootshelvesdto.LeaveMyRootShelvesRequestDto) *exceptions.Exception

	SearchPrivateRootShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception)
}

type RootShelfService struct {
	db                       *gorm.DB
	rootShelfScope           scopes.RootShelfScopeInterface
	rootShelfRepository      repositories.RootShelfRepositoryInterface
	usersToShelvesRepository repositories.UsersToShelvesRepositoryInterface
	blockPackRepository      repositories.BlockPackRepositoryInterface
	realtimeLeaseStore       *realtimelease.RealtimeLeaseStore
}

func NewRootShelfService(
	db *gorm.DB,
	rootShelfScope scopes.RootShelfScopeInterface,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
	usersToShelvesRepository repositories.UsersToShelvesRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	realtimeLeaseStore *realtimelease.RealtimeLeaseStore,
) RootShelfServiceInterface {
	if db == nil {
		db = data.NotezyDB
	}
	if realtimeLeaseStore == nil {
		realtimeLeaseStore = realtimelease.NewRealtimeLeaseStore()
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
) (*rootshelvesdto.RootShelfPermissionResponseDto, *exceptions.Exception) {
	if permission == enums.AccessControlPermission_Owner {
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"SaveMyRootShelfPermission",
			"Failed to begin the root shelf permission transaction",
			http.StatusInternalServerError,
			true,
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
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}

	targetPermission, targetException := s.usersToShelvesRepository.GetOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if targetException != nil && !errors.Is(targetException.Origin(), gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, targetException
	}
	if requireExisting != nil && *requireExisting != (targetPermission != nil) {
		tx.Rollback()
		if *requireExisting {
			return nil, targetException
		}
		return nil, exceptions.New(
			"NoChanges",
			"RootShelf",
			"Manage",
			"No root shelf changes were applied",
			http.StatusNotModified,
		)
	}
	if targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}
	if actorPermission != enums.AccessControlPermission_Owner && (permission == enums.AccessControlPermission_Admin || targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Admin) {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
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
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &rootshelvesdto.RootShelfPermissionResponseDto{
		UserPublicId: targetUser.PublicId,
		Permission:   relation.Permission.String(),
		UpdatedAt:    relation.UpdatedAt,
		CreatedAt:    relation.CreatedAt,
	}, nil
}

/* ============================== Service Methods for RootShelf ============================== */

func (s *RootShelfService) GetMyRootShelfById(
	ctx context.Context, requestDto *rootshelvesdto.GetMyRootShelfByIdRequestDto,
) (*rootshelvesdto.GetMyRootShelfByIdResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.IsDeleted != nil {
		if *requestDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	shelf, permission, exception := s.rootShelfRepository.GetOneById(
		requestDto.Param.RootShelfId,
		actorUserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.GetMyRootShelfByIdResponseDto{
		Id:             shelf.Id,
		Name:           shelf.Name,
		Permission:     permission.String(),
		SubShelfCount:  shelf.SubShelfCount,
		ItemCount:      shelf.ItemCount,
		LastAnalyzedAt: shelf.LastAnalyzedAt,
		DeletedAt:      shelf.DeletedAt,
		UpdatedAt:      shelf.UpdatedAt,
		CreatedAt:      shelf.CreatedAt,
	}, nil
}

func (s *RootShelfService) CreateRootShelf(
	ctx context.Context, requestDto *rootshelvesdto.CreateRootShelfRequestDto,
) (*rootshelvesdto.CreateRootShelfResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	db := s.db.WithContext(ctx)

	now := time.Now()
	newRootShelfId, exception := s.rootShelfRepository.CreateOne(
		actorUserId,
		inputs.CreateRootShelfInput{
			Id:             requestDto.Body.Id,
			Name:           requestDto.Body.Name,
			LastAnalyzedAt: &now,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.CreateRootShelfResponseDto{
		Id:             *newRootShelfId,
		LastAnalyzedAt: now,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *RootShelfService) CreateRootShelves(
	ctx context.Context, requestDto *rootshelvesdto.CreateRootShelvesRequestDto,
) (*rootshelvesdto.CreateRootShelvesResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	db := s.db.WithContext(ctx)

	now := time.Now()
	input := make([]inputs.CreateRootShelfInput, len(requestDto.Body.RootShelves))
	for index, createdRootShelf := range requestDto.Body.RootShelves {
		input[index] = inputs.CreateRootShelfInput{
			Id:             createdRootShelf.Id,
			Name:           createdRootShelf.Name,
			LastAnalyzedAt: &now,
		}
	}
	newRootShelfIds, exception := s.rootShelfRepository.CreateMany(
		actorUserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.CreateRootShelvesResponseDto{
		Ids:            newRootShelfIds,
		LastAnalyzedAt: now,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelfById(
	ctx context.Context, requestDto *rootshelvesdto.UpdateMyRootShelfByIdRequestDto,
) (*rootshelvesdto.UpdateMyRootShelfByIdResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	rootShelf, exception := s.rootShelfRepository.UpdateOneById(
		requestDto.Param.RootShelfId,
		actorUserId,
		inputs.PartialUpdateRootShelfInput{
			Values: inputs.UpdateRootShelfInput{
				Name: requestDto.Body.Values.Name,
			},
			SetNull: requestDto.Body.SetNull,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.UpdateMyRootShelfByIdResponseDto{
		UpdatedAt: rootShelf.UpdatedAt,
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelvesByIds(
	ctx context.Context,
	requestDto *rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto,
) (*rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	db := s.db.WithContext(ctx)
	input := make([]inputs.UpdateRootShelfByIdInput, len(requestDto.Body.UpdatedRootShelves))
	for index, updatedRootShelf := range requestDto.Body.UpdatedRootShelves {
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
		actorUserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) RestoreMyRootShelfById(
	ctx context.Context,
	requestDto *rootshelvesdto.RestoreMyRootShelfByIdRequestDto,
) (*rootshelvesdto.RestoreMyRootShelfByIdResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredRootShelf, exception := s.rootShelfRepository.RestoreSoftDeletedOneById(
		requestDto.Body.RootShelfId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.RestoreMyRootShelfByIdResponseDto{
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
	ctx context.Context,
	requestDto *rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto,
) (*rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredRootShelves, exception := s.rootShelfRepository.RestoreSoftDeletedManyByIds(
		requestDto.Body.RootShelfIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	responseDto := rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto{}
	for _, restoredRootShelf := range restoredRootShelves {
		responseDto = append(responseDto, rootshelvesdto.RestoreMyRootShelfByIdResponseDto{
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

	return &responseDto, nil
}

func (s *RootShelfService) DeleteMyRootShelfById(
	ctx context.Context,
	requestDto *rootshelvesdto.DeleteMyRootShelfByIdRequestDto,
) (*rootshelvesdto.DeleteMyRootShelfByIdResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, permission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Body.RootShelfId,
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
			return nil, exceptions.New(
				"FailedToUpdate",
				"RootShelf",
				"Manage",
				"Failed to update the root shelf",
				http.StatusInternalServerError,
				true,
			).WithOrigin(result.Error)
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, exceptions.New(
				"NoChanges",
				"RootShelf",
				"Manage",
				"No root shelf changes were applied",
				http.StatusNotModified,
			)
		}
	} else {
		exception = s.usersToShelvesRepository.DeleteOne(
			rootShelf.Id,
			actorUserId,
			options.WithTransactionDB(tx),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if permission == enums.AccessControlPermission_Owner {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
	} else {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(actorUserId, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
	}

	return &rootshelvesdto.DeleteMyRootShelfByIdResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelvesByIds(
	ctx context.Context,
	requestDto *rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto,
) (*rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		requestDto.Body.RootShelfIds,
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
		requestDto.Body.RootShelfIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) GetMyRootShelfPermission(
	ctx context.Context, requestDto *rootshelvesdto.GetMyRootShelfPermissionRequestDto,
) (*rootshelvesdto.GetMyRootShelfPermissionResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	if _, _, exception = s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Param.RootShelfId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	); exception != nil {
		return nil, exception
	}

	var targetUser schemas.User
	if result := db.Where("public_id = ?", requestDto.Param.UserPublicId).First(&targetUser); result.Error != nil {
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}
	relation, exception := s.usersToShelvesRepository.GetOne(
		requestDto.Param.RootShelfId,
		targetUser.Id,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &rootshelvesdto.GetMyRootShelfPermissionResponseDto{
		UserPublicId: targetUser.PublicId,
		Permission:   relation.Permission.String(),
		UpdatedAt:    relation.UpdatedAt,
		CreatedAt:    relation.CreatedAt,
	}, nil
}

func (s *RootShelfService) CreateMyRootShelfPermission(
	ctx context.Context, requestDto *rootshelvesdto.CreateMyRootShelfPermissionRequestDto,
) (*rootshelvesdto.CreateMyRootShelfPermissionResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	permission, err := enums.ConvertStringToAccessControlPermission(requestDto.Body.Permission)
	if err != nil {
		return nil, exceptions.InvalidInput("RootShelf").WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	requireExisting := false
	return s.saveMyRootShelfPermission(ctx, actorUserId, requestDto.Param.RootShelfId, requestDto.Param.UserPublicId, *permission, &requireExisting)
}

func (s *RootShelfService) UpsertMyRootShelfPermission(
	ctx context.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionRequestDto,
) (*rootshelvesdto.UpsertMyRootShelfPermissionResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	permission, err := enums.ConvertStringToAccessControlPermission(requestDto.Body.Permission)
	if err != nil {
		return nil, exceptions.InvalidInput("RootShelf").WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	return s.saveMyRootShelfPermission(ctx, actorUserId, requestDto.Param.RootShelfId, requestDto.Param.UserPublicId, *permission, nil)
}

func (s *RootShelfService) UpsertMyRootShelfPermissions(
	ctx context.Context, requestDto *rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto,
) (*rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	userPublicIds := make([]uuid.UUID, len(requestDto.Body.Permissions))
	permissionByPublicId := make(map[uuid.UUID]enums.AccessControlPermission, len(requestDto.Body.Permissions))
	for index, input := range requestDto.Body.Permissions {
		permission, err := enums.ConvertStringToAccessControlPermission(input.Permission)
		if err != nil {
			return nil, exceptions.InvalidInput("RootShelf").WithOrigin(err)
		}
		if *permission == enums.AccessControlPermission_Owner {
			return nil, exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
		}
		if _, exists := permissionByPublicId[input.UserPublicId]; exists {
			return nil, exceptions.New(
				"InvalidRequest",
				"RootShelf",
				"ValidateRequest",
				"Root shelf request is invalid",
				http.StatusBadRequest,
			)
		}

		userPublicIds[index] = input.UserPublicId
		permissionByPublicId[input.UserPublicId] = *permission
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"Manage",
			"Failed to begin the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Param.RootShelfId,
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

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", userPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}
	if len(targetUsers) != len(userPublicIds) {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		)
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
			return nil, exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			(permission == enums.AccessControlPermission_Admin ||
				existingPermissionByUserId[userId] == enums.AccessControlPermission_Admin) {
			tx.Rollback()
			return nil, exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
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
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	updatedPermissionByUserId := make(map[uuid.UUID]schemas.UsersToShelves, len(updatedPermissions))
	for _, updatedPermission := range updatedPermissions {
		updatedPermissionByUserId[updatedPermission.UserId] = updatedPermission
	}

	responseDtos := make([]rootshelvesdto.RootShelfPermissionResponseDto, len(userIds))
	for index, userId := range userIds {
		user := userById[userId]
		updatedPermission := updatedPermissionByUserId[userId]
		responseDtos[index] = rootshelvesdto.RootShelfPermissionResponseDto{
			UserPublicId: user.PublicId,
			Permission:   updatedPermission.Permission.String(),
			UpdatedAt:    updatedPermission.UpdatedAt,
			CreatedAt:    updatedPermission.CreatedAt,
		}
	}

	return &rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto{
		Permissions: responseDtos,
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelfPermission(
	ctx context.Context, requestDto *rootshelvesdto.UpdateMyRootShelfPermissionRequestDto,
) (*rootshelvesdto.UpdateMyRootShelfPermissionResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	permission, err := enums.ConvertStringToAccessControlPermission(requestDto.Body.Permission)
	if err != nil {
		return nil, exceptions.InvalidInput("RootShelf").WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	requireExisting := true
	return s.saveMyRootShelfPermission(ctx, actorUserId, requestDto.Param.RootShelfId, requestDto.Param.UserPublicId, *permission, &requireExisting)
}

func (s *RootShelfService) TransferMyRootShelfOwnership(
	ctx context.Context,
	requestDto *rootshelvesdto.TransferMyRootShelfOwnershipRequestDto,
) (*rootshelvesdto.TransferMyRootShelfOwnershipResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"TransferMyRootShelfOwnership",
			"Failed to begin the root shelf ownership transfer transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	rootShelf, permission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Param.RootShelfId,
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
	if permission != enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}

	var actorUser schemas.User
	if result := tx.Select("id, public_id").Where("id = ?", actorUserId).First(&actorUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}
	var targetUser schemas.User
	if result := tx.Select("id, public_id").Where("public_id = ?", requestDto.Body.TargetUserPublicId).First(&targetUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}
	if targetUser.Id == actorUserId {
		tx.Rollback()
		return nil, exceptions.New(
			"NoChanges",
			"RootShelf",
			"Manage",
			"No root shelf changes were applied",
			http.StatusNotModified,
		)
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
		return nil, exceptions.New(
			"NoChanges",
			"RootShelf",
			"Manage",
			"No root shelf changes were applied",
			http.StatusNotModified,
		)
	}

	var accounts []schemas.UserAccount
	result := tx.
		Clauses(clause.Locking{Strength: options.LockingStrengthUpdate}).
		Where("user_id IN ?", []uuid.UUID{actorUserId, targetUser.Id}).
		Order("user_id").
		Find(&accounts)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToUpdate",
			"RootShelf",
			"Manage",
			"Failed to update the root shelf",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}
	if len(accounts) != 2 {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		)
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
		return nil, exceptions.New(
			"FailedToUpdate",
			"RootShelf",
			"Manage",
			"Failed to update the root shelf",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
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
		return nil, exceptions.New(
			"QueryFailed",
			"BlockPack",
			"ManageRootShelf",
			"Failed to retrieve root shelf block packs",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	if _, exception = s.usersToShelvesRepository.UpdateOne(
		rootShelf.Id,
		actorUserId,
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
		return nil, exceptions.New(
			"FailedToUpdate",
			"RootShelf",
			"Manage",
			"Failed to update the root shelf",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"RootShelf",
			"Manage",
			"Root shelf was not found",
			http.StatusNotFound,
		)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &rootshelvesdto.TransferMyRootShelfOwnershipResponseDto{
		RootShelfId:               rootShelf.Id,
		PreviousOwnerUserPublicId: actorUser.PublicId,
		NewOwnerUserPublicId:      targetUser.PublicId,
		UpdatedAt:                 newOwnerMembership.UpdatedAt,
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelfPermission(
	ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionRequestDto,
) (*rootshelvesdto.DeleteMyRootShelfPermissionResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Param.RootShelfId,
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
	result := tx.
		Model(&schemas.User{}).
		Where("public_id = ?", requestDto.Param.UserPublicId).
		First(&targetUser)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}

	targetPermission, exception := s.usersToShelvesRepository.GetOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}
	if actorPermission != enums.AccessControlPermission_Owner &&
		targetPermission.Permission == enums.AccessControlPermission_Admin {
		tx.Rollback()
		return nil, exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}

	exception = s.usersToShelvesRepository.DeleteOne(
		rootShelf.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &rootshelvesdto.DeleteMyRootShelfPermissionResponseDto{}, nil
}

func (s *RootShelfService) DeleteMyRootShelfPermissions(
	ctx context.Context, requestDto *rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto,
) (*rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}

	userPublicIdSet := make(map[uuid.UUID]struct{}, len(requestDto.Body.UserPublicIds))
	for _, userPublicId := range requestDto.Body.UserPublicIds {
		if _, exists := userPublicIdSet[userPublicId]; exists {
			return nil, exceptions.New(
				"InvalidRequest",
				"RootShelf",
				"ValidateRequest",
				"Root shelf request is invalid",
				http.StatusBadRequest,
			)
		}

		userPublicIdSet[userPublicId] = struct{}{}
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"Manage",
			"Failed to begin the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}

	rootShelf, actorPermission, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
		requestDto.Param.RootShelfId,
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

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", requestDto.Body.UserPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}
	if len(targetUsers) != len(requestDto.Body.UserPublicIds) {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		)
	}

	userIdByPublicId := make(map[uuid.UUID]uuid.UUID, len(targetUsers))
	for _, targetUser := range targetUsers {
		userIdByPublicId[targetUser.PublicId] = targetUser.Id
	}

	userIds := make([]uuid.UUID, len(requestDto.Body.UserPublicIds))
	for index, userPublicId := range requestDto.Body.UserPublicIds {
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
		return nil, exception
	}
	if len(targetPermissions) != len(userIds) {
		tx.Rollback()
		return nil, exceptions.New(
			"NotFound",
			"RootShelf",
			"Manage",
			"Root shelf was not found",
			http.StatusNotFound,
		)
	}

	for _, targetPermission := range targetPermissions {
		if targetPermission.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return nil, exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			targetPermission.Permission == enums.AccessControlPermission_Admin {
			tx.Rollback()
			return nil, exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
		}
	}

	exception = s.usersToShelvesRepository.DeleteMany(
		rootShelf.Id,
		userIds,
		options.WithTransactionDB(tx),
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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	for _, targetUser := range targetUsers {
		if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(targetUser.Id, blockPackIds); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
		}
	}

	return &rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto{}, nil
}

func (s *RootShelfService) LeaveMyRootShelf(
	ctx context.Context, requestDto *rootshelvesdto.LeaveMyRootShelfRequestDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"LeaveMyRootShelf",
			"Failed to begin the root shelf leave transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	blockPacks, exception := s.blockPackRepository.GetManyByRootShelfIds(
		[]uuid.UUID{requestDto.Param.RootShelfId},
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
		requestDto.Param.RootShelfId,
		actorUserId,
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
		return exceptions.New(
			"PermissionDenied",
			"RootShelf",
			"ManagePermission",
			"You do not have permission to manage this root shelf",
			http.StatusBadRequest,
		)
	}
	if exception = s.usersToShelvesRepository.DeleteOne(
		rootShelf.Id,
		actorUserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(actorUserId, blockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return nil
}

func (s *RootShelfService) LeaveMyRootShelves(
	ctx context.Context, requestDto *rootshelvesdto.LeaveMyRootShelvesRequestDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return apiexceptions.Shelf.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return exception
	}
	rootShelfIdSet := make(map[uuid.UUID]struct{}, len(requestDto.Body.RootShelves))
	rootShelfIds := make([]uuid.UUID, len(requestDto.Body.RootShelves))
	for index, rootShelfRequestDto := range requestDto.Body.RootShelves {
		if _, exists := rootShelfIdSet[rootShelfRequestDto.RootShelfId]; exists {
			return exceptions.New(
				"InvalidRequest",
				"RootShelf",
				"ValidateRequest",
				"Root shelf request is invalid",
				http.StatusBadRequest,
			)
		}
		rootShelfIdSet[rootShelfRequestDto.RootShelfId] = struct{}{}
		rootShelfIds[index] = rootShelfRequestDto.RootShelfId
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.New(
			"TransactionBeginFailed",
			"RootShelf",
			"LeaveMyRootShelves",
			"Failed to begin the root shelf leave transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
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
		actorUserId,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if len(relations) != len(rootShelfIds) {
		tx.Rollback()
		return exceptions.New(
			"NotFound",
			"RootShelf",
			"Manage",
			"Root shelf was not found",
			http.StatusNotFound,
		)
	}
	for _, relation := range relations {
		if relation.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return exceptions.New(
				"PermissionDenied",
				"RootShelf",
				"ManagePermission",
				"You do not have permission to manage this root shelf",
				http.StatusBadRequest,
			)
		}
	}

	if exception = s.usersToShelvesRepository.DeleteManyByRootShelfIdsAndUserId(
		rootShelfIds,
		actorUserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.New(
			"TransactionCommitFailed",
			"RootShelf",
			"Manage",
			"Failed to commit the root shelf transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(actorUserId, blockPackIds); err != nil {
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
			return nil, exceptions.New(
				"CursorDecodeFailed",
				"Search",
				"SearchPrivateRootShelves",
				"Failed to decode the search cursor",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
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
		return nil, exceptions.New(
			"NotFound",
			"RootShelf",
			"Manage",
			"Root shelf was not found",
			http.StatusNotFound,
		).WithOrigin(err)
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
			return nil, exceptions.New(
				"NotFound",
				"User",
				"ResolveUser",
				"User was not found",
				http.StatusNotFound,
			).WithOrigin(err)
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
			return nil, exceptions.New(
				"CursorEncodeFailed",
				"Search",
				"SearchPrivateRootShelves",
				"Failed to encode the search cursor",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.New(
				"CursorEncodingFailed",
				"Search",
				"SearchPrivateRootShelves",
				"Failed to encode the search cursor",
				http.StatusInternalServerError,
				true,
			)
		}

		privateRootShelf := shelf.RootShelf.ToPrivateRootShelf(shelf.Permission)
		owner, exists := publicUsersById[shelf.OwnerId]
		if !exists {
			return nil, exceptions.New(
				"NotFound",
				"User",
				"ResolveUser",
				"User was not found",
				http.StatusNotFound,
			)
		}

		privateRootShelf.Owner = owner
		for _, usersToShelf := range shelf.UsersToShelves {
			if usersToShelf.UserId == shelf.OwnerId {
				continue
			}

			sharer, exists := publicUsersById[usersToShelf.UserId]
			if !exists {
				return nil, exceptions.New(
					"NotFound",
					"User",
					"ResolveUser",
					"User was not found",
					http.StatusNotFound,
				)
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
