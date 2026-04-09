package repositories

import (
	"fmt"
	"sort"
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

type BlockGroupRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HasPermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionAndGetValidIds(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	GetOneByPrevBlockGroupId(blockPackId uuid.UUID, prevBlockGroupId *uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	InsertOneByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	InsertManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, inputs []inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	AppendOneByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	AppendManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input []inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockGroupInput, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockGroupRepository struct{}

func NewBlockGroupRepository() BlockGroupRepositoryInterface {
	return &BlockGroupRepository{}
}

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
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.NotFound().WithOrigin(result.Error)},
		{First: len(blockGroups) == 0, Second: exceptions.BlockGroup.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return blockGroups, nil
}

func (r *BlockGroupRepository) CheckPermissionsAndGetManyByBlockPackId(
	blockPackId uuid.UUID,
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
		Where("bp.id = ? AND EXISTS (?)",
			blockPackId, subQuery,
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.NotFound().WithOrigin(result.Error)},
		{First: len(blockGroups) == 0, Second: exceptions.BlockGroup.NotFound()},
	}); exception != nil {
		return nil, exception
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
		return make([]uuid.UUID, len(ids)), exceptions.BlockGroup.NotFound().WithOrigin(err)
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

func (r *BlockGroupRepository) GetOneByPrevBlockGroupId(
	blockPackId uuid.UUID,
	prevBlockGroupId *uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("bp.id = ? AND \"BlockGroupTable\".prev_block_group_id = ? AND EXISTS (?)",
			blockPackId, prevBlockGroupId, subQuery,
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
	if err := query.First(&blockGroup).Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	return &blockGroup, nil
}

func (r *BlockGroupRepository) InsertOneByBlockPackId(
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
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: ownerId == nil || blockPack == nil, Second: exceptions.BlockPack.NoPermission("get owner's block pack")},
	}); exception != nil {
		return nil, exception
	}

	var newBlockGroup schemas.BlockGroup
	if input.BlockGroupId != nil {
		newBlockGroup.Id = *input.BlockGroupId
	}
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId
	newBlockGroup.PrevBlockGroupId = blockPack.FinalBlockGroupId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroup)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlockGroup.Id == uuid.Nil, Second: exceptions.BlockGroup.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	if input.PrevBlockGroupId != nil {
		collapsedBlockGroup, exception := r.GetOneByPrevBlockGroupId(
			blockPackId,
			input.PrevBlockGroupId,
			userId,
			nil,
			opts...,
		)
		if exception != nil {
			return nil, exception
		}

		if _, exception = r.UpdateOneById(
			collapsedBlockGroup.Id,
			userId,
			inputs.PartialUpdateBlockGroupInput{
				Values: inputs.UpdateBlockGroupInput{
					PrevBlockGroupId: &newBlockGroup.Id,
				},
				SetNull: nil,
			},
			opts...,
		); exception != nil {
			return nil, exception
		}
		if _, exception = r.UpdateOneById(
			newBlockGroup.Id,
			userId,
			inputs.PartialUpdateBlockGroupInput{
				Values: inputs.UpdateBlockGroupInput{
					PrevBlockGroupId: input.PrevBlockGroupId,
				},
				SetNull: nil,
			},
			opts...,
		); exception != nil {
			return nil, exception
		}
	}

	return &newBlockGroup.Id, nil
}

func (r *BlockGroupRepository) InsertManyByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input []inputs.CreateBlockGroupInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)
	shouldCommit := false
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		shouldCommit = true
	}

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
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: ownerId == nil || blockPack == nil, Second: exceptions.BlockPack.NoPermission("write owner's block pack")},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	nextBlockGroups := make(map[uuid.UUID]uuid.UUID)
	newBlockGroups := make([]schemas.BlockGroup, len(input))
	newBlockGroupIds := make([]uuid.UUID, len(input))
	isNewBlockGroupId := make(map[uuid.UUID]bool, len(input))
	for index, in := range input { // use index to modify the elements of slice, since the value from the `range` is only a copy
		if in.BlockGroupId != nil {
			newBlockGroups[index].Id = *in.BlockGroupId
		} else {
			newBlockGroups[index].Id = uuid.New() // generate the id here, so that we can return the newBlockGroupIds in the same order of the input
		}
		newBlockGroups[index].OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
		newBlockGroups[index].BlockPackId = blockPackId
		newBlockGroups[index].PrevBlockGroupId = in.PrevBlockGroupId

		fakeDeletedAt := time.Now()
		newBlockGroups[index].DeletedAt = &fakeDeletedAt

		newBlockGroupIds[index] = newBlockGroups[index].Id
		isNewBlockGroupId[newBlockGroups[index].Id] = true
		if in.PrevBlockGroupId == nil {
			nextBlockGroups[uuid.Nil] = newBlockGroups[index].Id
		} else {
			nextBlockGroups[*newBlockGroups[index].PrevBlockGroupId] = newBlockGroups[index].Id
		}
	}

	var find func(id uuid.UUID) uuid.UUID
	find = func(id uuid.UUID) uuid.UUID {
		next, exist := nextBlockGroups[id]
		if !exist {
			return id
		}
		tail := find(next)
		nextBlockGroups[id] = tail
		return tail
	}
	var updatePlaceholders []string
	var updateArgs []interface{}
	var collisionPrevIds []uuid.UUID
	var isCollisionIncludingHead bool
	for _, in := range input {
		prevBlockGroupId := uuid.Nil
		if in.PrevBlockGroupId != nil {
			prevBlockGroupId = *in.PrevBlockGroupId
		}

		if isNewBlockGroupId[prevBlockGroupId] {
			continue
		}

		tail := find(prevBlockGroupId)
		if prevBlockGroupId == uuid.Nil {
			isCollisionIncludingHead = true
		} else {
			collisionPrevIds = append(collisionPrevIds, prevBlockGroupId)
		}

		updatePlaceholders = append(updatePlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
		updateArgs = append(updateArgs, prevBlockGroupId, tail, prevBlockGroupId == uuid.Nil)
	}

	// Create deleted block groups on their desired positions first,
	// so that we can avoid the block pack id and prev block group id unique index constraint
	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Create(newBlockGroups)
	if err := result.Error; err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithOrigin(err)
	}

	if len(updatePlaceholders) > 0 {
		conditions := []string{}
		var conditionArgs []interface{}

		if len(collisionPrevIds) > 0 {
			conditions = append(conditions, "prev_block_group_id IN (?)")
			conditionArgs = append(conditionArgs, collisionPrevIds)
		}
		if isCollisionIncludingHead {
			conditions = append(conditions, "(prev_block_group_id IS NULL AND block_pack_id = ?)")
			conditionArgs = append(conditionArgs, blockPackId)
		}

		if len(conditions) > 0 {
			var collidingBlockGroups []schemas.BlockGroup
			collisionQuery := parsedOptions.DB.Model(&schemas.BlockGroup{}).
				Where("deleted_at IS NULL").
				Where(strings.Join(conditions, " OR "), conditionArgs...)
			if err := collisionQuery.Find(&collidingBlockGroups).Error; err != nil {
				parsedOptions.DB.Rollback()
				return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
			}

			var colliderPlaceholders []string
			var colliderArgs []interface{}

			for _, collider := range collidingBlockGroups {
				colliderPrevKey := uuid.Nil
				if collider.PrevBlockGroupId != nil {
					colliderPrevKey = *collider.PrevBlockGroupId
				}

				tail := find(colliderPrevKey)

				colliderPlaceholders = append(colliderPlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
				colliderArgs = append(colliderArgs, collider.Id, tail, false)
			}

			if len(colliderPlaceholders) > 0 {
				sql := fmt.Sprintf(`
                    UPDATE "BlockGroupTable" AS bg
                    SET
                        prev_block_group_id = CASE
                            WHEN v.is_prev_block_group_null::boolean THEN NULL
                            ELSE v.new_prev_block_group_id::uuid
                        END,
                        updated_at = NOW()
                    FROM (VALUES %s) AS v(target_block_group_id, new_prev_block_group_id, is_prev_block_group_null)
                    WHERE bg.id = v.target_block_group_id::uuid
                `, strings.Join(colliderPlaceholders, ","))
				if err := parsedOptions.DB.Exec(sql, colliderArgs...).Error; err != nil {
					parsedOptions.DB.Rollback()
					return nil, exceptions.BlockGroup.FailedToUpdate().WithOrigin(err)
				}
			}
		}
	}

	restoreResult := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NOT NULL", newBlockGroupIds).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: restoreResult.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(restoreResult.Error)},
		{First: restoreResult.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldCommit {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newBlockGroupIds, nil
}

func (r *BlockGroupRepository) AppendOneByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
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
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: ownerId == nil || blockPack == nil, Second: exceptions.BlockPack.NoPermission("get owner's block pack")},
	}); exception != nil {
		return nil, exception
	}

	var newBlockGroup schemas.BlockGroup
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId
	newBlockGroup.PrevBlockGroupId = blockPack.FinalBlockGroupId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroup)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlockGroup.Id == uuid.Nil, Second: exceptions.BlockGroup.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &newBlockGroup.Id, nil
}

func (r *BlockGroupRepository) AppendManyByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input []inputs.CreateBlockGroupInput, // input[0] should be nil in this case
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

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
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: ownerId == nil || blockPack == nil, Second: exceptions.BlockPack.NoPermission("get owner's block pack")},
	}); exception != nil {
		return nil, exception
	}

	var newBlockGroups []schemas.BlockGroup
	if err := copier.Copy(&newBlockGroups, &input); err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithOrigin(err)
	}
	ids := make([]uuid.UUID, len(input))
	for index := range newBlockGroups { // use index to modify the elements of slice, since the value from the `range` is only a copy
		newBlockGroups[index].Id = uuid.New()    // generate the id here, so that we can return the ids in the same order of the input
		newBlockGroups[index].OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
		newBlockGroups[index].BlockPackId = blockPackId
		ids[index] = newBlockGroups[index].Id
	}
	newBlockGroups[0].PrevBlockGroupId = blockPack.FinalBlockGroupId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroups)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return nil, exception
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
		).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldCommit := false
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		shouldCommit = true
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	restoredBlockGroup, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.BlockGroup.NoPermission("restore a deleted block group")
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("block_pack_id = ? AND prev_block_group_id = ? AND deleted_at IS NULL",
			restoredBlockGroup.BlockPackId, restoredBlockGroup.PrevBlockGroupId,
		).Update("prev_block_group_id", restoredBlockGroup.Id)
	if err := result.Error; err != nil {
		// skip the error handling if there's no next block group to maintain its prev block group id
	}

	result = parsedOptions.DB.Model(&restoredBlockGroup).
		Clauses(clause.Returning{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldCommit {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return restoredBlockGroup, nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.BlockGroup, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)
	shouldCommit := false
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		shouldCommit = true
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockGroups, exception := r.CheckPermissionsAndGetManyByIds(
		ids,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	prevBlockGroupIdMap := make(map[uuid.UUID][]schemas.BlockGroup)
	for _, blockGroup := range blockGroups {
		key := uuid.Nil
		if blockGroup.PrevBlockGroupId != nil {
			key = *blockGroup.PrevBlockGroupId
		}
		prevBlockGroupIdMap[key] = append(prevBlockGroupIdMap[key], blockGroup)
	}

	var updatePlaceholders []string
	var updateArgs []interface{}

	var collidingPrevBlockGroupIds []uuid.UUID // for those block groups with prev block group id not equal to nil
	var collidingBlockPackIds []uuid.UUID      // for those block groups with prev block group id of nil,
	// this should have at most one element, but we extends this for restoring soft deleted block groups across any given block packs in the future

	restoredChainTails := make(map[uuid.UUID]uuid.UUID)
	restoredHeadChainTails := make(map[uuid.UUID]uuid.UUID)
	nodeToChainTail := make(map[uuid.UUID]uuid.UUID)

	for prevBlockGroupId, groupedBlockGroups := range prevBlockGroupIdMap {
		if len(groupedBlockGroups) == 0 {
			continue
		}

		sort.Slice(groupedBlockGroups, func(i, j int) bool {
			if groupedBlockGroups[i].DeletedAt == nil {
				return false
			}
			if groupedBlockGroups[j].DeletedAt == nil {
				return true
			}
			return groupedBlockGroups[i].DeletedAt.Before(*groupedBlockGroups[j].DeletedAt)
		})

		for i := 1; i < len(groupedBlockGroups); i++ {
			updatePlaceholders = append(updatePlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
			updateArgs = append(updateArgs, groupedBlockGroups[i].Id, groupedBlockGroups[i-1].Id, false)
		}

		tailBlockGroup := groupedBlockGroups[len(groupedBlockGroups)-1]

		for _, bg := range groupedBlockGroups {
			nodeToChainTail[bg.Id] = tailBlockGroup.Id
		}

		if prevBlockGroupId != uuid.Nil {
			collidingPrevBlockGroupIds = append(collidingPrevBlockGroupIds, prevBlockGroupId)
			restoredChainTails[prevBlockGroupId] = tailBlockGroup.Id
		} else {
			blockPackId := tailBlockGroup.BlockPackId
			collidingBlockPackIds = append(collidingBlockPackIds, blockPackId)
			restoredHeadChainTails[uuid.Nil] = tailBlockGroup.Id
		}
	}

	var collidingBlockGroups []schemas.BlockGroup
	collisionQuery := parsedOptions.DB.Model(&schemas.BlockGroup{}).Where("deleted_at IS NULL")

	conditions := []string{}
	var conditionArgs []interface{}

	if len(collidingPrevBlockGroupIds) > 0 {
		conditions = append(conditions, "prev_block_group_id IN (?)")
		conditionArgs = append(conditionArgs, collidingPrevBlockGroupIds)
	}
	if len(collidingBlockPackIds) > 0 {
		conditions = append(conditions, "(prev_block_group_id IS NULL AND block_pack_id IN (?))")
		conditionArgs = append(conditionArgs, collidingBlockPackIds)
	}
	if len(ids) > 0 {
		conditions = append(conditions, "prev_block_group_id IN (?)")
		conditionArgs = append(conditionArgs, ids)
	}

	if len(conditions) > 0 {
		collisionQuery = collisionQuery.Where(strings.Join(conditions, " OR "), conditionArgs...)
		if err := collisionQuery.Find(&collidingBlockGroups).Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
		}
	}

	for _, collider := range collidingBlockGroups {
		var tailId uuid.UUID
		found := false

		if collider.PrevBlockGroupId != nil {
			if tid, ok := restoredChainTails[*collider.PrevBlockGroupId]; ok {
				tailId = tid
				found = true
			}
			if !found {
				if tid, ok := nodeToChainTail[*collider.PrevBlockGroupId]; ok {
					tailId = tid
					found = true
				}
			}
		} else {
			if tid, ok := restoredHeadChainTails[uuid.Nil]; ok {
				tailId = tid
				found = true
			}
		}

		if found {
			updatePlaceholders = append(updatePlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
			updateArgs = append(updateArgs, collider.Id, tailId, false)
		}
	}

	if len(updatePlaceholders) > 0 {
		sql := fmt.Sprintf(`
            UPDATE "BlockGroupTable" AS bg
            SET
                prev_block_group_id = CASE 
                    WHEN v.is_prev_block_group_null::boolean THEN NULL 
                    ELSE v.new_prev_block_group_id::uuid 
                END, 
                updated_at = NOW()
            FROM (VALUES %s) AS v(target_block_group_id, new_prev_block_group_id, is_prev_block_group_null)
            WHERE bg.id = v.target_block_group_id::uuid
        `, strings.Join(updatePlaceholders, ","))
		if err := parsedOptions.DB.Exec(sql, updateArgs...).Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToUpdate().WithOrigin(err)
		}
	}

	var restoredBlockGroups []schemas.BlockGroup
	result := parsedOptions.DB.Model(&restoredBlockGroups).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldCommit {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return restoredBlockGroups, nil
}

func (r *BlockGroupRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	return r.SoftDeleteManyByIds([]uuid.UUID{id}, userId, opts...)
}

func (r *BlockGroupRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldCommit := false
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		shouldCommit = true
	}

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		if !r.HasPermissions(ids, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.NoPermission("soft delete block groups")
		}
	}

	var targets []struct {
		Id               uuid.UUID
		PrevBlockGroupId *uuid.UUID
		BlockPackId      uuid.UUID
	}
	if err := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Scan(&targets).Error; err != nil {
		parsedOptions.DB.Rollback()
		return exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	if len(targets) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.BlockGroup.NotFound()
	}

	deletionSet := make(map[uuid.UUID]bool)   // ID -> IsDeleting
	prevMap := make(map[uuid.UUID]*uuid.UUID) // ID -> PrevID (Adjacency List)

	for _, target := range targets {
		deletionSet[target.Id] = true
		prevMap[target.Id] = target.PrevBlockGroupId
	}

	memoEffectivePrev := make(map[uuid.UUID]*uuid.UUID)

	var getEffectivePrev func(nodeId uuid.UUID) *uuid.UUID
	getEffectivePrev = func(nodeId uuid.UUID) *uuid.UUID {
		if res, ok := memoEffectivePrev[nodeId]; ok {
			return res
		}

		originalPrev := prevMap[nodeId]

		if originalPrev == nil {
			memoEffectivePrev[nodeId] = nil
			return nil
		}

		if !deletionSet[*originalPrev] {
			memoEffectivePrev[nodeId] = originalPrev
			return originalPrev
		}

		ancestorPrev := getEffectivePrev(*originalPrev)
		memoEffectivePrev[nodeId] = ancestorPrev
		return ancestorPrev
	}

	var valuePlaceholders []string
	var valueArgs []interface{}

	for _, target := range targets {
		effectivePrev := getEffectivePrev(target.Id)

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
		isNewPrevNull := effectivePrev == nil
		var newPrevVal interface{}
		if effectivePrev != nil {
			newPrevVal = *effectivePrev
		} else {
			newPrevVal = nil
		}

		valueArgs = append(valueArgs, target.Id, newPrevVal, isNewPrevNull)
	}

	if len(valuePlaceholders) > 0 {
		sql := fmt.Sprintf(`
            UPDATE "BlockGroupTable" AS bg
            SET
                prev_block_group_id = CASE 
                    WHEN v.is_prev_block_group_id_null::boolean THEN NULL 
                    ELSE v.new_prev_block_group_id::uuid 
                END, 
                updated_at = NOW()
            FROM (VALUES %s) AS v(old_prev_block_group_id, new_prev_block_group_id, is_prev_block_group_id_null)
            WHERE bg.prev_block_group_id = v.old_prev_block_group_id::uuid AND bg.deleted_at IS NULL
        `, strings.Join(valuePlaceholders, ","))
		result := parsedOptions.DB.Exec(sql, valueArgs...)
		if result.Error != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)
		}
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ?", ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldCommit {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			return exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *BlockGroupRepository) HardDeleteManyByIds(
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
