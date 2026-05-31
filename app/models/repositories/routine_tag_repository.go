package repositories

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	scopes "notezy-backend/app/models/scopes"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	array "notezy-backend/shared/lib/array"
	types "notezy-backend/shared/types"
)

type RoutineTagRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RoutineTag, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateRoutineTagInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyByUserId(userId uuid.UUID, input []inputs.BulkCreateRoutineTagInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRoutineTagInput, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateRoutineTagInput, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type RoutineTagRepository struct {
	routineTagScope scopes.RoutineTagScopeInterface
}

func NewRoutineTagRepository(routineTagScope scopes.RoutineTagScopeInterface) RoutineTagRepositoryInterface {
	return &RoutineTagRepository{
		routineTagScope: routineTagScope,
	}
}

func (r *RoutineTagRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Select("1").
		Scopes(r.routineTagScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *RoutineTagRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Select("DISTINCT \"RoutineTagTable\".id").
		Scopes(r.routineTagScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *RoutineTagRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTagRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTag, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTag schemas.RoutineTag
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Scopes(r.routineTagScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineTagScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&routineTag)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.NotFound().WithOrigin(result.Error)},
		{First: routineTag.Id == uuid.Nil, Second: exceptions.RoutineTag.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &routineTag, nil
}

func (r *RoutineTagRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTagRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTag, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTags []schemas.RoutineTag
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Scopes(r.routineTagScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineTagScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&routineTags)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.NotFound().WithOrigin(result.Error)},
		{First: len(routineTags) == 0, Second: exceptions.RoutineTag.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return routineTags, nil
}

func (r *RoutineTagRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTagRelation,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTag, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(id, userId, preloads, allowedPermissions, opts...)
}

func (r *RoutineTagRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateRoutineTagInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	newRoutineTag := schemas.RoutineTag{
		Id:    uuid.New(),
		Color: "#FFFFFF",
	}
	if err := copier.Copy(&newRoutineTag, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTag.InvalidInput().WithOrigin(err)
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Create(&newRoutineTag)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newUsersToRoutineTag := schemas.UsersToRoutineTags{
		UserId:     userId,
		TagId:      newRoutineTag.Id,
		Permission: enums.AccessControlPermission_Owner,
	}
	result = parsedOptions.DB.
		Model(&schemas.UsersToRoutineTags{}).
		Create(&newUsersToRoutineTag)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTag.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newRoutineTag.Id, nil
}

func (r *RoutineTagRepository) BulkCreateManyByUserId(
	userId uuid.UUID,
	input []inputs.BulkCreateRoutineTagInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.RoutineTag.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	newRoutineTags := make([]schemas.RoutineTag, 0, len(input))
	for _, in := range input {
		newRoutineTag := schemas.RoutineTag{
			Id:    uuid.New(),
			Color: "#FFFFFF",
		}
		if err := copier.Copy(&newRoutineTag, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTag.InvalidInput().WithOrigin(err)
		}
		if newRoutineTag.Id == uuid.Nil {
			newRoutineTag.Id = uuid.New()
		}
		if newRoutineTag.Color == "" {
			newRoutineTag.Color = "#FFFFFF"
		}
		newRoutineTags = append(newRoutineTags, newRoutineTag)
	}

	if len(newRoutineTags) == 0 {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTag.NoChanges()
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		CreateInBatches(&newRoutineTags, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newUsersToRoutineTags := make([]schemas.UsersToRoutineTags, len(newRoutineTags))
	newRoutineTagIds := make([]uuid.UUID, len(newRoutineTags))
	for index, newRoutineTag := range newRoutineTags {
		newRoutineTagIds[index] = newRoutineTag.Id
		newUsersToRoutineTags[index] = schemas.UsersToRoutineTags{
			UserId:     userId,
			TagId:      newRoutineTag.Id,
			Permission: enums.AccessControlPermission_Owner,
		}
	}
	result = parsedOptions.DB.
		Model(&schemas.UsersToRoutineTags{}).
		CreateInBatches(&newUsersToRoutineTags, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTag.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newRoutineTagIds, nil
}

func (r *RoutineTagRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRoutineTagInput,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTag, *exceptions.Exception) {
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
	existingRoutineTag, exception := r.CheckPermissionAndGetOneById(id, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRoutineTag)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingRoutineTag).WithOrigin(err)
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Where("\"RoutineTagTable\".id = ?", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTag.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *RoutineTagRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateRoutineTagInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(input) == 0 {
		return exceptions.RoutineTag.NoChanges()
	}

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
	ids := make([]uuid.UUID, len(input))
	for index, in := range input {
		ids[index] = in.Id
	}
	validRoutineTags, exception := r.CheckPermissionsAndGetManyByIds(ids, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exceptions.RoutineTag.NoPermission("update these routine tags")
	}

	routineTagById := make(map[uuid.UUID]schemas.RoutineTag, len(validRoutineTags))
	for _, validRoutineTag := range validRoutineTags {
		routineTagById[validRoutineTag.Id] = validRoutineTag
	}

	for _, in := range input {
		existingRoutineTag, exist := routineTagById[in.Id]
		if !exist {
			continue
		}
		updates, err := util.PartialUpdatePreprocess(in.PartialUpdateInput.Values, in.PartialUpdateInput.SetNull, existingRoutineTag)
		if err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Util.FailedToPreprocessPartialUpdate(
				in.PartialUpdateInput.Values,
				in.PartialUpdateInput.SetNull,
				existingRoutineTag,
			).WithOrigin(err)
		}
		result := parsedOptions.DB.
			Model(&schemas.RoutineTag{}).
			Where("\"RoutineTagTable\".id = ?", in.Id).
			Select("*").
			Updates(&updates)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToUpdate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return exception
		}
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.RoutineTag.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *RoutineTagRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Scopes(r.routineTagScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Where("\"RoutineTagTable\".id = ?", id).
		Delete(&schemas.RoutineTag{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineTagRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.RoutineTag.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Scopes(r.routineTagScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Where("\"RoutineTagTable\".id IN ?", ids).
		Delete(&schemas.RoutineTag{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
