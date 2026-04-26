package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	"notezy-backend/app/monitor/logs"
	"notezy-backend/app/monitor/traces"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
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
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateBlocksInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Block, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Block, *exceptions.Exception)
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Block, *exceptions.Exception)
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockRepository struct{}

func NewBlockRepository() BlockRepositoryInterface {
	return &BlockRepository{}
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

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		logs.FInfo(traces.GetTrace(0).FileLineString(), "YEE: %v", err)
		return false
	}

	return count > 0
}

func (r *BlockRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var block schemas.Block
	result := query.First(&block)
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

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blocks []schemas.Block
	result := query.Find(&blocks)
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository()

		if !blockGroupRepository.HasPermission(
			blockGroupId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	var newBlock schemas.Block
	if err := copier.Copy(&newBlock, &input); err != nil {
		return nil, exceptions.Block.InvalidInput().WithOrigin(err)
	}
	newBlock.BlockGroupId = blockGroupId

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlock)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlock.Id == uuid.Nil, Second: exceptions.Block.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository()

		if !blockGroupRepository.HasPermission(
			blockGroupId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	newBlocks := make([]schemas.Block, len(input))
	for index, in := range input {
		var newBlock schemas.Block
		if err := copier.Copy(&newBlock, &in); err != nil {
			return nil, exceptions.Block.InvalidInput().WithOrigin(err)
		}
		newBlock.BlockGroupId = blockGroupId
		newBlocks[index] = newBlock
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newBlocks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	newBlockIds := make([]uuid.UUID, len(newBlocks))
	for index, newBlock := range newBlocks {
		newBlockIds[index] = newBlock.Id
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
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository()

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
						return nil, exceptions.Block.InvalidInput().WithOrigin(err)
					}
					newBlock.BlockGroupId = in.BlockGroupId
					newBlocks = append(newBlocks, newBlock)
				}
			}
		}

		result := parsedOptions.DB.Model(&schemas.Block{}).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			CreateInBatches(&newBlocks, parsedOptions.BatchSize)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
		}); exception != nil {
			return nil, exception
		}

		newBlockIds := make([]uuid.UUID, len(newBlocks))
		for index, newBlock := range newBlocks {
			newBlockIds[index] = newBlock.Id
		}

		return newBlockIds, nil
	}

	var newBlocks []schemas.Block
	for _, in := range input {
		for _, inputBlock := range in.Blocks {
			var newBlock schemas.Block
			if err := copier.Copy(&newBlock, &inputBlock); err != nil {
				return nil, exceptions.Block.InvalidInput().WithOrigin(err)
			}
			newBlock.BlockGroupId = in.BlockGroupId
			newBlocks = append(newBlocks, newBlock)
		}
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newBlocks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
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
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlock)
	if err != nil {
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
		return nil, exception
	}

	return &updates, nil
}

func (r *BlockRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateBlocksInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	// since there're no nullable fields may update in this repository function,
	// so we don't have to use partial update process here actually,
	// we can simply use COALESCE to maintain the original value in each fields for each null fields in the passing input
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

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

		setParentBlockIdNull := false
		if in.PartialUpdateInput.SetNull != nil {
			for field, setNull := range *in.PartialUpdateInput.SetNull {
				if strings.ToLower(field) == "parentblockid" && setNull {
					setParentBlockIdNull = true
					break
				}
			}
		}
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::\"BlockType\", ?::jsonb, ?::jsonb, ?::uuid, ?::uuid, ?::boolean)")
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
		return exception
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}
	blocks, exception := r.CheckPermissionsAndGetManyByIds(
		ids,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil
	}
	logs.Info(traces.GetTrace(0).FileLineString(), "---")
	for _, block := range blocks {
		logs.FInfo(traces.GetTrace(0).FileLineString(), "id: %s, parent_block_id: %v, block_group_id: %s", block.Id.String(), block.ParentBlockId, block.BlockGroupId.String())
	}
	logs.Info(traces.GetTrace(0).FileLineString(), "---")

	return nil
}

func (r *BlockRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Block, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("restore a deleted block")
		}
	}

	var restoredBlock schemas.Block
	result := parsedOptions.DB.Model(&restoredBlock).
		Clauses(clause.Returning{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HavePermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("restore deleted blocks")
		}
	}

	var restoredBlocks []schemas.Block
	result := parsedOptions.DB.Model(restoredBlocks).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("soft delete a block")
		}
	}

	var deletedBlock schemas.Block
	result := parsedOptions.DB.Model(&deletedBlock).
		Clauses(clause.Returning{}).
		Where("id = ? AND deleted_at IS NULL", id).
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HavePermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Block.NoPermission("soft delete blocks")
		}
	}

	var deletedBlocks []schemas.Block
	result := parsedOptions.DB.Model(&deletedBlocks).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockGroup.NoPermission("hard delete a block")
		}
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
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

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HavePermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.Block.NoPermission("hard delete blocks")
		}
	}

	result := parsedOptions.DB.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.Block{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Block.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
