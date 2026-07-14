package services

import (
	"context"
	"fmt"
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
