package repositories

import (
	"fmt"
	"strings"
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

type BlockPackRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockPack, *exceptions.Exception)
	CheckPermissionAndGetOneWithOwnerIdById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*uuid.UUID, *schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyWithOwnerIdsByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]uuid.UUID, []schemas.BlockPack, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	CreateOneBySubShelfId(subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockPackInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyBySubShelfIds(userId uuid.UUID, input []inputs.BulkCreateBlockPackInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockPackInput, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateBlockPackInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.BlockPack, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.BlockPack, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockPackRepository struct{}

func NewBlockPackRepository() BlockPackRepositoryInterface {
	return &BlockPackRepository{}
}

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

func (r *BlockPackRepository) HavePermissions(
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
		Where("\"UsersToShelvesTable\".root_shelf_id = ss.root_shelf_id").
		Where("\"UsersToShelvesTable\".user_id = ? AND \"UsersToShelvesTable\".permission IN ?",
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

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPack schemas.BlockPack
	result := query.First(&blockPack)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.NotFound().WithOrigin(result.Error)},
		{First: blockPack.Id == uuid.Nil, Second: exceptions.BlockPack.NotFound()},
	}); exception != nil {
		return nil, exception
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
		Where("\"UsersToShelvesTable\".root_shelf_id = ss.root_shelf_id").
		Where("\"UsersToShelvesTable\".user_id = ? AND \"UsersToShelvesTable\".permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
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

	var blockPacks []schemas.BlockPack
	result := query.Find(&blockPacks)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.NotFound().WithOrigin(result.Error)},
		{First: len(blockPacks) == 0, Second: exceptions.BlockPack.NotFound()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.NotFound().WithOrigin(result.Error)},
		{First: blockPackWithOwnerId.OwnerId == uuid.Nil, Second: exceptions.BlockPack.NotFound()},
	}); exception != nil {
		return nil, nil, exception
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
		Where("\"UsersToShelvesTable\".root_shelf_id = ss.root_shelf_id").
		Where("\"UsersToShelvesTable\".user_id = ? AND \"UsersToShelvesTable\".permission IN ?",
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.NotFound().WithOrigin(result.Error)},
		{First: len(blockPacksWithOwnerIds) == 0, Second: exceptions.BlockPack.NotFound()},
	}); exception != nil {
		return nil, nil, exception
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
		return nil, exceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}
	newBlockPack.ParentSubShelfId = subShelfId

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockPack)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlockPack.Id == uuid.Nil, Second: exceptions.BlockPack.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &newBlockPack.Id, nil
}

func (r *BlockPackRepository) BulkCreateManyBySubShelfIds(
	userId uuid.UUID,
	input []inputs.BulkCreateBlockPackInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	isParentSubShelfIdValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		isParentSubShelfExist := make(map[uuid.UUID]bool)
		var parentSubShelfIds []uuid.UUID
		for _, in := range input {
			if isParentSubShelfExist[in.ParentSubShelfId] {
				continue
			}
			isParentSubShelfExist[in.ParentSubShelfId] = true
			parentSubShelfIds = append(parentSubShelfIds, in.ParentSubShelfId)
		}

		subShelfRepository := NewSubShelfRepository()
		validParentSubShelves, exception := subShelfRepository.CheckPermissionsAndGetManyByIds(
			parentSubShelfIds,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			return nil, exception
		}

		for _, validParentSubShelf := range validParentSubShelves {
			isParentSubShelfIdValid[validParentSubShelf.Id] = true
		}
	}

	var newBlockPacks []schemas.BlockPack
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isParentSubShelfIdValid[in.ParentSubShelfId] {
			continue
		}
		var newBlockPack schemas.BlockPack
		if err := copier.Copy(&newBlockPack, &in); err != nil {
			return nil, exceptions.BlockPack.InvalidInput().WithOrigin(err)
		}
		newBlockPacks = append(newBlockPacks, newBlockPack)
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newBlockPacks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Block.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	newBlockPackIds := make([]uuid.UUID, len(newBlockPacks))
	for index, newBlockPack := range newBlockPacks {
		newBlockPackIds[index] = newBlockPack.Id
	}

	return newBlockPackIds, nil
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
		).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}

func (r *BlockPackRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateBlockPackInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	isSubShelfValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		isParentSubShelfExist := make(map[uuid.UUID]bool)
		var parentSubShelfIds []uuid.UUID
		for _, in := range input {
			if in.PartialUpdateInput.Values.ParentSubShelfId == nil {
				if isParentSubShelfExist[uuid.Nil] {
					continue
				}

				parentSubShelfIds = append(parentSubShelfIds, uuid.Nil)
				isParentSubShelfExist[uuid.Nil] = true
			} else {
				if isParentSubShelfExist[*in.PartialUpdateInput.Values.ParentSubShelfId] {
					continue
				}

				parentSubShelfIds = append(parentSubShelfIds, *in.PartialUpdateInput.Values.ParentSubShelfId)
				isParentSubShelfExist[*in.PartialUpdateInput.Values.ParentSubShelfId] = true
			}
		}

		subShelfRepository := NewSubShelfRepository()
		validSubShelves, exception := subShelfRepository.CheckPermissionsAndGetManyByIds(
			parentSubShelfIds,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			return exception
		}

		for _, validSubShelf := range validSubShelves {
			isSubShelfValid[validSubShelf.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck &&
			in.PartialUpdateInput.Values.ParentSubShelfId != nil &&
			!isSubShelfValid[*in.PartialUpdateInput.Values.ParentSubShelfId] {
			continue
		}

		setIconNull := false
		setHeaderBackgroundNull := false
		if in.PartialUpdateInput.SetNull != nil {
			for field, setNull := range *in.PartialUpdateInput.SetNull {
				if setNull {
					switch strings.ToLower(field) {
					case "icon":
						setIconNull = true
					case "headerbackgroundurl":
						setHeaderBackgroundNull = true
					}
				}
				if setIconNull && setHeaderBackgroundNull {
					break
				}
			}
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::string, ?::\"SupportedBlockPackIcon\", ?::string, ?::boolean, ?::boolean)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.ParentSubShelfId,
			in.PartialUpdateInput.Values.Name,
			in.PartialUpdateInput.Values.Icon,
			in.PartialUpdateInput.Values.HeaderBackgroundURL,
			setIconNull,
			setHeaderBackgroundNull,
		)
	}

	sql := fmt.Sprintf(`
		UPDATE "BlockPackTable" bp
		SET
			parent_sub_shelf_id = COALESCE(v.parent_sub_shelf_id::uuid, bp.parent_sub_shelf_id),
			name = COALESCE(v.name::string, bp.name),
			icon = CASE
				WHEN v.set_icon_null::boolean THEN NULL
				ELSE COALESCE(v.icon::"SupportedBlockPackIcon", bp.icon)
			END,
			header_background_url = CASE
				WHEN v.set_header_background_url_null::boolean THEN NULL
				ELSE COALESCE(v.header_background_url::string, bp.header_background_url)
			END,
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, parent_sub_shelf_id, name, icon, header_background_url, set_icon_null, set_header_background_url_null)
		WHERE bp.id = v.id::uuid AND bp.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *BlockPackRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.BlockPack, *exceptions.Exception) {
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
			return nil, exceptions.BlockPack.NoPermission("restore a deleted block pack")
		}
	}

	var restoredBlockPack schemas.BlockPack
	result := parsedOptions.DB.Model(&restoredBlockPack).
		Clauses(clause.Returning{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredBlockPack, nil
}

func (r *BlockPackRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.BlockPack, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HavePermissions(
			ids,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.BlockPack.NoPermission("restore deleted block packs")
		}
	}

	var restoredBlockPacks []schemas.BlockPack
	result := parsedOptions.DB.Model(restoredBlockPacks).
		Clauses(&clause.Returning{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredBlockPacks, nil
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return exception
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

		if !r.HavePermissions(
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return exception
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

		if !r.HavePermissions(
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
