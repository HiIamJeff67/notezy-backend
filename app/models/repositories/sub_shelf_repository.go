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

type SubShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	GetAllByRootShelfId(rootShelfId uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	CreateOneByRootShelfId(rootShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateSubShelfInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyByRootShelfIds(userId uuid.UUID, input []inputs.BulkCreateSubShelfInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateSubShelfInput, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateSubShelfInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type SubShelfRepository struct {
	subShelfScope scopes.SubShelfScopeInterface
}

func NewSubShelfRepository(subShelfScope scopes.SubShelfScopeInterface) SubShelfRepositoryInterface {
	return &SubShelfRepository{
		subShelfScope: subShelfScope,
	}
}

func (r *SubShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.SubShelf{}).
		Select("1").
		Scopes(r.subShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *SubShelfRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.SubShelf{}).
		Select(`DISTINCT "SubShelfTable".id`).
		Scopes(r.subShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *SubShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subShelf := schemas.SubShelf{}
	result := parsedOptions.DB.
		Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.subShelfScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&subShelf)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: subShelf.Id == uuid.Nil, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &subShelf, nil
}

func (r *SubShelfRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subShelves := []schemas.SubShelf{}
	result := parsedOptions.DB.
		Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.subShelfScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&subShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(subShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return subShelves, nil
}

func (r *SubShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
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

func (r *SubShelfRepository) GetAllByRootShelfId(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelves := []schemas.SubShelf{}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where(`root_shelf_id = "SubShelfTable".root_shelf_id AND user_id = ? AND permission IN ?`,
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("root_shelf_id = ? AND EXISTS (?)", rootShelfId, subQuery)
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&subShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(subShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return subShelves, nil
}

func (r *SubShelfRepository) CreateOneByRootShelfId(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateSubShelfInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB), options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	var newSubShelf schemas.SubShelf
	if input.PrevSubShelfId != nil {
		prevSubShelf, exception := r.CheckPermissionAndGetOneById(
			*input.PrevSubShelfId,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception = exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
			{First: prevSubShelf.RootShelfId != rootShelfId, Second: exceptions.Shelf.InvalidDto("the given prev sub shelf is not one of the children of the given root shelf")},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return nil, exception
		}
		prevSubShelf.Path = append(prevSubShelf.Path, prevSubShelf.Id)
		newSubShelf.Path = prevSubShelf.Path
	} else {
		rootShelfRepository := NewRootShelfRepository(scopes.NewRootShelfScope())

		if !rootShelfRepository.HasPermission(
			rootShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.NoPermission("create sub shelf by the given root shelf")
		}
	}

	if err := copier.Copy(&newSubShelf, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Shelf.InvalidInput().WithOrigin(err)
	}
	if newSubShelf.Id == uuid.Nil {
		newSubShelf.Id = uuid.New()
	}
	newSubShelf.RootShelfId = rootShelfId

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Create(&newSubShelf)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: newSubShelf.Id == uuid.Nil, Second: exceptions.Shelf.FailedToCreate()},
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

	return &newSubShelf.Id, nil
}

func (r *SubShelfRepository) BulkCreateManyByRootShelfIds(
	userId uuid.UUID,
	input []inputs.BulkCreateSubShelfInput,
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
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB), options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	isPrevSubShelfExist := make(map[uuid.UUID]bool)
	isRootShelfExist := make(map[uuid.UUID]bool)
	prevSubShelfIds := make([]uuid.UUID, len(input))
	rootShelfIds := make([]uuid.UUID, len(input))
	for index, in := range input {
		if in.PrevSubShelfId != nil {
			if isPrevSubShelfExist[*in.PrevSubShelfId] {
				prevSubShelfIds[index] = *in.PrevSubShelfId
			}
			isPrevSubShelfExist[*in.PrevSubShelfId] = true
		}
		if isRootShelfExist[in.RootShelfId] {
			rootShelfIds[index] = in.RootShelfId
		}
		isRootShelfExist[in.RootShelfId] = true
	}

	validPrevSubShelves, exception := r.CheckPermissionsAndGetManyByIds(
		prevSubShelfIds,
		userId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	isPrevSubShelfValid := make(map[uuid.UUID]*uuid.UUID)
	for _, validPrevSubShelf := range validPrevSubShelves {
		isPrevSubShelfValid[validPrevSubShelf.Id] = &validPrevSubShelf.RootShelfId
	}

	rootShelfRepository := NewRootShelfRepository(scopes.NewRootShelfScope())

	validRootShelves, _, exception := rootShelfRepository.CheckPermissionsAndGetManyByIds(
		rootShelfIds,
		userId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	isRootShelfValid := make(map[uuid.UUID]bool)
	for _, validRootShelf := range validRootShelves {
		isRootShelfValid[validRootShelf.Id] = true
	}

	var newSubShelves []schemas.SubShelf
	for _, in := range input {
		if !isRootShelfValid[in.RootShelfId] ||
			(in.PrevSubShelfId != nil && (isPrevSubShelfValid[*in.PrevSubShelfId] != &in.RootShelfId)) {
			continue
		}
		var newSubShelf schemas.SubShelf
		if err := copier.Copy(&newSubShelf, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.InvalidInput().WithOrigin(err)
		}
		if newSubShelf.Id == uuid.Nil {
			newSubShelf.Id = uuid.New()
		}
		newSubShelf.RootShelfId = in.RootShelfId
		newSubShelves = append(newSubShelves, newSubShelf)
	}

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		CreateInBatches(&newSubShelves, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newSubShelfIds := make([]uuid.UUID, len(newSubShelves))
	for index, newSubShelf := range newSubShelves {
		newSubShelfIds[index] = newSubShelf.Id
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newSubShelfIds, nil
}

func (r *SubShelfRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateSubShelfInput,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if !parsedOptions.IsTransactionStarted {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB), options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingSubShelf, exception := r.CheckPermissionAndGetOneById(
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingSubShelf)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingSubShelf).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
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

func (r *SubShelfRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateSubShelfInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB), options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	isSubShelfValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		subShelfIds := make([]uuid.UUID, len(input))
		for index, in := range input {
			subShelfIds[index] = in.Id
		}

		validSubShelves, exception := r.CheckPermissionsAndGetManyByIds(
			subShelfIds,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			parsedOptions.DB.Rollback()
			return exception
		}

		for _, validSubShelf := range validSubShelves {
			isSubShelfValid[validSubShelf.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isSubShelfValid[in.Id] {
			continue
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::string)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Name,
		)
	}

	sql := fmt.Sprintf(`
		UPDATE "SubShelfTable" AS s
		SET
			name = COALESCE(v.name::string, s.name)
		FROM (VALUES %s) AS v(id, name)
		WHERE s.id = v.id::uuid AND s.deleted_at IS NULL
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

func (r *SubShelfRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredSubShelf schemas.SubShelf
	result := parsedOptions.DB.Model(&restoredSubShelf).
		Scopes(r.subShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where(`"SubShelfTable".id = ?`, id).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredSubShelf.Id == uuid.Nil, Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredSubShelf, nil
}

func (r *SubShelfRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredSubShelves []schemas.SubShelf
	result := parsedOptions.DB.Model(&restoredSubShelves).
		Scopes(r.subShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where(`"SubShelfTable".id IN ?`, ids).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredSubShelves) == 0, Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredSubShelves, nil
}

func (r *SubShelfRepository) SoftDeleteOneById(
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

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"SubShelfTable".id = ?`, id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) SoftDeleteManyByIds(
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

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"SubShelfTable".id IN ?`, ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteOneById(
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

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"SubShelfTable".id = ?`, id).
		Delete(&schemas.SubShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteManyByIds(
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

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Scopes(r.subShelfScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.subShelfScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"SubShelfTable".id IN ?`, ids).
		Delete(&schemas.SubShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
