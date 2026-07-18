package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
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

type BlockServiceInterface interface {
	GetMyBlockById(ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception)
	GetMyBlocksByIds(ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception)
	GetMyBlocksByBlockPackId(ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception)
	Apply(ctx context.Context, blockPackId uuid.UUID, input dtos.ApplyBlockProjectionInput) (*dtos.ApplyBlockProjectionResult, error)
	ApplyMany(ctx context.Context, inputs []dtos.ApplyBlockProjectionDocumentInput) (dtos.ApplyBlockProjectionDocumentResult, error)

	SearchPrivateBlocks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockInput) (*gqlmodels.SearchBlockConnection, *exceptions.Exception)
}

type BlockService struct {
	db                   *gorm.DB
	blockScope           scopes.BlockScopeInterface
	blockPackScope       scopes.BlockPackScopeInterface
	subShelfScope        scopes.SubShelfScopeInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
	editableBlockAdapter adapters.EditableBlockAdapterInterface
}

func NewBlockService(
	db *gorm.DB,
	blockScope scopes.BlockScopeInterface,
	blockPackScope scopes.BlockPackScopeInterface,
	subShelfScope scopes.SubShelfScopeInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockServiceInterface {
	return &BlockService{
		db:                   db,
		blockScope:           blockScope,
		blockPackScope:       blockPackScope,
		subShelfScope:        subShelfScope,
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
		editableBlockAdapter: adapters.NewEditableBlockAdapter(),
	}
}

func NewBlockProjectionService(db *gorm.DB) BlockServiceInterface {
	blockScope := scopes.NewBlockScope()
	blockPackScope := scopes.NewBlockPackScope()
	subShelfScope := scopes.NewSubShelfScope()

	return NewBlockService(
		db,
		blockScope,
		blockPackScope,
		subShelfScope,
		repositories.NewBlockPackRepository(blockPackScope),
		repositories.NewBlockRepository(blockScope),
	)
}

func (s *BlockService) GetMyBlockById(
	ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto,
) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	block, exception := s.blockRepository.GetOneById(
		reqDto.Param.BlockId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	res := dtos.GetMyBlockByIdResDto{
		Id:            block.Id,
		BlockPackId:   block.BlockPackId,
		ParentBlockId: block.ParentBlockId,
		PrevBlockId:   block.PrevBlockId,
		NextBlockId:   block.NextBlockId,
		Type:          block.Type,
		Props:         block.Props,
		Content:       block.Content,
		UpdatedAt:     block.UpdatedAt,
		CreatedAt:     block.CreatedAt,
	}

	return &res, nil
}

func (s *BlockService) GetMyBlocksByIds(
	ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto,
) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Param.BlockIds,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	res := make(dtos.GetMyBlocksByIdsResDto, len(blocks))
	for index, block := range blocks {
		res[index] = dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) GetMyBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	if !s.blockPackRepository.HasPermission(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Block.NoPermission("get the block pack of blocks")
	}

	var blocks []schemas.Block
	if err := db.Model(&schemas.Block{}).
		Where("block_pack_id = ?", reqDto.Param.BlockPackId).
		Order("created_at ASC").
		Order("id ASC").
		Find(&blocks).Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	res := make(dtos.GetMyBlocksByBlockPackIdResDto, len(blocks))
	for index, block := range blocks {
		res[index] = dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) Apply(
	ctx context.Context,
	blockPackId uuid.UUID,
	input dtos.ApplyBlockProjectionInput,
) (*dtos.ApplyBlockProjectionResult, error) {
	if blockPackId == uuid.Nil {
		return nil, fmt.Errorf("block projection requires a block pack id")
	}
	if input.SchemaId != constants.YjsBlockPackSchemaId ||
		input.SchemaVersion != constants.YjsBlockPackSchemaVersion {
		return nil, fmt.Errorf("block projection source schema is not supported")
	}
	if input.ProjectedSequence < 0 {
		return nil, fmt.Errorf("block projection target update sequence must not be negative")
	}

	flattenedBlocks, _, exception := s.editableBlockAdapter.FlattenManyToRaw(input.Blocks)
	if exception != nil {
		return nil, fmt.Errorf("failed to flatten block projection: %w", exception)
	}

	blockIds := make([]uuid.UUID, len(flattenedBlocks))
	projectedBlocks := make([]schemas.Block, len(flattenedBlocks))
	for index, flattenedBlock := range flattenedBlocks {
		blockIds[index] = flattenedBlock.Id
		projectedBlocks[index] = schemas.Block{
			Id:            flattenedBlock.Id,
			BlockPackId:   blockPackId,
			ParentBlockId: flattenedBlock.ParentBlockId,
			PrevBlockId:   flattenedBlock.PrevBlockId,
			NextBlockId:   flattenedBlock.NextBlockId,
			Type:          flattenedBlock.Type,
			Props:         flattenedBlock.Props,
			Content:       flattenedBlock.Content,
		}
	}

	tx := s.db.WithContext(ctx).Begin()

	lockingStrength := "UPDATE"
	var blockPack schemas.BlockPack
	if err := tx.Model(&schemas.BlockPack{}).
		Select("id").
		Scopes(scopes.Locking(&lockingStrength)).
		Where("id = ? AND deleted_at IS NULL", blockPackId).
		First(&blockPack).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock block pack for projection: %w", err)
	}

	var document schemas.BlockPackYjsDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Scopes(scopes.Locking(&lockingStrength)).
		Where("block_pack_id = ? AND deleted_at IS NULL", blockPackId).
		First(&document).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock yjs document for projection: %w", err)
	}

	if input.ProjectedSequence <= document.ProjectedUntilSequence {
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("failed to commit stale block projection: %w", err)
		}

		return &dtos.ApplyBlockProjectionResult{
			Applied:                false,
			ProjectedUntilSequence: document.ProjectedUntilSequence,
		}, nil
	}

	if input.ProjectedSequence > document.LastUpdateSequence {
		tx.Rollback()

		return nil, fmt.Errorf("block projection target update sequence exceeds durable yjs state")
	}

	type existingBlock struct {
		Id          uuid.UUID `gorm:"column:id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}

	existingBlocks := []existingBlock{}
	if len(blockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Select("id, block_pack_id").
			Scopes(scopes.Locking(&lockingStrength)).
			Where("id IN ?", blockIds).
			Find(&existingBlocks).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to lock projected blocks: %w", err)
		}
	}

	for _, existingBlock := range existingBlocks {
		if existingBlock.BlockPackId != blockPackId {
			tx.Rollback()

			return nil, fmt.Errorf("block projection contains an id owned by another block pack")
		}
	}

	now := time.Now()

	if len(projectedBlocks) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"block_pack_id":   blockPackId,
				"parent_block_id": gorm.Expr("EXCLUDED.parent_block_id"),
				"prev_block_id":   gorm.Expr("EXCLUDED.prev_block_id"),
				"next_block_id":   gorm.Expr("EXCLUDED.next_block_id"),
				"type":            gorm.Expr("EXCLUDED.type"),
				"props":           gorm.Expr("EXCLUDED.props"),
				"content":         gorm.Expr("EXCLUDED.content"),
				"updated_at":      now,
			}),
		}).CreateInBatches(&projectedBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to bulk upsert block projection: %w", err)
		}
	}

	deleteQuery := tx.Where("block_pack_id = ?", blockPackId)
	if len(blockIds) > 0 {
		deleteQuery = deleteQuery.Where("id NOT IN ?", blockIds)
	}
	if err := deleteQuery.Delete(&schemas.Block{}).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to delete removed projected blocks: %w", err)
	}

	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Where("id = ?", document.Id).
		Updates(map[string]any{
			"projected_until_sequence": input.ProjectedSequence,
			"updated_at":               now,
		}).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to update block projection checkpoint: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit block projection: %w", err)
	}

	return &dtos.ApplyBlockProjectionResult{
		Applied:                true,
		ProjectedUntilSequence: input.ProjectedSequence,
	}, nil
}

func (s *BlockService) ApplyMany(
	ctx context.Context,
	inputs []dtos.ApplyBlockProjectionDocumentInput,
) (dtos.ApplyBlockProjectionDocumentResult, error) {
	if len(inputs) == 0 {
		return dtos.ApplyBlockProjectionDocumentResult{}, nil
	}

	type preparedProjection struct {
		blockPackId  uuid.UUID
		projectedSeq int64
		blocks       []schemas.Block
		blockIds     []uuid.UUID
	}

	preparedProjections := make([]preparedProjection, 0, len(inputs))
	blockPackIdSet := make(map[uuid.UUID]bool, len(inputs))
	blockIdSet := make(map[uuid.UUID]bool)
	for _, input := range inputs {
		if input.BlockPackId == uuid.Nil {
			return nil, fmt.Errorf("block projection requires a block pack id")
		}
		if input.Projection.SchemaId != constants.YjsBlockPackSchemaId ||
			input.Projection.SchemaVersion != constants.YjsBlockPackSchemaVersion {
			return nil, fmt.Errorf("block projection source schema is not supported")
		}
		if input.Projection.ProjectedSequence < 0 {
			return nil, fmt.Errorf("block projection target update sequence must not be negative")
		}
		if blockPackIdSet[input.BlockPackId] {
			return nil, fmt.Errorf("duplicate block projection block pack id")
		}
		blockPackIdSet[input.BlockPackId] = true

		flattenedBlocks, _, exception := s.editableBlockAdapter.FlattenManyToRaw(input.Projection.Blocks)
		if exception != nil {
			return nil, fmt.Errorf("failed to flatten block projection: %w", exception)
		}

		blocks := make([]schemas.Block, len(flattenedBlocks))
		blockIds := make([]uuid.UUID, len(flattenedBlocks))
		for index, flattenedBlock := range flattenedBlocks {
			if blockIdSet[flattenedBlock.Id] {
				return nil, fmt.Errorf("duplicate block id in projection batch")
			}
			blockIdSet[flattenedBlock.Id] = true
			blockIds[index] = flattenedBlock.Id
			blocks[index] = schemas.Block{
				Id:            flattenedBlock.Id,
				BlockPackId:   input.BlockPackId,
				ParentBlockId: flattenedBlock.ParentBlockId,
				PrevBlockId:   flattenedBlock.PrevBlockId,
				NextBlockId:   flattenedBlock.NextBlockId,
				Type:          flattenedBlock.Type,
				Props:         flattenedBlock.Props,
				Content:       flattenedBlock.Content,
			}
		}

		preparedProjections = append(preparedProjections, preparedProjection{
			blockPackId:  input.BlockPackId,
			projectedSeq: input.Projection.ProjectedSequence,
			blocks:       blocks,
			blockIds:     blockIds,
		})
	}

	blockPackIds := make([]uuid.UUID, 0, len(blockPackIdSet))
	for blockPackId := range blockPackIdSet {
		blockPackIds = append(blockPackIds, blockPackId)
	}
	sort.Slice(blockPackIds, func(left, right int) bool {
		return string(blockPackIds[left][:]) < string(blockPackIds[right][:])
	})

	tx := s.db.WithContext(ctx).Begin()

	lockingStrength := options.LockingStrengthUpdate
	var blockPacks []schemas.BlockPack
	if err := tx.Model(&schemas.BlockPack{}).
		Select("id").
		Where("id IN ? AND deleted_at IS NULL", blockPackIds).
		Scopes(scopes.Locking(&lockingStrength)).
		Order("id ASC").
		Find(&blockPacks).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock block packs for projection: %w", err)
	}

	activeBlockPackIdSet := make(map[uuid.UUID]bool, len(blockPacks))
	for _, blockPack := range blockPacks {
		activeBlockPackIdSet[blockPack.Id] = true
	}

	var documents []schemas.BlockPackYjsDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Where("block_pack_id IN ? AND deleted_at IS NULL", blockPackIds).
		Scopes(scopes.Locking(&lockingStrength)).
		Order("block_pack_id ASC").
		Find(&documents).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock yjs documents for projection: %w", err)
	}

	documentByBlockPackId := make(map[uuid.UUID]schemas.BlockPackYjsDocument, len(documents))
	for _, document := range documents {
		documentByBlockPackId[document.BlockPackId] = document
	}

	applicableProjections := make([]preparedProjection, 0, len(preparedProjections))
	for _, projection := range preparedProjections {
		if !activeBlockPackIdSet[projection.blockPackId] {
			continue
		}
		document, exists := documentByBlockPackId[projection.blockPackId]
		if !exists {
			continue
		}
		if projection.projectedSeq > document.LastUpdateSequence {
			tx.Rollback()

			return nil, fmt.Errorf("block projection target update sequence exceeds durable yjs state")
		}
		if projection.projectedSeq <= document.ProjectedUntilSequence {
			continue
		}
		applicableProjections = append(applicableProjections, projection)
	}
	if len(applicableProjections) == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("failed to commit stale block projections: %w", err)
		}

		return dtos.ApplyBlockProjectionDocumentResult{}, nil
	}

	projectedBlockIds := make([]uuid.UUID, 0)
	allProjectedBlocks := make([]schemas.Block, 0)
	applicableBlockPackIds := make([]uuid.UUID, 0, len(applicableProjections))
	for _, projection := range applicableProjections {
		applicableBlockPackIds = append(applicableBlockPackIds, projection.blockPackId)
		projectedBlockIds = append(projectedBlockIds, projection.blockIds...)
		allProjectedBlocks = append(allProjectedBlocks, projection.blocks...)
	}

	type existingBlock struct {
		Id          uuid.UUID `gorm:"column:id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}
	var existingBlocks []existingBlock
	if len(projectedBlockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Select("id, block_pack_id").
			Where("id IN ?", projectedBlockIds).
			Scopes(scopes.Locking(&lockingStrength)).
			Find(&existingBlocks).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to lock projected blocks: %w", err)
		}
	}

	projectedBlockPackIdById := make(map[uuid.UUID]uuid.UUID, len(allProjectedBlocks))
	for _, block := range allProjectedBlocks {
		projectedBlockPackIdById[block.Id] = block.BlockPackId
	}
	for _, existingBlock := range existingBlocks {
		if projectedBlockPackIdById[existingBlock.Id] != existingBlock.BlockPackId {
			tx.Rollback()

			return nil, fmt.Errorf("block projection contains an id owned by another block pack")
		}
	}

	deleteQuery := tx.Where("block_pack_id IN ?", applicableBlockPackIds)
	if len(projectedBlockIds) > 0 {
		deleteQuery = deleteQuery.Where("id NOT IN ?", projectedBlockIds)
	}
	result := deleteQuery.Delete(&schemas.Block{})
	if err := result.Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to delete removed projected blocks: %w", err)
	}

	now := time.Now()
	if len(allProjectedBlocks) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"block_pack_id":   gorm.Expr("EXCLUDED.block_pack_id"),
				"parent_block_id": gorm.Expr("EXCLUDED.parent_block_id"),
				"prev_block_id":   gorm.Expr("EXCLUDED.prev_block_id"),
				"next_block_id":   gorm.Expr("EXCLUDED.next_block_id"),
				"type":            gorm.Expr("EXCLUDED.type"),
				"props":           gorm.Expr("EXCLUDED.props"),
				"content":         gorm.Expr("EXCLUDED.content"),
				"updated_at":      now,
			}),
		}).CreateInBatches(&allProjectedBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to bulk upsert block projections: %w", err)
		}
	}

	valueRows := make([]string, 0, len(applicableProjections))
	args := make([]any, 0, len(applicableProjections)*2)
	for _, projection := range applicableProjections {
		valueRows = append(valueRows, "(?::uuid, ?::bigint)")
		args = append(args, projection.blockPackId, projection.projectedSeq)
	}
	args = append(args, now)
	updateQuery := `
		WITH target(block_pack_id, projected_until_sequence) AS (
			VALUES ` + strings.Join(valueRows, ",") + `
		)
		UPDATE "BlockPackYjsDocumentTable" AS document
		SET projected_until_sequence = target.projected_until_sequence, updated_at = ?
		FROM target
		WHERE document.block_pack_id = target.block_pack_id
			AND document.deleted_at IS NULL
			AND document.projected_until_sequence < target.projected_until_sequence
			AND document.last_update_sequence >= target.projected_until_sequence
		RETURNING document.block_pack_id`

	var appliedDocuments []struct {
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}
	result = tx.Raw(updateQuery, args...).Scan(&appliedDocuments)
	if err := result.Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to update block projection checkpoints: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit block projections: %w", err)
	}

	appliedBlockPackIds := make([]uuid.UUID, len(appliedDocuments))
	for index, document := range appliedDocuments {
		appliedBlockPackIds[index] = document.BlockPackId
	}

	return dtos.ApplyBlockProjectionDocumentResult(appliedBlockPackIds), nil
}

func (s *BlockService) SearchPrivateBlocks(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockInput,
) (*gqlmodels.SearchBlockConnection, *exceptions.Exception) {
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.Block{}).
		Select(`"BlockTable".*`).
		Joins(`INNER JOIN "BlockPackTable" ON "BlockPackTable".id = "BlockTable".block_pack_id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = "SubShelfTable".root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		}).
		Scopes(s.blockPackScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(s.subShelfScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(s.blockScope.IncludePreloads([]schemas.BlockRelation{schemas.BlockRelation_Children}))

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		pattern := "%" + gqlInput.Query + "%"
		query = query.Where(`"BlockTable".content::text ILIKE ? OR "BlockTable".props::text ILIKE ? OR "BlockTable".type::text ILIKE ?`, pattern, pattern, pattern)
	}

	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchBlockCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"BlockTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchBlockSortByType:
			query = query.Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".created_at ` + cending)
		case gqlmodels.SearchBlockSortByLastUpdate:
			query = query.Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".type ` + cending).Order(`"BlockTable".created_at ` + cending)
		case gqlmodels.SearchBlockSortByCreatedAt:
			query = query.Order(`"BlockTable".created_at ` + cending).Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending)
		default:
			query = query.Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var blocks []schemas.Block
	if err := query.Find(&blocks).Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	hasNextPage := len(blocks) > limit
	searchEdges := make([]*gqlmodels.SearchBlockEdge, len(blocks))
	for index, block := range blocks {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchBlockCursorFields]{Fields: gqlmodels.SearchBlockCursorFields{ID: block.Id}}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchBlockEdge{EncodedSearchCursor: *encodedSearchCursor, Node: block.ToPrivateBlock()}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6

	return &gqlmodels.SearchBlockConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
