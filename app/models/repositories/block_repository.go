package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ============================== Definitions ============================== */

type BlockRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOption) bool
	HasPermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOption) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOption) (*schemas.Block, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOption) ([]schemas.Block, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, opts ...options.RepositoryOption) (*schemas.Block, *exceptions.Exception)
	CreateOneByBlockGroupId(blockGroupId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockInput, opts ...options.RepositoryOption) (*uuid.UUID, *exceptions.Exception)
	CreateManyByBlockGroupId(blockGroupId uuid.UUID, userId uuid.UUID, input []inputs.CreateBlockInput, opts ...options.RepositoryOption) ([]schemas.Block, *exceptions.Exception)
	CreateManyByBlockGroupIds(userId uuid.UUID, input []inputs.CreateBlockGroupContentInput, opts ...options.RepositoryOption) ([]schemas.Block, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockInput, opts ...options.RepositoryOption) (*schemas.Block, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOption) *exceptions.Exception
}

type BlockRepository struct{}

func NewBlockRepository() BlockRepositoryInterface {
	return &BlockRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOption,
) bool {
	options := options.ParseRepositoryOptions(opts...)
	if options.DB == nil {
		options.DB = models.NotezyDB
	}

	subQuery := options.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := options.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch options.OnlyDeleted {
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

func (r *BlockRepository) HasPermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOption,
) bool {
	options := options.ParseRepositoryOptions(opts...)

	subQuery := options.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := options.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch options.OnlyDeleted {
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

func (r *BlockRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOption,
) (*schemas.Block, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

	subQuery := options.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := options.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch options.OnlyDeleted {
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
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound().WithError(err)
	}

	return &block, nil
}

func (r *BlockRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOption,
) ([]schemas.Block, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

	subQuery := options.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := options.DB.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch options.OnlyDeleted {
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
	result := query.First(&blocks)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound().WithError(err)
	}

	return blocks, nil
}

func (r *BlockRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	opts ...options.RepositoryOption,
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
	opts ...options.RepositoryOption,
) (*uuid.UUID, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository()

		if !blockGroupRepository.HasPermission(
			options.DB,
			blockGroupId,
			userId,
			allowedPermissions,
			options.OnlyDeleted,
		) {
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	var newBlock schemas.Block
	if err := copier.Copy(&newBlock, &input); err != nil {
		return nil, exceptions.Block.InvalidInput().WithError(err)
	}
	newBlock.BlockGroupId = blockGroupId

	result := options.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlock)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToCreate().WithError(err)
	}

	return &newBlock.Id, nil
}

func (r *BlockRepository) CreateManyByBlockGroupId(
	blockGroupId uuid.UUID,
	userId uuid.UUID,
	input []inputs.CreateBlockInput,
	opts ...options.RepositoryOption,
) ([]schemas.Block, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		blockGroupRepository := NewBlockGroupRepository()

		if !blockGroupRepository.HasPermission(
			options.DB,
			blockGroupId,
			userId,
			allowedPermissions,
			options.OnlyDeleted,
		) {
			return nil, exceptions.Block.NoPermission("get owner's block group")
		}
	}

	newBlocks := make([]schemas.Block, len(input))
	for index, in := range input {
		var newBlock schemas.Block
		if err := copier.Copy(&newBlock, &in); err != nil {
			return nil, exceptions.Block.InvalidInput().WithError(err)
		}
		newBlock.BlockGroupId = blockGroupId
		newBlocks[index] = newBlock
	}

	result := options.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newBlocks, options.BatchSize)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToCreate().WithError(err)
	}

	return newBlocks, nil
}

func (r *BlockRepository) CreateManyByBlockGroupIds(
	userId uuid.UUID,
	input []inputs.CreateBlockGroupContentInput,
	opts ...options.RepositoryOption,
) ([]schemas.Block, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
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
			options.DB,
			blockGroupIds,
			userId,
			allowedPermissions,
			options.OnlyDeleted,
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
						return nil, exceptions.Block.InvalidInput().WithError(err)
					}
					newBlock.BlockGroupId = in.BlockGroupId
					newBlocks = append(newBlocks, newBlock)
				}
			}
		}

		result := options.DB.Model(&schemas.Block{}).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			CreateInBatches(&newBlocks, options.BatchSize)
		if err := result.Error; err != nil {
			return nil, exceptions.Block.FailedToCreate().WithError(err)
		}

		return newBlocks, nil
	}

	var newBlocks []schemas.Block
	for _, in := range input {
		for _, inputBlock := range in.Blocks {
			var newBlock schemas.Block
			if err := copier.Copy(&newBlock, &inputBlock); err != nil {
				return nil, exceptions.Block.InvalidInput().WithError(err)
			}
			newBlock.BlockGroupId = in.BlockGroupId
			newBlocks = append(newBlocks, newBlock)
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newBlocks, options.BatchSize)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToCreate().WithError(err)
	}

	return newBlocks, nil
}

func (r *BlockRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockInput,
	opts ...options.RepositoryOption,
) (*schemas.Block, *exceptions.Exception) {
	options := options.ParseRepositoryOptions(opts...)

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
		).WithError(err)
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Block.NoChanges()
	}

	return &updates, nil
}

func (r *BlockRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
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
			return exceptions.Block.NoPermission("restore a deleted block")
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.Block.NoPermission("restore deleted blocks")
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
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
			return exceptions.Block.NoPermission("soft delete a block")
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.Block.NoPermission("soft delete blocks")
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
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

	result := options.DB.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.Block{})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOption,
) *exceptions.Exception {
	options := options.ParseRepositoryOptions(opts...)

	if !options.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.Block.NoPermission("hard delete blocks")
		}
	}

	result := options.DB.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.Block{})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}
