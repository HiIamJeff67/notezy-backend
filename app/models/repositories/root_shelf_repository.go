package repositories

import (
	"fmt"
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
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RootShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RootShelf, enums.AccessControlPermission, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RootShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RootShelf, []enums.AccessControlPermission, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RootShelfRelation, opts ...options.RepositoryOptions) (*schemas.RootShelf, enums.AccessControlPermission, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateRootShelfInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	CreateManyByOwnerId(ownerId uuid.UUID, input []inputs.CreateRootShelfInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRootShelfInput, opts ...options.RepositoryOptions) (*schemas.RootShelf, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateRootShelfInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.RootShelf, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.RootShelf, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(sids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type RootShelfRepository struct {
	rootShelfScope scopes.RootShelfScopeInterface
}

func NewRootShelfRepository(rootShelfScope scopes.RootShelfScopeInterface) RootShelfRepositoryInterface {
	return &RootShelfRepository{
		rootShelfScope: rootShelfScope,
	}
}

func (r *RootShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.RootShelf{}).
		Select("1").
		Scopes(r.rootShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *RootShelfRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.RootShelf{}).
		Select("DISTINCT \"RootShelfTable\".id").
		Scopes(r.rootShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *RootShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RootShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.RootShelf, enums.AccessControlPermission, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var rootShelf schemas.RootShelf
	result := parsedOptions.DB.
		Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.rootShelfScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&rootShelf)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: rootShelf.Id == uuid.Nil, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, "", exception
	}

	var permission enums.AccessControlPermission
	result = parsedOptions.DB.
		Model(&schemas.UsersToShelves{}).
		Select("permission").
		Where(
			"root_shelf_id = ? AND user_id = ? AND permission IN ?",
			rootShelf.Id,
			userId,
			allowedPermissions,
		).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&permission)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: permission == "", Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, "", exception
	}

	return &rootShelf, permission, nil
}

func (r *RootShelfRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RootShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.RootShelf, []enums.AccessControlPermission, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var rootShelves []schemas.RootShelf
	result := parsedOptions.DB.
		Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.rootShelfScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&rootShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(rootShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, nil, exception
	}

	var usersToShelves []schemas.UsersToShelves
	result = parsedOptions.DB.
		Model(&schemas.UsersToShelves{}).
		Select("root_shelf_id, permission").
		Where(
			"root_shelf_id IN ? AND user_id = ? AND permission IN ?",
			ids,
			userId,
			allowedPermissions,
		).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&usersToShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(usersToShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, nil, exception
	}

	permissionByRootShelfId := make(map[uuid.UUID]enums.AccessControlPermission, len(usersToShelves))
	for _, usersToShelf := range usersToShelves {
		permissionByRootShelfId[usersToShelf.RootShelfId] = usersToShelf.Permission
	}

	permissions := make([]enums.AccessControlPermission, len(rootShelves))
	for index, rootShelf := range rootShelves {
		permission, exist := permissionByRootShelfId[rootShelf.Id]
		if !exist {
			return nil, nil, exceptions.Shelf.NotFound()
		}
		permissions[index] = permission
	}

	return rootShelves, permissions, nil
}

func (r *RootShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RootShelfRelation,
	opts ...options.RepositoryOptions,
) (*schemas.RootShelf, enums.AccessControlPermission, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(
		id,
		userId,
		preloads,
		allowedPermissions,
		opts...,
	)
}

func (r *RootShelfRepository) CreateOneByOwnerId(
	ownerId uuid.UUID,
	input inputs.CreateRootShelfInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	var newRootShelf schemas.RootShelf
	newRootShelf.OwnerId = ownerId
	if err := copier.Copy(&newRootShelf, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Shelf.FailedToCreate().WithOrigin(err)
	}
	if newRootShelf.Id == uuid.Nil {
		newRootShelf.Id = uuid.New()
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newRootShelf)
	if err := result.Error; err != nil {
		parsedOptions.DB.Rollback()
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"shelf_idx_owner_id_name\" (SQLSTATE 23505)":
			return nil, exceptions.Shelf.DuplicateName(input.Name)
		default:
			return nil, exceptions.Shelf.FailedToCreate().WithOrigin(err)
		}
	}

	// create the users to shelves relation with the permission of admin
	newUsersToShelves := schemas.UsersToShelves{
		UserId:      ownerId,
		RootShelfId: newRootShelf.Id,
		Permission:  enums.AccessControlPermission_Owner,
	}
	result = parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Create(&newUsersToShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newRootShelf.Id, nil
}

func (r *RootShelfRepository) CreateManyByOwnerId(
	ownerId uuid.UUID,
	input []inputs.CreateRootShelfInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	var newRootShelves []schemas.RootShelf
	for _, in := range input {
		var newRootShelf schemas.RootShelf
		newRootShelf.OwnerId = ownerId
		if err := copier.Copy(&newRootShelf, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.InvalidDto().WithOrigin(err)
		}
		if newRootShelf.Id == uuid.Nil {
			newRootShelf.Id = uuid.New()
		}
		newRootShelves = append(newRootShelves, newRootShelf)
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		CreateInBatches(&newRootShelves, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newRootShelfIds := make([]uuid.UUID, len(newRootShelves))
	newUsersToShelves := make([]schemas.UsersToShelves, len(newRootShelves))
	for index, newRootShelf := range newRootShelves {
		newRootShelfIds[index] = newRootShelf.Id
		newUsersToShelves[index] = schemas.UsersToShelves{
			UserId:      ownerId,
			RootShelfId: newRootShelf.Id,
			Permission:  enums.AccessControlPermission_Owner,
		}
	}
	result = parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		CreateInBatches(&newUsersToShelves, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newRootShelfIds, nil
}

func (r *RootShelfRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRootShelfInput,
	opts ...options.RepositoryOptions,
) (*schemas.RootShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingRootShelf, _, exception := r.CheckPermissionAndGetOneById(
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRootShelf)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values, input.SetNull, *existingRootShelf,
		).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *RootShelfRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateRootShelfInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction { // only start the transaction when the permission check is required
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	isRootShelfValid := make(map[uuid.UUID]bool)
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

		validRootShelves, _, exception := r.CheckPermissionsAndGetManyByIds(
			ids,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Shelf.NoPermission("update these root shelves")
		}

		for _, validRootShelf := range validRootShelves {
			isRootShelfValid[validRootShelf.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isRootShelfValid[in.Id] {
			continue
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::string, ?::integer, ?::integer, ?::timestamptz)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Name,
			in.PartialUpdateInput.Values.SubShelfCount,
			in.PartialUpdateInput.Values.ItemCount,
			in.PartialUpdateInput.Values.LastAnalyzedAt,
		)
	}

	sql := fmt.Sprintf(`
		UPDATE "RootShelfTable" AS r
		SET
			name = COALESCE(v.name::string, r.name),
			sub_shelf_count = COALESCE(v.sub_shelf_count::integer, r.sub_shelf_count),
			item_count = COALESCE(v.item_count:integer, r.item_count),
			last_analyzed_at = COALESCE(v.last_analyzed_at, r.last_analyzed_at),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, name, sub_shelf_count, item_count, last_analyzed_at)
		WHERE r.id = v.id::uuid AND r.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *RootShelfRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.RootShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredRootShelf schemas.RootShelf
	result := parsedOptions.DB.Model(&restoredRootShelf).
		Scopes(r.rootShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where("\"RootShelfTable\".id = ?", id).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredRootShelf.Id == uuid.Nil, Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredRootShelf, nil
}

func (r *RootShelfRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.RootShelf, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredRootShelves []schemas.RootShelf
	result := parsedOptions.DB.Model(&restoredRootShelves).
		Scopes(r.rootShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where("\"RootShelfTable\".id IN ?", ids).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredRootShelves) != len(ids), Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredRootShelves, nil
}

func (r *RootShelfRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RootShelfTable\".id = ?", id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RootShelfRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RootShelfTable\".id IN ?", ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RootShelfRepository) SoftDeleteManyByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("owner_id = ?", userId).
		Delete(&schemas.RootShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RootShelfTable\".id = ?", id).
		Delete(&schemas.RootShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RootShelfTable\".id IN ?", ids).
		Delete(&schemas.RootShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteManyByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.RootShelf{}).
		Scopes(r.rootShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("owner_id = ?", userId).
		Delete(&schemas.RootShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
