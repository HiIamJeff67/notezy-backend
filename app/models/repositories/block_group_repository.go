package repositories

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	blockgroupsql "github.com/HiIamJeff67/notezy-backend/app/models/sqls/block_group"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockGroupRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionAndGetValidIds(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	CollectOrphanedBlockGroupsByIds(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) *exceptions.Exception
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	GetOneByPrevBlockGroupId(blockPackId uuid.UUID, prevBlockGroupId *uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	GetManyByPrevBlockGroupIds(BlockPackIds []uuid.UUID, PrevBlockGroupIds []*uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	InsertOneByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	InsertManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input []inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	InsertManyByBlockPackIds(userId uuid.UUID, input []inputs.BulkCreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	AppendOneByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	AppendManyByBlockPackId(blockPackId uuid.UUID, userId uuid.UUID, input []inputs.CreateBlockGroupInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockGroupInput, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateBlockGroupsInput, opts ...options.RepositoryOptions) *exceptions.Exception
	IncrementSizeById(id uuid.UUID, userId uuid.UUID, sizeDelta int64, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	IncrementSizesByIds(userId uuid.UUID, sizeDeltaById map[uuid.UUID]int64, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.BlockGroup, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.BlockGroup, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type BlockGroupRepository struct {
	blockGroupScope scopes.BlockGroupScopeInterface
}

func NewBlockGroupRepository(blockGroupScope scopes.BlockGroupScopeInterface) BlockGroupRepositoryInterface {
	return &BlockGroupRepository{
		blockGroupScope: blockGroupScope,
	}
}

func (r *BlockGroupRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.BlockGroup{}).
		Select("1").
		Scopes(r.blockGroupScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *BlockGroupRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.BlockGroup{}).
		Select(`DISTINCT "BlockGroupTable".id`).
		Scopes(r.blockGroupScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *BlockGroupRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var blockGroup schemas.BlockGroup
	result := parsedOptions.DB.
		Model(&schemas.BlockGroup{}).
		Scopes(r.blockGroupScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockGroupScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&blockGroup)
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

	var blockGroups []schemas.BlockGroup
	result := parsedOptions.DB.
		Model(&schemas.BlockGroup{}).
		Scopes(r.blockGroupScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockGroupScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&blockGroups)
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
		Joins(`INNER JOIN "BlockPackTable" bp ON block_pack_id = bp.id`).
		Joins(`INNER JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id`).
		Where("bp.id = ? AND EXISTS (?)",
			blockPackId, subQuery,
		)

	var blockGroups []schemas.BlockGroup
	result := query.
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockGroupScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&blockGroups)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.NotFound().WithOrigin(result.Error)},
		{First: len(blockGroups) == 0, Second: exceptions.BlockGroup.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return blockGroups, nil
}

// Similar to the `HavePermissions`, but with best effort strategy,
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
		Joins(`INNER JOIN "BlockPackTable" bp ON block_pack_id = bp.id`).
		Joins(`INNER JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id`).
		Where(`"BlockGroupTable".id IN ? AND EXISTS (?)`,
			ids, subQuery,
		)

	var validIds []uuid.UUID
	if err := query.
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Scan(&validIds).Error; err != nil {
		return make([]uuid.UUID, len(ids)), exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	return validIds, nil
}

func (r *BlockGroupRepository) CollectOrphanedBlockGroupsByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	var orphanedBlockGroupIds []uuid.UUID
	result := parsedOptions.DB.Raw(
		blockgroupsql.GetGarbageCollectedOrphanedBlockGroupIdsSQL,
		ids,
		userId,
		allowedPermissions,
	).Scan(&orphanedBlockGroupIds)
	if result.Error != nil {
		parsedOptions.DB.Rollback()
		return exceptions.BlockGroup.NotFound().WithOrigin(result.Error)
	}

	if len(orphanedBlockGroupIds) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.BlockGroup.NoChanges()
	}

	if exception := r.SoftDeleteManyByIds(
		orphanedBlockGroupIds,
		userId,
		opts...,
	); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
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
		Joins(`INNER JOIN "BlockPackTable" bp ON block_pack_id = bp.id`).
		Joins(`INNER JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id`).
		Where(`bp.id = ? AND "BlockGroupTable".prev_block_group_id = ? AND EXISTS (?)`,
			blockPackId, prevBlockGroupId, subQuery,
		)

	var blockGroup schemas.BlockGroup
	if err := query.
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockGroupScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&blockGroup).Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	return &blockGroup, nil
}

func (r *BlockGroupRepository) GetManyByPrevBlockGroupIds(
	BlockPackIds []uuid.UUID,
	PrevBlockGroupIds []*uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.BlockGroup, *exceptions.Exception) {
	if len(BlockPackIds) != len(PrevBlockGroupIds) {
		return []schemas.BlockGroup{}, nil
	}

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
		Where("user_id = ? AND permission IN ?", userId, allowedPermissions)
	query := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Joins(`INNER JOIN "BlockPackTable" bp ON block_pack_id = bp.id`).
		Joins(`INNER JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id`).
		Where("EXISTS (?)", subQuery)

	var nilPrevBlockPackIds []uuid.UUID
	var nonNilConditions []string
	var nonNilArgs []interface{}

	for index, _ := range BlockPackIds {
		if PrevBlockGroupIds[index] == nil {
			nilPrevBlockPackIds = append(nilPrevBlockPackIds, BlockPackIds[index])
		} else {
			nonNilConditions = append(nonNilConditions, `("BlockGroupTable".block_pack_id = ? AND "BlockGroupTable".prev_block_group_id = ?)`)
			nonNilArgs = append(nonNilArgs, BlockPackIds[index], *PrevBlockGroupIds[index])
		}
	}

	var combinedConditions []string
	var combinedArgs []interface{}

	if len(nilPrevBlockPackIds) > 0 {
		combinedConditions = append(combinedConditions, `("BlockGroupTable".block_pack_id IN ? AND "BlockGroupTable".prev_block_group_id IS NULL)`)
		combinedArgs = append(combinedArgs, nilPrevBlockPackIds)
	}

	if len(nonNilConditions) > 0 {
		combinedConditions = append(combinedConditions, strings.Join(nonNilConditions, " OR "))
		combinedArgs = append(combinedArgs, nonNilArgs...)
	}

	if len(combinedConditions) > 0 {
		query = query.Where(fmt.Sprintf("(%s)", strings.Join(combinedConditions, " OR ")), combinedArgs...)
	}

	var blockGroups []schemas.BlockGroup
	result := query.
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.blockGroupScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&blockGroups)
	if result.Error != nil {
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(result.Error)
	}

	return blockGroups, nil
}

func (r *BlockGroupRepository) InsertOneByBlockPackId(
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockGroupInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository(scopes.NewBlockPackScope())

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
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	var newBlockGroup schemas.BlockGroup
	if input.BlockGroupId != nil {
		newBlockGroup.Id = *input.BlockGroupId
	}
	if newBlockGroup.Id == uuid.Nil {
		newBlockGroup.Id = uuid.New()
	}
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId
	newBlockGroup.PrevBlockGroupId = blockPack.FinalBlockGroupId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Create(&newBlockGroup)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlockGroup.Id == uuid.Nil, Second: exceptions.BlockGroup.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
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
			parsedOptions.DB.Rollback()
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
			parsedOptions.DB.Rollback()
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
			parsedOptions.DB.Rollback()
			return nil, exception
		}
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository(scopes.NewBlockPackScope())

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
	fakeDeletedAt := time.Now()
	for index, in := range input { // use index to modify the elements of slice, since the value from the `range` is only a copy
		if in.BlockGroupId != nil && *in.BlockGroupId != uuid.Nil {
			newBlockGroups[index].Id = *in.BlockGroupId
		} else {
			newBlockGroups[index].Id = uuid.New() // generate the id here, so that we can return the newBlockGroupIds in the same order of the input
		}
		newBlockGroups[index].OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
		newBlockGroups[index].BlockPackId = blockPackId
		newBlockGroups[index].PrevBlockGroupId = in.PrevBlockGroupId
		newBlockGroups[index].DeletedAt = &fakeDeletedAt

		newBlockGroupIds[index] = newBlockGroups[index].Id
		isNewBlockGroupId[newBlockGroups[index].Id] = true
		if in.PrevBlockGroupId == nil {
			nextBlockGroups[uuid.Nil] = newBlockGroups[index].Id
		} else {
			nextBlockGroups[*newBlockGroups[index].PrevBlockGroupId] = newBlockGroups[index].Id
		}
	}

	// Create deleted block groups on their desired positions first,
	// so that we can avoid the block pack id and prev block group id unique index constraint
	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		CreateInBatches(newBlockGroups, parsedOptions.BatchSize)
	if err := result.Error; err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithOrigin(err)
	}

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

		if prevBlockGroupId == uuid.Nil {
			isCollisionIncludingHead = true
		} else {
			collisionPrevIds = append(collisionPrevIds, prevBlockGroupId)
		}
	}

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

	if len(conditions) > 0 && len(conditions) == len(conditionArgs) {
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
			colliderArgs = append(colliderArgs, collider.Id, tail, tail == uuid.Nil)
		}

		if len(colliderPlaceholders) > 0 && len(colliderArgs) > 0 {
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

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newBlockGroupIds, nil
}

func (r *BlockGroupRepository) InsertManyByBlockPackIds(
	userId uuid.UUID,
	input []inputs.BulkCreateBlockGroupInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository(scopes.NewBlockPackScope())

	blockPackIds := make([]uuid.UUID, len(input))
	for index, in := range input {
		blockPackIds[index] = in.BlockPackId
	}

	ownerIds, validBlockPacks, exception := blockPackRepository.CheckPermissionsAndGetManyWithOwnerIdsByIds(
		blockPackIds,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: len(ownerIds) == 0 || len(validBlockPacks) == 0, Second: exceptions.BlockPack.NoPermission("write all owner's block packs")},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	validBlockPackIdToOwnerId := make(map[uuid.UUID]uuid.UUID)
	for index, _ := range validBlockPacks {
		validBlockPackIdToOwnerId[validBlockPacks[index].Id] = ownerIds[index]
	}

	nextBlockGroups := make(map[uuid.UUID]uuid.UUID)
	newBlockGroups := make([]schemas.BlockGroup, len(input))
	newBlockGroupIds := make([]uuid.UUID, len(input))
	isNewBlockGroupId := make(map[uuid.UUID]bool)
	blockPackIdToBlockGroups := make(map[uuid.UUID][]schemas.BlockGroup)
	fakeDeletedAt := time.Now()
	for index, in := range input {
		ownerId, exist := validBlockPackIdToOwnerId[in.BlockPackId]
		if !exist {
			continue // best effort
		}

		if in.BlockGroupId != nil && *in.BlockGroupId != uuid.Nil {
			newBlockGroups[index].Id = *in.BlockGroupId
		} else {
			newBlockGroups[index].Id = uuid.New()
		}
		newBlockGroups[index].OwnerId = ownerId
		newBlockGroups[index].BlockPackId = in.BlockPackId
		newBlockGroups[index].PrevBlockGroupId = in.PrevBlockGroupId
		newBlockGroups[index].DeletedAt = &fakeDeletedAt

		newBlockGroupIds[index] = newBlockGroups[index].Id
		isNewBlockGroupId[newBlockGroups[index].Id] = true
		if in.PrevBlockGroupId == nil {
			nextBlockGroups[uuid.Nil] = newBlockGroups[index].Id
		} else {
			nextBlockGroups[*in.PrevBlockGroupId] = newBlockGroups[index].Id
		}
		blockPackIdToBlockGroups[in.BlockPackId] = append(blockPackIdToBlockGroups[in.BlockPackId], newBlockGroups[index])
	}

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		CreateInBatches(newBlockGroups, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	conditions := []string{}
	var conditionArgs []interface{}
	for _, validBlockPack := range validBlockPacks {
		var collisionPrevIds []uuid.UUID
		var isCollisionIncludingHead bool
		for _, blockGroup := range blockPackIdToBlockGroups[validBlockPack.Id] {
			prevBlockGroupId := uuid.Nil
			if blockGroup.PrevBlockGroupId != nil {
				prevBlockGroupId = *blockGroup.PrevBlockGroupId
			}

			if isNewBlockGroupId[prevBlockGroupId] {
				continue
			}

			if prevBlockGroupId == uuid.Nil {
				isCollisionIncludingHead = true
			} else {
				collisionPrevIds = append(collisionPrevIds, prevBlockGroupId)
			}
		}

		if len(collisionPrevIds) > 0 {
			conditions = append(conditions, "prev_block_group_id IN (?)")
			conditionArgs = append(conditionArgs, collisionPrevIds)
		}
		if isCollisionIncludingHead {
			conditions = append(conditions, "(prev_block_group_id IS NULL AND block_pack_id = ?)")
			conditionArgs = append(conditionArgs, validBlockPack.Id)
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

	if len(conditions) > 0 && len(conditions) == len(conditionArgs) {
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
			colliderArgs = append(colliderArgs, collider.Id, tail, tail == uuid.Nil)
		}

		if len(colliderPlaceholders) > 0 && len(colliderPlaceholders) == len(colliderArgs) {
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
			result := parsedOptions.DB.Exec(sql, colliderArgs...)
			if exception = exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
				{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
				{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
			}); exception != nil {
				parsedOptions.DB.Rollback()
				return nil, exception
			}
		}
	}

	restoredResult := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NOT NULL", newBlockGroupIds).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception = exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: restoredResult.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(restoredResult.Error)},
		{First: restoredResult.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository(scopes.NewBlockPackScope())

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
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	var newBlockGroup schemas.BlockGroup
	newBlockGroup.Id = uuid.New()
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId
	newBlockGroup.PrevBlockGroupId = blockPack.FinalBlockGroupId

	result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
		Create(&newBlockGroup)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToCreate().WithOrigin(result.Error)},
		{First: newBlockGroup.Id == uuid.Nil, Second: exceptions.BlockGroup.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository(scopes.NewBlockPackScope())

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
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	var newBlockGroups []schemas.BlockGroup
	if err := copier.Copy(&newBlockGroups, &input); err != nil {
		parsedOptions.DB.Rollback()
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
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

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
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlockGroup)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			existingBlockGroup,
		).WithOrigin(err)
	}

	result := parsedOptions.DB.
		Model(&schemas.BlockGroup{}).
		Scopes(r.blockGroupScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *BlockGroupRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateBlockGroupsInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	isBlockGroupValid := make(map[uuid.UUID]bool)
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

		validBlockGroups, exception := r.CheckPermissionsAndGetManyByIds(
			ids,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.NoPermission("update these block groups")
		}

		for _, validBlockGroup := range validBlockGroups {
			isBlockGroupValid[validBlockGroup.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isBlockGroupValid[in.Id] {
			continue
		}

		setPrevBlockGroupIdNull := false
		if in.PartialUpdateInput.SetNull != nil {
			for field, setNull := range *in.PartialUpdateInput.SetNull {
				if strings.ToLower(field) == "prevblockgroupid" && setNull {
					setPrevBlockGroupIdNull = true
					break
				}
			}
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::bigint, ?::boolean)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.PrevBlockGroupId,
			in.PartialUpdateInput.Values.Size,
			setPrevBlockGroupIdNull,
		)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.BlockGroup.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "BlockGroupTable" AS bg
		SET
			prev_block_group_id = CASE
				WHEN v.set_prev_block_group_id_null::boolean THEN NULL
				ELSE COALESCE(v.prev_block_group_id::uuid, bg.prev_block_group_id)
			END,
			size = GREATEST(0, COALESCE(v.size::bigint, bg.size)),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, prev_block_group_id, size, set_prev_block_group_id_null)
		WHERE bg.id = v.id::uuid AND bg.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *BlockGroupRepository) IncrementSizeById(
	id uuid.UUID,
	userId uuid.UUID,
	sizeDelta int64,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	exception := r.IncrementSizesByIds(
		userId,
		map[uuid.UUID]int64{
			id: sizeDelta,
		},
		opts...,
	)
	if exception != nil {
		return nil, exception
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	return r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		options.WithDB(parsedOptions.DB),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
}

func (r *BlockGroupRepository) IncrementSizesByIds(
	userId uuid.UUID,
	sizeDeltaById map[uuid.UUID]int64,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(sizeDeltaById) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	isBlockGroupValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		ids := make([]uuid.UUID, 0, len(sizeDeltaById))
		for id := range sizeDeltaById {
			ids = append(ids, id)
		}

		validBlockGroups, exception := r.CheckPermissionsAndGetManyByIds(
			ids,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.NoPermission("update these block groups")
		}

		for _, validBlockGroup := range validBlockGroups {
			isBlockGroupValid[validBlockGroup.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for id, sizeDelta := range sizeDeltaById {
		if sizeDelta == 0 {
			continue
		}
		if !parsedOptions.SkipPermissionCheck && !isBlockGroupValid[id] {
			continue
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::bigint)")
		valueArgs = append(valueArgs, id, sizeDelta)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		if !parsedOptions.SkipPermissionCheck {
			return exceptions.BlockGroup.NoPermission("update these block groups")
		}
		return exceptions.BlockGroup.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "BlockGroupTable" AS bg
		SET
			size = GREATEST(0, bg.size + v.size_delta::bigint),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, size_delta)
		WHERE bg.id = v.id::uuid AND bg.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.BlockGroup, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
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

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
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

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockGroups, exception := r.CheckPermissionsAndGetManyByIds(
		ids,
		userId,
		[]schemas.BlockGroupRelation{
			schemas.BlockGroupRelation_NextBlockGroup,
		},
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	isBlockGroupDeleted := make(map[uuid.UUID]bool)
	blockGroupIdToPrevBlockGroupId := make(map[uuid.UUID]*uuid.UUID)
	blockGroupIdToNextBlockGroupId := make(map[uuid.UUID]*uuid.UUID)
	for _, blockGroup := range blockGroups {
		isBlockGroupDeleted[blockGroup.Id] = true
		blockGroupIdToPrevBlockGroupId[blockGroup.Id] = blockGroup.PrevBlockGroupId
		if blockGroup.NextBlockGroup == nil {
			blockGroupIdToNextBlockGroupId[blockGroup.Id] = nil
		} else {
			blockGroupIdToNextBlockGroupId[blockGroup.Id] = &blockGroup.NextBlockGroup.Id
		}
	}

	ancestors := make(map[uuid.UUID]*uuid.UUID)
	var findAncestor func(id uuid.UUID) *uuid.UUID
	findAncestor = func(id uuid.UUID) *uuid.UUID {
		if ancestor, exist := ancestors[id]; exist {
			return ancestor
		}

		originalPrevBlockGroupId := blockGroupIdToPrevBlockGroupId[id]
		if originalPrevBlockGroupId == nil {
			ancestors[id] = nil
			return nil
		}

		if !isBlockGroupDeleted[*originalPrevBlockGroupId] {
			ancestors[id] = originalPrevBlockGroupId
			return originalPrevBlockGroupId
		}

		ancestor := findAncestor(*originalPrevBlockGroupId)
		ancestors[id] = ancestor
		return ancestor
	}

	descendants := make(map[uuid.UUID]*uuid.UUID)
	var findDescendant func(id uuid.UUID) *uuid.UUID
	findDescendant = func(id uuid.UUID) *uuid.UUID {
		if descendant, exist := descendants[id]; exist {
			return descendant
		}

		originalNextBlockGroupId := blockGroupIdToNextBlockGroupId[id]
		if originalNextBlockGroupId == nil {
			descendants[id] = nil
			return nil
		}

		if !isBlockGroupDeleted[*originalNextBlockGroupId] {
			descendants[id] = originalNextBlockGroupId
			return originalNextBlockGroupId
		}

		descendant := findDescendant(*originalNextBlockGroupId)
		descendants[id] = descendant
		return descendant
	}

	isPlaced := make(map[string]bool)
	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, blockGroup := range blockGroups {
		descendant := findDescendant(blockGroup.Id)
		if descendant == nil {
			continue
		}
		descendantString := descendant.String()

		ancestor := findAncestor(blockGroup.Id)
		ancestorString := "nil"
		if ancestor != nil {
			ancestorString = ancestor.String()
		}

		if isPlaced[ancestorString+descendantString] {
			continue
		}
		isPlaced[ancestorString+descendantString] = true

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::boolean)")
		valueArgs = append(valueArgs, descendant, ancestor, ancestor == nil)
	}

	// delete first so that we can avoid the unique index of block group id and prev block group id while relink the descendants to their ancestors below
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

	if len(valuePlaceholders) > 0 && len(valueArgs) > 0 {
		sql := fmt.Sprintf(`
            UPDATE "BlockGroupTable" AS bg
            SET
                prev_block_group_id = CASE 
                    WHEN v.is_prev_block_group_id_null::boolean THEN NULL 
                    ELSE v.prev_block_group_id::uuid 
                END, 
                updated_at = NOW()
            FROM (VALUES %s) AS v(id, prev_block_group_id, is_prev_block_group_id_null)
            WHERE bg.id = v.id::uuid AND bg.deleted_at IS NULL
        `, strings.Join(valuePlaceholders, ","))
		result := parsedOptions.DB.Exec(sql, valueArgs...)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return exception
		}
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
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

		result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
			Where("id = ? AND deleted_at IS NOT NULL", id).
			Scopes(r.blockGroupScope.PassPermissionCheck(id, userId, allowedPermissions)).
			Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
			Delete(&schemas.BlockGroup{})
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToDelete().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
		}); exception != nil {
			return exception
		}

		return nil
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

		result := parsedOptions.DB.Model(&schemas.BlockGroup{}).
			Where("id IN ? AND deleted_at IS NOT NULL", ids).
			Scopes(r.blockGroupScope.PassPermissionChecks(ids, userId, allowedPermissions)).
			Scopes(r.blockGroupScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
			Delete(&schemas.BlockGroup{})
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToDelete().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
		}); exception != nil {
			return exception
		}

		return nil
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
