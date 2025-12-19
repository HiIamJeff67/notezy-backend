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

type BlockPackRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HasPermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockPack, *exceptions.Exception)
	CheckPermissionAndGetOneWithOwnerIdById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*uuid.UUID, *schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyWithOwnerIdsByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]uuid.UUID, []schemas.BlockPack, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	CreateOneBySubShelfId(subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockPackInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockPackInput, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockPackRepository struct{}

func NewBlockPackRepository() BlockPackRepositoryInterface {
	return &BlockPackRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockPackRepository) HasPermission(
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
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("\"BlockPackTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockPackRepository) HasPermissions(
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
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id == ss.id").
		Where("\"BlockPackTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockPackRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.BlockPack, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPack schemas.BlockPack
	result := query.First(&blockPack)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &blockPack, nil
}

func (r *BlockPackRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.BlockPack, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPacks []schemas.BlockPack
	result := query.Find(&blockPacks)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}
	if len(blockPacks) == 0 {
		return nil, exceptions.BlockPack.NotFound()
	}

	return blockPacks, nil
}

func (r *BlockPackRepository) CheckPermissionAndGetOneWithOwnerIdById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *schemas.BlockPack, *exceptions.Exception) { // we should also return the owner id for the block groups and blocks
	parsedOptions := options.ParseRepositoryOptions(opts...)

	// note that the subQuery is querying the permission of the current user,
	// 			 and the query is querying the data and the owner id(which may be different from the current user)
	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Select("\"BlockPackTable\".*, owner_uts.user_id AS owner_id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Joins("INNER JOIN \"UsersToShelvesTable\" owner_uts ON ss.root_shelf_id = owner_uts.root_shelf_id AND owner_uts.permission = 'Owner'").
		Where("\"BlockPackTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPackWithOwnerId struct {
		schemas.BlockPack
		OwnerId uuid.UUID `gorm:"column:owner_id;"`
	}
	result := query.First(&blockPackWithOwnerId)
	if err := result.Error; err != nil {
		return nil, nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &blockPackWithOwnerId.OwnerId, &blockPackWithOwnerId.BlockPack, nil
}

func (r *BlockPackRepository) CheckPermissionsAndGetManyWithOwnerIdsByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, []schemas.BlockPack, *exceptions.Exception) { // we should also return the owner id for the block groups and blocks
	parsedOptions := options.ParseRepositoryOptions(opts...)

	// note that the subQuery is querying the permission of the current user,
	// 			 and the query is querying the data and the owner id(which may be different from the current user)
	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Select("\"BlockPackTable\".*, owner_uts.user_id AS owner_id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Joins("INNER JOIN \"UsersToShelvesTable\" owner_uts ON ss.root_shelf_id = owner_uts.root_shelf_id AND owner_uts.permission = 'Owner'").
		Where("\"BlockPackTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPacksWithOwnerIds []struct {
		schemas.BlockPack
		ownerId uuid.UUID `gorm:"column:owner_id;"`
	}
	result := query.Find(&blockPacksWithOwnerIds)
	if err := result.Error; err != nil {
		return nil, nil, exceptions.BlockPack.NotFound().WithError(err)
	}
	if len(blockPacksWithOwnerIds) == 0 {
		return nil, nil, exceptions.BlockPack.NotFound()
	}

	ownerIds := make([]uuid.UUID, len(blockPacksWithOwnerIds))
	blockPacks := make([]schemas.BlockPack, len(blockPacksWithOwnerIds))
	for index, element := range blockPacksWithOwnerIds {
		ownerIds[index] = element.ownerId
		blockPacks[index] = element.BlockPack
	}

	return ownerIds, blockPacks, nil
}

func (r *BlockPackRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.BlockPack, *exceptions.Exception) {
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

func (r *BlockPackRepository) CreateOneBySubShelfId(
	subShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockPackInput,
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

		subShelfRepository := NewSubShelfRepository()

		if !subShelfRepository.HasPermission(
			subShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Shelf.NoPermission("create a block pack under this shelf")
		}
	}

	var newBlockPack schemas.BlockPack
	if err := copier.Copy(&newBlockPack, &input); err != nil {
		return nil, exceptions.BlockPack.FailedToCreate().WithError(err)
	}
	newBlockPack.ParentSubShelfId = subShelfId

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockPack)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToCreate().WithError(err)
	}

	return &newBlockPack.Id, nil
}

func (r *BlockPackRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockPackInput,
	opts ...options.RepositoryOptions,
) (*schemas.BlockPack, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingBlockPack, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}
	if existingBlockPack == nil {
		return nil, exceptions.BlockPack.NotFound()
	}

	if input.Values.ParentSubShelfId != nil && (input.SetNull == nil || !(*input.SetNull)["ParentSubShelfId"]) {
		subShelfRepository := NewSubShelfRepository()
		if !subShelfRepository.HasPermission(
			*input.Values.ParentSubShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Shelf.NoPermission("move a block pack to this shelf")
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlockPack)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlockPack,
		).WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockPack.NoChanges()
	}

	return &updates, nil
}

func (r *BlockPackRepository) RestoreSoftDeletedOneById(
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
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("restore a deleted block pack")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) RestoreSoftDeletedManyByIds(
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
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("restore deleted block packs")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) SoftDeleteOneById(
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
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("soft delete a block pack")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("soft delete block packs")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) HardDeleteOneById(
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
		}

		if !r.HasPermission(
			id,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("hard delete a block pack")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.BlockPack{})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) HardDeleteManyByIds(
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
		}

		if !r.HasPermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return exceptions.BlockPack.NoPermission("hard delete block packs")
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.BlockPack{})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}
