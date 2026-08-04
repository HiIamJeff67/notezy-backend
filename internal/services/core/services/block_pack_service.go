package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	eventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	blockpacksql "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/sqls/block_pack"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
	validator "github.com/go-playground/validator/v10"
)

type BlockPackServiceInterface interface {
	GetMyBlockPackById(ctx context.Context, requestDto *blockpacksdto.GetMyBlockPackByIdRequestDto) (*blockpacksdto.GetMyBlockPackByIdResponseDto, *exceptions.Exception)
	GetMyBlockPackAndItsParentById(ctx context.Context, requestDto *blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto) (*blockpacksdto.GetMyBlockPackAndItsParentByIdResponseDto, *exceptions.Exception)
	GetMyBlockPacksByParentSubShelfId(ctx context.Context, requestDto *blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto) (*blockpacksdto.GetMyBlockPacksByParentSubShelfIdResponseDto, *exceptions.Exception)
	GetAllMyBlockPacksByRootShelfId(ctx context.Context, requestDto *blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto) (*blockpacksdto.GetAllMyBlockPacksByRootShelfIdResponseDto, *exceptions.Exception)
	CreateBlockPack(ctx context.Context, requestDto *blockpacksdto.CreateBlockPackRequestDto) (*blockpacksdto.CreateBlockPackResponseDto, *exceptions.Exception)
	CreateBlockPacks(ctx context.Context, requestDto *blockpacksdto.CreateBlockPacksRequestDto) (*blockpacksdto.CreateBlockPacksResponseDto, *exceptions.Exception)
	UpdateMyBlockPackById(ctx context.Context, requestDto *blockpacksdto.UpdateMyBlockPackByIdRequestDto) (*blockpacksdto.UpdateMyBlockPackByIdResponseDto, *exceptions.Exception)
	UpdateMyBlockPacksByIds(ctx context.Context, requestDto *blockpacksdto.UpdateMyBlockPacksByIdsRequestDto) (*blockpacksdto.UpdateMyBlockPacksByIdsResponseDto, *exceptions.Exception)
	MoveMyBlockPackByParentSubShelfId(ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto) (*blockpacksdto.MoveMyBlockPackByParentSubShelfIdResponseDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfId(ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto) (*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdResponseDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfIds(ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto) (*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsResponseDto, *exceptions.Exception)
	RestoreMyBlockPackById(ctx context.Context, requestDto *blockpacksdto.RestoreMyBlockPackByIdRequestDto) (*blockpacksdto.RestoreMyBlockPackByIdResponseDto, *exceptions.Exception)
	RestoreMyBlockPacksByIds(ctx context.Context, requestDto *blockpacksdto.RestoreMyBlockPacksByIdsRequestDto) (*blockpacksdto.RestoreMyBlockPacksByIdsResponseDto, *exceptions.Exception)
	DeleteMyBlockPackById(ctx context.Context, requestDto *blockpacksdto.DeleteMyBlockPackByIdRequestDto) (*blockpacksdto.DeleteMyBlockPackByIdResponseDto, *exceptions.Exception)
	DeleteMyBlockPacksByIds(ctx context.Context, requestDto *blockpacksdto.DeleteMyBlockPacksByIdsRequestDto) (*blockpacksdto.DeleteMyBlockPacksByIdsResponseDto, *exceptions.Exception)

	SearchPrivateBlockPacks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockPackInput) (*gqlmodels.SearchBlockPackConnection, *exceptions.Exception)
}

type BlockPackService struct {
	validator           *validator.Validate
	db                  *gorm.DB
	blockPackScope      scopes.BlockPackScopeInterface
	subShelfRepository  repositories.SubShelfRepositoryInterface
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewBlockPackService(
	validator *validator.Validate,
	db *gorm.DB,
	blockPackScope scopes.BlockPackScopeInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) BlockPackServiceInterface {
	return &BlockPackService{
		validator:           validator,
		db:                  db,
		blockPackScope:      blockPackScope,
		subShelfRepository:  subShelfRepository,
		blockPackRepository: blockPackRepository,
	}
}

/* ============================== Main Methods ============================== */

func (s *BlockPackService) GetMyBlockPackById(
	ctx context.Context, requestDto *blockpacksdto.GetMyBlockPackByIdRequestDto,
) (*blockpacksdto.GetMyBlockPackByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.IsDeleted != nil {
		if *requestDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	blockPack, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		requestDto.Param.BlockPackId,
		actorUserId,
		[]schemas.BlockPackRelation{schemas.BlockPackRelation_YjsDocument},
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := &blockpacksdto.GetMyBlockPackByIdResponseDto{
		Id:                  blockPack.Id,
		ParentSubShelfId:    blockPack.ParentSubShelfId,
		Name:                blockPack.Name,
		Icon:                blockPack.Icon.ToContractable(),
		HeaderBackgroundURL: blockPack.HeaderBackgroundURL,
		BlockCount:          blockPack.BlockCount,
		DeletedAt:           blockPack.DeletedAt,
		UpdatedAt:           blockPack.UpdatedAt,
		CreatedAt:           blockPack.CreatedAt,
	}
	if blockPack.YjsDocument != nil {
		resDto.LastUpdateSequence = blockPack.YjsDocument.LastUpdateSequence
		resDto.CompactedUntilSequence = blockPack.YjsDocument.CompactedUntilSequence
		resDto.ProjectedUntilSequence = blockPack.YjsDocument.ProjectedUntilSequence
		resDto.IsProjectionCurrent = blockPack.YjsDocument.ProjectedUntilSequence >= blockPack.YjsDocument.LastUpdateSequence
	}

	return resDto, nil
}

func (s *BlockPackService) GetMyBlockPackAndItsParentById(
	ctx context.Context, requestDto *blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto,
) (*blockpacksdto.GetMyBlockPackAndItsParentByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.IsDeleted != nil {
		if *requestDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := blockpacksdto.GetMyBlockPackAndItsParentByIdResponseDto{}
	err := db.Raw(blockpacksql.GetMyBlockPackAndItsParentByIdSQL,
		requestDto.Param.BlockPackId, actorUserId, pg.Array(allowedPermissions), onlyDeleted,
	).Row().
		Scan(&resDto.Id,
			&resDto.Name,
			&resDto.Icon,
			&resDto.HeaderBackgroundURL,
			&resDto.BlockCount,
			&resDto.LastUpdateSequence,
			&resDto.CompactedUntilSequence,
			&resDto.ProjectedUntilSequence,
			&resDto.IsProjectionCurrent,
			&resDto.DeletedAt,
			&resDto.UpdatedAt,
			&resDto.CreatedAt,
			&resDto.RootShelfId,
			&resDto.Permission,
			&resDto.ParentSubShelfId,
			&resDto.ParentSubShelfName,
			&resDto.ParentSubShelfPrevSubShelfId,
			&resDto.ParentSubShelfPath,
			&resDto.ParentSubShelfDeletedAt,
			&resDto.ParentSubShelfUpdatedAt,
			&resDto.ParentSubShelfCreatedAt)
	if err != nil {
		return nil, apiexceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetMyBlockPacksByParentSubShelfId(
	ctx context.Context, requestDto *blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto,
) (*blockpacksdto.GetMyBlockPacksByParentSubShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.AreDeleted != nil {
		if *requestDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := blockpacksdto.GetMyBlockPacksByParentSubShelfIdResponseDto{}
	result := db.Model(&schemas.BlockPack{}).
		Select(`
			"BlockPackTable".*,
			COALESCE(ydoc.last_update_sequence, 0) AS last_update_sequence,
			COALESCE(ydoc.compacted_until_sequence, 0) AS compacted_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) AS projected_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) >= COALESCE(ydoc.last_update_sequence, 0) AS is_projection_current
		`).
		Joins(`LEFT JOIN "BlockPackYjsDocumentTable" ydoc ON ydoc.block_pack_id = "BlockPackTable".id AND ydoc.deleted_at IS NULL`).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "BlockPackTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			requestDto.Param.ParentSubShelfId,
			actorUserId,
			allowedPermissions,
		).Scopes(scopes.NewBlockPackScope().FilterOnlyDeleted(onlyDeleted)).
		Order("name ASC").
		Limit(int(constants.MaxBlockPackOfSubShelf)).
		Scan(&resDto)
	if err := result.Error; err != nil {
		return nil, apiexceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetAllMyBlockPacksByRootShelfId(
	ctx context.Context, requestDto *blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto,
) (*blockpacksdto.GetAllMyBlockPacksByRootShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.AreDeleted != nil {
		if *requestDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := blockpacksdto.GetAllMyBlockPacksByRootShelfIdResponseDto{}
	result := db.Model(&schemas.BlockPack{}).
		Select(`
			"BlockPackTable".*,
			COALESCE(ydoc.last_update_sequence, 0) AS last_update_sequence,
			COALESCE(ydoc.compacted_until_sequence, 0) AS compacted_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) AS projected_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) >= COALESCE(ydoc.last_update_sequence, 0) AS is_projection_current
		`).
		Joins(`LEFT JOIN "BlockPackYjsDocumentTable" ydoc ON ydoc.block_pack_id = "BlockPackTable".id AND ydoc.deleted_at IS NULL`).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "BlockPackTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			requestDto.Param.RootShelfId, actorUserId, allowedPermissions,
		).Scopes(scopes.NewBlockPackScope().FilterOnlyDeleted(onlyDeleted)).
		Limit(int(constants.MaxBlockPackOfRootShelf)).
		Order("name ASC").
		Scan(&resDto)
	if err := result.Error; err != nil {
		return nil, apiexceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) CreateBlockPack(
	ctx context.Context, requestDto *blockpacksdto.CreateBlockPackRequestDto,
) (*blockpacksdto.CreateBlockPackResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
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

	newBlockPackId, exception := s.blockPackRepository.CreateOneBySubShelfId(
		requestDto.Body.ParentSubShelfId,
		actorUserId,
		inputs.CreateBlockPackInput{
			Id:                  requestDto.Body.Id,
			Name:                requestDto.Body.Name,
			Icon:                (*enums.SupportedIcon)(requestDto.Body.Icon).ToStorable(),
			HeaderBackgroundURL: requestDto.Body.HeaderBackgroundURL,
		},
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	document := schemas.BlockPackYjsDocument{BlockPackId: *newBlockPackId}
	if err := tx.Create(&document).Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	return &blockpacksdto.CreateBlockPackResponseDto{
		Id:        *newBlockPackId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) CreateBlockPacks(
	ctx context.Context, requestDto *blockpacksdto.CreateBlockPacksRequestDto,
) (*blockpacksdto.CreateBlockPacksResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
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

	input := make([]inputs.CreateBlockPackBySubShelfIdInput, len(requestDto.Body.CreatedBlockPacks))
	for index, createdBlockPack := range requestDto.Body.CreatedBlockPacks {
		input[index] = inputs.CreateBlockPackBySubShelfIdInput{
			Id:                  createdBlockPack.Id,
			ParentSubShelfId:    createdBlockPack.ParentSubShelfId,
			Name:                createdBlockPack.Name,
			Icon:                (*enums.SupportedIcon)(createdBlockPack.Icon).ToStorable(),
			HeaderBackgroundURL: createdBlockPack.HeaderBackgroundURL,
		}
	}
	newBlockPackIds, exception := s.blockPackRepository.CreateManyBySubShelfIds(
		actorUserId,
		input,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	documents := make([]schemas.BlockPackYjsDocument, len(newBlockPackIds))
	for index, newBlockPackId := range newBlockPackIds {
		documents[index] = schemas.BlockPackYjsDocument{BlockPackId: newBlockPackId}
	}
	if err := tx.CreateInBatches(&documents, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()

		return nil, apiexceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		return nil, apiexceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	return &blockpacksdto.CreateBlockPacksResponseDto{
		Ids:       newBlockPackIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) UpdateMyBlockPackById(
	ctx context.Context, requestDto *blockpacksdto.UpdateMyBlockPackByIdRequestDto,
) (*blockpacksdto.UpdateMyBlockPackByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	blockPack, exception := s.blockPackRepository.UpdateOneById(
		requestDto.Param.BlockPackId,
		actorUserId,
		inputs.PartialUpdateBlockPackInput{
			Values: inputs.UpdateBlockPackInput{
				Name:                requestDto.Body.Values.Name,
				Icon:                (*enums.SupportedIcon)(requestDto.Body.Values.Icon).ToStorable(),
				HeaderBackgroundURL: requestDto.Body.Values.HeaderBackgroundURL,
			},
			SetNull: requestDto.Body.SetNull,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &blockpacksdto.UpdateMyBlockPackByIdResponseDto{
		UpdatedAt: blockPack.UpdatedAt,
	}, nil
}

func (s *BlockPackService) UpdateMyBlockPacksByIds(
	ctx context.Context, requestDto *blockpacksdto.UpdateMyBlockPacksByIdsRequestDto,
) (*blockpacksdto.UpdateMyBlockPacksByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, len(requestDto.Body.UpdatedBlockPacks))
	for index, updatedBlockPack := range requestDto.Body.UpdatedBlockPacks {
		input[index] = inputs.UpdateBlockPackByIdInput{
			Id: updatedBlockPack.BlockPackId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
				Values: inputs.UpdateBlockPackInput{
					Name:                updatedBlockPack.Values.Name,
					Icon:                (*enums.SupportedIcon)(updatedBlockPack.Values.Icon).ToStorable(),
					HeaderBackgroundURL: updatedBlockPack.Values.HeaderBackgroundURL,
				},
			},
		}
	}
	exception = s.blockPackRepository.UpdateManyByIds(
		actorUserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &blockpacksdto.UpdateMyBlockPacksByIdsResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPackByParentSubShelfId(
	ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto,
) (*blockpacksdto.MoveMyBlockPackByParentSubShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
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
			"BlockPack",
			"MoveMyBlockPackByParentSubShelfId",
			"Failed to begin the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	_, exception = s.blockPackRepository.UpdateOneById(
		requestDto.Body.BlockPackId,
		actorUserId,
		inputs.PartialUpdateBlockPackInput{
			Values: inputs.UpdateBlockPackInput{
				ParentSubShelfId: &requestDto.Body.DestinationParentSubShelfId,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if err := repositories.NewOutboxEventRepository().EnqueueBlockPackAccessRevocations(
		tx,
		requestDto.Body.BlockPackId.String(),
		[]uuid.UUID{requestDto.Body.BlockPackId},
		nil,
		eventscontract.BlockPackAccessRevocationReason_PermissionRevoked,
	); err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToCreate",
			"Outbox",
			"MoveMyBlockPackByParentSubShelfId",
			"Failed to create lifecycle outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"MoveMyBlockPackByParentSubShelfId",
			"Failed to commit the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &blockpacksdto.MoveMyBlockPackByParentSubShelfIdResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPacksByParentSubShelfId(
	ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto,
) (*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, len(requestDto.Body.BlockPackIds))
	for index, blockPackId := range requestDto.Body.BlockPackIds {
		input[index] = inputs.UpdateBlockPackByIdInput{
			Id: blockPackId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
				Values: inputs.UpdateBlockPackInput{
					ParentSubShelfId: &requestDto.Body.DestinationParentSubShelfId,
				},
			},
		}
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"BlockPack",
			"MoveMyBlockPacksByParentSubShelfId",
			"Failed to begin the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	exception = s.blockPackRepository.UpdateManyByIds(
		actorUserId,
		input,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if err := repositories.NewOutboxEventRepository().EnqueueBlockPackAccessRevocations(
		tx,
		"block-pack-bulk-move",
		requestDto.Body.BlockPackIds,
		nil,
		eventscontract.BlockPackAccessRevocationReason_PermissionRevoked,
	); err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToCreate",
			"Outbox",
			"MoveMyBlockPacksByParentSubShelfId",
			"Failed to create lifecycle outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"MoveMyBlockPacksByParentSubShelfId",
			"Failed to commit the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &blockpacksdto.MoveMyBlockPacksByParentSubShelfIdResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPacksByParentSubShelfIds(
	ctx context.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto,
) (*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, 0)
	for _, movedBlockPack := range requestDto.Body.MovedBlockPacks {
		for _, blockPackId := range movedBlockPack.BlockPackIds {
			input = append(input, inputs.UpdateBlockPackByIdInput{
				Id: blockPackId,
				PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
					Values: inputs.UpdateBlockPackInput{
						ParentSubShelfId: &movedBlockPack.DestinationParentSubShelfId,
					},
				},
			})
		}
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"BlockPack",
			"MoveMyBlockPacksByParentSubShelfIds",
			"Failed to begin the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	if exception = s.blockPackRepository.UpdateManyByIds(
		actorUserId,
		input,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}
	blockPackIds := make([]uuid.UUID, len(input))
	for index, movedBlockPack := range input {
		blockPackIds[index] = movedBlockPack.Id
	}
	if err := repositories.NewOutboxEventRepository().EnqueueBlockPackAccessRevocations(
		tx,
		"block-pack-multi-parent-move",
		blockPackIds,
		nil,
		eventscontract.BlockPackAccessRevocationReason_PermissionRevoked,
	); err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToCreate",
			"Outbox",
			"MoveMyBlockPacksByParentSubShelfIds",
			"Failed to create lifecycle outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"MoveMyBlockPacksByParentSubShelfIds",
			"Failed to commit the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPackById(
	ctx context.Context, requestDto *blockpacksdto.RestoreMyBlockPackByIdRequestDto,
) (*blockpacksdto.RestoreMyBlockPackByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredBlockPack, exception := s.blockPackRepository.RestoreSoftDeletedOneById(
		requestDto.Param.BlockPackId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &blockpacksdto.RestoreMyBlockPackByIdResponseDto{
		Id:                  restoredBlockPack.Id,
		ParentSubShelfId:    restoredBlockPack.ParentSubShelfId,
		Name:                restoredBlockPack.Name,
		Icon:                restoredBlockPack.Icon.ToContractable(),
		HeaderBackgroundURL: restoredBlockPack.HeaderBackgroundURL,
		BlockCount:          restoredBlockPack.BlockCount,
		DeletedAt:           restoredBlockPack.DeletedAt,
		UpdatedAt:           restoredBlockPack.UpdatedAt,
		CreatedAt:           restoredBlockPack.CreatedAt,
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPacksByIds(
	ctx context.Context, requestDto *blockpacksdto.RestoreMyBlockPacksByIdsRequestDto,
) (*blockpacksdto.RestoreMyBlockPacksByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredBlockPacks, exception := s.blockPackRepository.RestoreSoftDeletedManyByIds(
		requestDto.Body.BlockPackIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := blockpacksdto.RestoreMyBlockPacksByIdsResponseDto{}
	for _, restoredBlockPack := range restoredBlockPacks {
		resDto = append(resDto, blockpacksdto.RestoreMyBlockPackByIdResponseDto{
			Id:                  restoredBlockPack.Id,
			ParentSubShelfId:    restoredBlockPack.ParentSubShelfId,
			Name:                restoredBlockPack.Name,
			Icon:                restoredBlockPack.Icon.ToContractable(),
			HeaderBackgroundURL: restoredBlockPack.HeaderBackgroundURL,
			BlockCount:          restoredBlockPack.BlockCount,
			DeletedAt:           restoredBlockPack.DeletedAt,
			UpdatedAt:           restoredBlockPack.UpdatedAt,
			CreatedAt:           restoredBlockPack.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *BlockPackService) DeleteMyBlockPackById(
	ctx context.Context, requestDto *blockpacksdto.DeleteMyBlockPackByIdRequestDto,
) (*blockpacksdto.DeleteMyBlockPackByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
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
			"BlockPack",
			"DeleteMyBlockPackById",
			"Failed to begin the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	if exception = s.blockPackRepository.SoftDeleteOneById(
		requestDto.Param.BlockPackId,
		actorUserId,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if err := repositories.NewOutboxEventRepository().EnqueueBlockPackAccessRevocations(
		tx,
		requestDto.Param.BlockPackId.String(),
		[]uuid.UUID{requestDto.Param.BlockPackId},
		nil,
		eventscontract.BlockPackAccessRevocationReason_ResourceUnavailable,
	); err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToCreate",
			"Outbox",
			"DeleteMyBlockPackById",
			"Failed to create lifecycle outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"DeleteMyBlockPackById",
			"Failed to commit the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &blockpacksdto.DeleteMyBlockPackByIdResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) DeleteMyBlockPacksByIds(
	ctx context.Context, requestDto *blockpacksdto.DeleteMyBlockPacksByIdsRequestDto,
) (*blockpacksdto.DeleteMyBlockPacksByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.BlockPack.InvalidDto().WithOrigin(err)
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
			"BlockPack",
			"DeleteMyBlockPacksByIds",
			"Failed to begin the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}
	if exception = s.blockPackRepository.SoftDeleteManyByIds(
		requestDto.Body.BlockPackIds,
		actorUserId,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if err := repositories.NewOutboxEventRepository().EnqueueBlockPackAccessRevocations(
		tx,
		"block-pack-bulk-delete",
		requestDto.Body.BlockPackIds,
		nil,
		eventscontract.BlockPackAccessRevocationReason_ResourceUnavailable,
	); err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToCreate",
			"Outbox",
			"DeleteMyBlockPacksByIds",
			"Failed to create lifecycle outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"DeleteMyBlockPacksByIds",
			"Failed to commit the block pack transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &blockpacksdto.DeleteMyBlockPacksByIdsResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL BlockPack ============================== */

func (s *BlockPackService) SearchPrivateBlockPacks(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockPackInput,
) (*gqlmodels.SearchBlockPackConnection, *exceptions.Exception) {
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

	query := db.Model(&schemas.BlockPack{}).
		Select(`"BlockPackTable".*`).
		Joins(`INNER JOIN "SubShelfTable" ss ON ss.id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.blockPackScope.FilterOnlyDeleted(onlyDeleted))

	if gqlInput.ParentSubShelfID != nil {
		query = query.Where(`"BlockPackTable".parent_sub_shelf_id = ?`, *gqlInput.ParentSubShelfID)
	}

	if gqlInput.RootShelfID != nil {
		query = query.Where("ss.root_shelf_id = ?", *gqlInput.RootShelfID)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			`"BlockPackTable".name ILIKE ?`,
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchBlockPackCursorFields](*gqlInput.After)
		if err != nil {
			return nil, apiexceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"BlockPackTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchBlockPackSortByName:
			query = query.Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByBlockCount:
			query = query.Order(`"BlockPackTable".block_count ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByLastUpdate:
			query = query.Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByCreatedAt:
			query = query.Order(`"BlockPackTable".created_at ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending)
		default:
			query = query.Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var blockPacks []schemas.BlockPack
	if err := query.Scopes(s.blockPackScope.IncludePreloads(
		[]schemas.BlockPackRelation{
			schemas.BlockPackRelation_Blocks,
		},
	)).Find(&blockPacks).Error; err != nil {
		return nil, apiexceptions.BlockPack.NotFound().WithOrigin(err)
	}

	hasNextPage := len(blockPacks) > limit
	searchEdges := make([]*gqlmodels.SearchBlockPackEdge, len(blockPacks))

	for index, blockPack := range blockPacks {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchBlockPackCursorFields]{
			Fields: gqlmodels.SearchBlockPackCursorFields{
				ID: blockPack.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, apiexceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, apiexceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchBlockPackEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                blockPack.ToPrivateBlockPack(),
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

	return &gqlmodels.SearchBlockPackConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
