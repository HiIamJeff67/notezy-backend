package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Block, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	CreateOneByBlockGroupId(blockGroupId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	CreateManyByBlockGroupId(blockGroupId uuid.UUID, userId uuid.UUID, input []inputs.CreateBlockInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	CreateManyByBlockGroupIds(userId uuid.UUID, input []inputs.CreateBlockGroupContentInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockInput, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	UpdateManyByIds(userId uuid.UUID, input []inputs.UpdateBlockByIdInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Block, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Block, *exceptions.Exception)
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception

	/* ============================== System Only Method ============================== */

	BulkCheckPermissionsAndGetManyByIds(inputs []inputs.BulkCheckBlockPermissionInput, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]bool, []schemas.Block, *exceptions.Exception)
	BulkCreateMany(inputs []inputs.BulkCreateBlockGroupContentInput, opts ...options.RepositoryOptions) ([]bool, *exceptions.Exception)
	BulkUpdateMany(inputs []inputs.BulkUpdateBlockInput, opts ...options.RepositoryOptions) ([]bool, *exceptions.Exception)
	BulkDeleteMany(inputs []inputs.BulkDeleteBlockInput, opts ...options.RepositoryOptions) ([]bool, *exceptions.Exception)
}

type BlockRepository struct {
	blockScope scopes.BlockScopeInterface
}

func NewBlockRepository(blockScope scopes.BlockScopeInterface) BlockRepositoryInterface {
	return &BlockRepository{
		blockScope: blockScope,
	}
}

func (r *BlockRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	if parsedOptions.DB == nil {
		parsedOptions.DB = models.NotezyDB
	}

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.Block{}).
		Select("1").
		Scopes(r.blockScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *BlockRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.Block{}).
		Select(`DISTINCT "BlockTable".id`).
		Scopes(r.blockScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *BlockRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var block schemas.Block
	result := parsedOptions.DB.
		Model(&schemas.Block{}).
		Scopes(r.blockScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&block)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.NotFound().WithOrigin(result.Error)},
		{First: block.Id == uuid.Nil, Second: exceptions.Block.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &block, nil
}

func (r *BlockRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.Block, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var blocks []schemas.Block
	result := parsedOptions.DB.
		Model(&schemas.Block{}).
		Scopes(r.blockScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&blocks)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.NotFound().WithOrigin(result.Error)},
		{First: len(blocks) == 0, Second: exceptions.Block.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return blocks, nil
}

func (r *BlockRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
}

func (r *BlockRepository) CreateOneByBlockGroupId(
	blockGroupId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository(scopes.NewBlockGroupScope())

		if !blockGroupRepository.HasPermission(
			blockGroupId,
			userId,
			allowedPermissions,
			opts...,
		) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	var newBlock schemas.Block
	if err := copier.Copy(&newBlock, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Block.InvalidInput().WithOrigin(err)
	}
	newBlock.BlockGroupId = blockGroupId

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"block_group_id", "parent_block_id", "type", "props", "content", "updated_at",
				}),
			},
		).
		Create(&newBlock)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlock.Id == uuid.Nil, Second: exceptions.Block.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newBlock.Id, nil
}

func (r *BlockRepository) CreateManyByBlockGroupId(
	blockGroupId uuid.UUID,
	userId uuid.UUID,
	input []inputs.CreateBlockInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository(scopes.NewBlockGroupScope())

		if !blockGroupRepository.HasPermission(
			blockGroupId,
			userId,
			allowedPermissions,
			opts...,
		) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	newBlocks := make([]schemas.Block, len(input))
	for index, in := range input {
		var newBlock schemas.Block
		if err := copier.Copy(&newBlock, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.InvalidInput().WithOrigin(err)
		}
		newBlock.BlockGroupId = blockGroupId
		newBlocks[index] = newBlock
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"block_group_id", "parent_block_id", "type", "props", "content", "updated_at",
				}),
			},
		).
		CreateInBatches(&newBlocks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newBlockIds := make([]uuid.UUID, len(newBlocks))
	for index, newBlock := range newBlocks {
		newBlockIds[index] = newBlock.Id
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newBlockIds, nil
}

func (r *BlockRepository) CreateManyByBlockGroupIds(
	userId uuid.UUID,
	input []inputs.CreateBlockGroupContentInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if !parsedOptions.SkipPermissionCheck {
		shouldStartTransaction := !parsedOptions.IsTransactionStarted
		if shouldStartTransaction {
			parsedOptions.DB = parsedOptions.DB.Begin()
			opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
			opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
		}

		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository(scopes.NewBlockGroupScope())

		blockGroupIds := make([]uuid.UUID, len(input))
		for index, in := range input {
			blockGroupIds[index] = in.BlockGroupId
		}

		validIds, exception := blockGroupRepository.CheckPermissionAndGetValidIds(
			blockGroupIds,
			userId,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return nil, exception
		}

		validIdMap := make(map[uuid.UUID]bool)
		for _, validId := range validIds {
			validIdMap[validId] = true
		}

		var newBlocks []schemas.Block
		for _, in := range input {
			if validIdMap[in.BlockGroupId] {
				for _, inputBlock := range in.Blocks {
					var newBlock schemas.Block
					if err := copier.Copy(&newBlock, &inputBlock); err != nil {
						parsedOptions.DB.Rollback()
						return nil, exceptions.Block.InvalidInput().WithOrigin(err)
					}
					newBlock.BlockGroupId = in.BlockGroupId
					newBlocks = append(newBlocks, newBlock)
				}
			}
		}

		result := parsedOptions.DB.Model(&schemas.Block{}).
			Clauses(
				clause.Returning{Columns: []clause.Column{{Name: "id"}}},
				clause.OnConflict{
					Columns: []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"block_group_id", "parent_block_id", "type", "props", "content", "updated_at",
					}),
				},
			).
			CreateInBatches(&newBlocks, parsedOptions.BatchSize)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return nil, exception
		}

		newBlockIds := make([]uuid.UUID, len(newBlocks))
		for index, newBlock := range newBlocks {
			newBlockIds[index] = newBlock.Id
		}

		if shouldStartTransaction {
			if err := parsedOptions.DB.Commit().Error; err != nil {
				parsedOptions.DB.Rollback()
				return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
			}
		}

		return newBlockIds, nil
	}

	var newBlocks []schemas.Block
	for _, in := range input {
		for _, inputBlock := range in.Blocks {
			var newBlock schemas.Block
			if err := copier.Copy(&newBlock, &inputBlock); err != nil {
				parsedOptions.DB.Rollback()
				return nil, exceptions.Block.InvalidInput().WithOrigin(err)
			}
			newBlock.BlockGroupId = in.BlockGroupId
			newBlocks = append(newBlocks, newBlock)
		}
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"block_group_id", "parent_block_id", "type", "props", "content", "updated_at",
				}),
			},
		).
		CreateInBatches(&newBlocks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newBlockIds := make([]uuid.UUID, len(newBlocks))
	for index, newBlock := range newBlocks {
		newBlockIds[index] = newBlock.Id
	}

	return newBlockIds, nil
}

func (r *BlockRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockInput,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	// maybe we need a more efficient way to update the field of blocks
	// since they will be used quite frequently

	existingBlock, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlock)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlock,
		).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *BlockRepository) UpdateManyByIds(
	userId uuid.UUID,
	input []inputs.UpdateBlockByIdInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	// since there're no nullable fields may update in this repository function,
	// so we don't have to use partial update process here actually,
	// we can simply use COALESCE to maintain the original value in each fields for each null fields in the passing input
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	isBlockValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		ids := make([]uuid.UUID, len(input))
		for index, in := range input {
			ids[index] = in.Id
		}

		validBlocks, exception := r.CheckPermissionsAndGetManyByIds(ids, userId, nil, allowedPermissions, opts...)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Block.NoPermission("update these blocks")
		}

		for _, validBlock := range validBlocks {
			isBlockValid[validBlock.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	var ids []uuid.UUID
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isBlockValid[in.Id] {
			continue
		}

		setParentBlockIdNull := util.CheckSetNull(in.PartialUpdateInput.SetNull, "ParentBlockId")
		valuePlaceholders = append(valuePlaceholders, `(?::uuid, ?::"BlockType", ?::jsonb, ?::jsonb, ?::uuid, ?::uuid, ?::boolean)`)
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Type,
			in.PartialUpdateInput.Values.Props,
			in.PartialUpdateInput.Values.Content,
			in.PartialUpdateInput.Values.BlockGroupId,
			in.PartialUpdateInput.Values.ParentBlockId,
			setParentBlockIdNull,
		)
		ids = append(ids, in.Id)
	}

	sql := fmt.Sprintf(`
		UPDATE "BlockTable" AS b
		SET
			type = COALESCE(v.type::"BlockType", b.type),
			props = COALESCE(v.props::jsonb, b.props),
			content = COALESCE(v.content::jsonb, b.content),
			block_group_id = COALESCE(v.block_group_id::uuid, b.block_group_id),
			parent_block_id = CASE 
				WHEN v.set_parent_block_id_null::boolean THEN NULL 
				ELSE COALESCE(v.parent_block_id::uuid, b.parent_block_id)
			END,
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, type, props, content, block_group_id, parent_block_id, set_parent_block_id_null)
		WHERE b.id = v.id::uuid AND b.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *BlockRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var restoredBlock schemas.Block
	query := parsedOptions.DB.Model(&restoredBlock).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionCheck(id, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"BlockTable".id = ?`, id).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredBlock, nil
}

func (r *BlockRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Block, *exceptions.Exception) {
	if len(ids) == 0 {
		return []schemas.Block{}, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var restoredBlocks []schemas.Block
	query := parsedOptions.DB.Model(&restoredBlocks).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionChecks(ids, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"BlockTable".id IN ?`, ids).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredBlocks, nil
}

func (r *BlockRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var deletedBlock schemas.Block
	query := parsedOptions.DB.Model(&deletedBlock).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionCheck(id, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"BlockTable".id = ?`, id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &deletedBlock, nil
}

func (r *BlockRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Block, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var deletedBlocks []schemas.Block
	query := parsedOptions.DB.Model(&deletedBlocks).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionChecks(ids, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"BlockTable".id IN ?`, ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return deletedBlocks, nil
}

func (r *BlockRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	query := parsedOptions.DB.Model(&schemas.Block{}).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionCheck(id, userId, allowedPermissions))
	}

	result := query.
		Where(`"BlockTable".id = ?`, id).
		Delete(&schemas.Block{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *BlockRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	query := parsedOptions.DB.Model(&schemas.Block{}).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		query = query.Scopes(r.blockScope.PassPermissionChecks(ids, userId, allowedPermissions))
	}

	result := query.
		Where(`"BlockTable".id IN ?`, ids).
		Delete(&schemas.Block{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

/* ============================== System Only Method ============================== */

func (r *BlockRepository) BulkCheckPermissionsAndGetManyByIds(
	inputs []inputs.BulkCheckBlockPermissionInput,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]bool, []schemas.Block, *exceptions.Exception) {
	if len(inputs) == 0 {
		return []bool{}, []schemas.Block{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	successes := make([]bool, len(inputs))
	ids := make([]uuid.UUID, 0, len(inputs))
	userIds := make([]uuid.UUID, 0, len(inputs))
	for _, in := range inputs {
		ids = append(ids, in.Id)
		userIds = append(userIds, in.UserId)
	}

	var validTargets []struct {
		Id     uuid.UUID `gorm:"column:id"`
		UserId uuid.UUID `gorm:"column:user_id"`
	}
	result := parsedOptions.DB.Model(&schemas.Block{}).
		Select(`"BlockTable".id, uts.user_id`).
		Joins(`INNER JOIN "BlockGroupTable" AS bg ON bg.id = "BlockTable".block_group_id`).
		Joins(`INNER JOIN "BlockPackTable" AS bp ON bp.id = bg.block_pack_id`).
		Joins(`INNER JOIN "SubShelfTable" AS ss ON ss.id = bp.parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" AS uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where(`"BlockTable".id IN ?`, ids).
		Where("uts.user_id IN ? AND uts.permission IN ?", userIds, allowedPermissions).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scan(&validTargets)
	if result.Error != nil {
		return nil, nil, exceptions.Block.NotFound().WithOrigin(result.Error)
	}

	validTargetByUserId := make(map[[2]uuid.UUID]bool, len(validTargets))
	for _, validTarget := range validTargets {
		validTargetByUserId[[2]uuid.UUID{validTarget.Id, validTarget.UserId}] = true
	}

	validIdSet := make(map[uuid.UUID]bool, len(validTargets))
	for _, in := range inputs {
		if validTargetByUserId[[2]uuid.UUID{in.Id, in.UserId}] {
			validIdSet[in.Id] = true
		}
	}

	validIds := make([]uuid.UUID, 0, len(validIdSet))
	for validId := range validIdSet {
		validIds = append(validIds, validId)
	}
	if len(validIds) == 0 {
		return successes, []schemas.Block{}, nil
	}

	var blocks []schemas.Block
	result = parsedOptions.DB.Model(&schemas.Block{}).
		Where(`"BlockTable".id IN ?`, validIds).
		Scopes(r.blockScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&blocks)
	if result.Error != nil {
		return nil, nil, exceptions.Block.NotFound().WithOrigin(result.Error)
	}

	foundIdSet := make(map[uuid.UUID]bool, len(blocks))
	for _, block := range blocks {
		foundIdSet[block.Id] = true
	}
	for index, in := range inputs {
		if validTargetByUserId[[2]uuid.UUID{in.Id, in.UserId}] && foundIdSet[in.Id] {
			successes[index] = true
		}
	}

	return successes, blocks, nil
}

func (r *BlockRepository) BulkCreateMany(
	bulkInputs []inputs.BulkCreateBlockGroupContentInput,
	opts ...options.RepositoryOptions,
) ([]bool, *exceptions.Exception) {
	if len(bulkInputs) == 0 {
		return []bool{}, exceptions.Block.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	checkInputs := make([]inputs.BulkCheckBlockGroupPermissionInput, len(bulkInputs))
	for index, in := range bulkInputs {
		checkInputs[index] = inputs.BulkCheckBlockGroupPermissionInput{
			UserId: in.UserId,
			Id:     in.BlockGroupId,
		}
	}
	blockGroupRepository := NewBlockGroupRepository(scopes.NewBlockGroupScope())
	checkOptions := append(opts, options.WithTransactionDB(parsedOptions.DB))
	checkOptions = append(checkOptions, options.WithOnlyDeleted(types.Ternary_Negative))
	checkOptions = append(checkOptions, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	successes, _, exception := blockGroupRepository.BulkCheckPermissionsAndGetManyByIds(checkInputs, nil, allowedPermissions, checkOptions...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newBlocks := make([]schemas.Block, 0)
	for index, in := range bulkInputs {
		if !successes[index] {
			continue
		}
		for _, inputBlock := range in.Blocks {
			newBlocks = append(newBlocks, schemas.Block{
				Id:            inputBlock.Id,
				ParentBlockId: inputBlock.ParentBlockId,
				BlockGroupId:  in.BlockGroupId,
				Type:          inputBlock.Type,
				Props:         inputBlock.Props,
				Content:       inputBlock.Content,
			})
		}
	}
	if len(newBlocks) == 0 {
		if shouldStartTransaction {
			parsedOptions.DB.Rollback()
		}
		return successes, nil
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
			clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"block_group_id", "parent_block_id", "type", "props", "content", "updated_at",
				}),
			},
		).
		CreateInBatches(&newBlocks, parsedOptions.BatchSize)
	if result.Error != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Block.FailedToCreate().WithOrigin(result.Error)
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return successes, nil
}

func (r *BlockRepository) BulkUpdateMany(
	bulkInputs []inputs.BulkUpdateBlockInput,
	opts ...options.RepositoryOptions,
) ([]bool, *exceptions.Exception) {
	if len(bulkInputs) == 0 {
		return []bool{}, exceptions.Block.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	checkInputs := make([]inputs.BulkCheckBlockPermissionInput, len(bulkInputs))
	for index, in := range bulkInputs {
		checkInputs[index] = inputs.BulkCheckBlockPermissionInput{
			UserId: in.UserId,
			Id:     in.Id,
		}
	}
	checkOptions := append(opts, options.WithTransactionDB(parsedOptions.DB))
	checkOptions = append(checkOptions, options.WithOnlyDeleted(types.Ternary_Negative))
	checkOptions = append(checkOptions, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	successes, _, exception := r.BulkCheckPermissionsAndGetManyByIds(checkInputs, nil, allowedPermissions, checkOptions...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	valuePlaceholders := make([]string, 0, len(bulkInputs))
	valueArgs := make([]interface{}, 0, len(bulkInputs)*8)
	for index, in := range bulkInputs {
		if !successes[index] {
			continue
		}

		setParentBlockIdNull := util.CheckSetNull(in.PartialUpdateInput.SetNull, "ParentBlockId")
		valuePlaceholders = append(valuePlaceholders, `(?::int, ?::uuid, ?::"BlockType", ?::jsonb, ?::jsonb, ?::uuid, ?::uuid, ?::boolean)`)
		valueArgs = append(valueArgs,
			index,
			in.Id,
			in.PartialUpdateInput.Values.Type,
			in.PartialUpdateInput.Values.Props,
			in.PartialUpdateInput.Values.Content,
			in.PartialUpdateInput.Values.BlockGroupId,
			in.PartialUpdateInput.Values.ParentBlockId,
			setParentBlockIdNull,
		)
	}
	if len(valuePlaceholders) == 0 {
		if shouldStartTransaction {
			parsedOptions.DB.Rollback()
		}
		return successes, nil
	}

	sql := fmt.Sprintf(`
		WITH payload(idx, id, type, props, content, block_group_id, parent_block_id, set_parent_block_id_null) AS (
			VALUES %s
		),
		updated AS (
			UPDATE "BlockTable" AS b
			SET
				type = COALESCE(v.type::"BlockType", b.type),
				props = COALESCE(v.props::jsonb, b.props),
				content = COALESCE(v.content::jsonb, b.content),
				block_group_id = COALESCE(v.block_group_id::uuid, b.block_group_id),
				parent_block_id = CASE
					WHEN v.set_parent_block_id_null::boolean THEN NULL
					ELSE COALESCE(v.parent_block_id::uuid, b.parent_block_id)
				END,
				updated_at = NOW()
			FROM payload AS v
			WHERE b.id = v.id::uuid
				AND b.deleted_at IS NULL
			RETURNING b.id
		)
		SELECT v.idx
		FROM payload AS v
		INNER JOIN updated AS u ON u.id = v.id::uuid
	`, strings.Join(valuePlaceholders, ","))

	var updatedIndexes []struct {
		Index int `gorm:"column:idx"`
	}
	result := parsedOptions.DB.Raw(sql, valueArgs...).Scan(&updatedIndexes)
	if result.Error != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(result.Error)
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	successes = make([]bool, len(bulkInputs))
	for _, updatedIndex := range updatedIndexes {
		if updatedIndex.Index >= 0 && updatedIndex.Index < len(successes) {
			successes[updatedIndex.Index] = true
		}
	}

	return successes, nil
}

func (r *BlockRepository) BulkDeleteMany(
	bulkInputs []inputs.BulkDeleteBlockInput,
	opts ...options.RepositoryOptions,
) ([]bool, *exceptions.Exception) {
	if len(bulkInputs) == 0 {
		return []bool{}, exceptions.Block.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	checkInputs := make([]inputs.BulkCheckBlockPermissionInput, len(bulkInputs))
	for index, in := range bulkInputs {
		checkInputs[index] = inputs.BulkCheckBlockPermissionInput{
			UserId: in.UserId,
			Id:     in.Id,
		}
	}
	checkOptions := append(opts, options.WithTransactionDB(parsedOptions.DB))
	checkOptions = append(checkOptions, options.WithOnlyDeleted(types.Ternary_Negative))
	checkOptions = append(checkOptions, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	successes, _, exception := r.BulkCheckPermissionsAndGetManyByIds(checkInputs, nil, allowedPermissions, checkOptions...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	validIds := make([]uuid.UUID, 0, len(bulkInputs))
	for index, in := range bulkInputs {
		if successes[index] {
			validIds = append(validIds, in.Id)
		}
	}
	if len(validIds) == 0 {
		if shouldStartTransaction {
			parsedOptions.DB.Rollback()
		}
		return successes, nil
	}

	var deletedBlocks []schemas.Block
	result := parsedOptions.DB.Model(&deletedBlocks).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NULL", validIds).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Block.FailedToDelete().WithOrigin(result.Error)
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	deletedIdSet := make(map[uuid.UUID]bool, len(deletedBlocks))
	for _, deletedBlock := range deletedBlocks {
		deletedIdSet[deletedBlock.Id] = true
	}
	for index, in := range bulkInputs {
		if successes[index] && deletedIdSet[in.Id] {
			successes[index] = true
		} else {
			successes[index] = false
		}
	}

	return successes, nil
}
