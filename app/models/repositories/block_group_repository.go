package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ============================== Definitions ============================== */

type BlockGroupRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HasPermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionAndGetValidIds(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	CreateOneByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	CreateManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, inputs []inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockGroupInput, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockGroupRepository struct{}

func NewBlockGroupRepository() BlockGroupRepositoryInterface {
	return &BlockGroupRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockGroupRepository) HasPermission(
	id uuid.UUID,
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
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockGroupTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockGroupRepository) HasPermissions(
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
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockGroupTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockGroupRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockGroupTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockGroup schemas.BlockGroup
	result := query.First(&blockGroup)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &blockGroup, nil
}

func (r *BlockGroupRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.BlockGroup, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockGroupTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockGroups []schemas.BlockGroup
	result := query.Find(&blockGroups)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}
	if len(blockGroups) == 0 {
		return nil, exceptions.BlockGroup.NotFound()
	}

	return blockGroups, nil
}

// Similar to the `HasPermissions`, but with best effort strategy,
// if some of the ids is not valid or exist, they'll be not returned at the end.
//
// Note that the `HasPermission` doesn't need this best effort strategy.
func (r *BlockGroupRepository) CheckPermissionAndGetValidIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("\"BlockGroupTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockGroupTable\".deleted_at IS NULL")
	}

	var validIds []uuid.UUID
	if err := query.Scan(&validIds).Error; err != nil {
		return make([]uuid.UUID, len(ids)), exceptions.BlockGroup.NotFound().WithError(err)
	}

	return validIds, nil
}

func (r *BlockGroupRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
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

func (r *BlockGroupRepository) CreateOneByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockGroupInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository()

	ownerId, blockPack, exception := blockPackRepository.CheckPermissionAndGetOneWithOwnerIdById(
		blockPackId,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}
	if ownerId == nil || blockPack == nil {
		return nil, exceptions.BlockPack.NoPermission("get owner's block pack")
	}

	var newBlockGroup schemas.BlockGroup
	if err := copier.Copy(&newBlockGroup, &input); err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroup)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}

	return &newBlockGroup.Id, nil
}

func (r *BlockGroupRepository) CreateManyByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input []inputs.CreateBlockGroupInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository()

	ownerId, blockPack, exception := blockPackRepository.CheckPermissionAndGetOneWithOwnerIdById(
		blockPackId,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}
	if ownerId == nil || blockPack == nil {
		return nil, exceptions.BlockPack.NoPermission("get owner's block pack")
	}

	var newBlockGroups []schemas.BlockGroup
	if err := copier.Copy(&newBlockGroups, &input); err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}
	ids := make([]uuid.UUID, len(input))
	for index := range newBlockGroups { // use index to modify the elements of slice, since the value from the `range` is only a copy
		newBlockGroups[index].Id = uuid.New()    // generate the id here, so that we can return the ids in the same order of the input
		newBlockGroups[index].OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
		newBlockGroups[index].BlockPackId = blockPackId
		ids[index] = newBlockGroups[index].Id
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroups)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}

	return ids, nil
}

func (r *BlockGroupRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockGroupInput,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingBlockGroup, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlockGroup)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlockGroup,
		).WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	return &updates, nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedOneById(
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
			return exceptions.BlockGroup.NoPermission("restore a deleted block group")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
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

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockGroup.NoPermission("restore deleted block groups")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
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
			return exceptions.BlockGroup.NoPermission("soft delete a block group")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if !parsedOptions.SkipPermissionCheck {
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
			return exceptions.BlockGroup.NoPermission("soft delete block groups")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) HardDeleteOneById(
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
			return exceptions.BlockGroup.NoPermission("hard delete a block group")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.BlockGroup{})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
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

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockGroup.NoPermission("hard delete block groups")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.BlockGroup{})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}
