package repositories

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

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

type RoutineTagRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RoutineTag, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTagRelation, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	GetAllByUserId(userId uuid.UUID, preloads []schemas.RoutineTagRelation, opts ...options.RepositoryOptions) ([]schemas.RoutineTag, *exceptions.Exception)
	CreateOne(userId uuid.UUID, input inputs.CreateRoutineTagInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	CreateMany(userId uuid.UUID, input []inputs.CreateRoutineTagInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRoutineTagInput, opts ...options.RepositoryOptions) (*schemas.RoutineTag, *exceptions.Exception)
	UpdateManyByIds(userId uuid.UUID, input []inputs.UpdateRoutineTagByIdInput, opts ...options.RepositoryOptions) *exceptions.Exception
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
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
		Select(`DISTINCT "RoutineTagTable".id`).
		Scopes(r.routineTagScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
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

func (r *RoutineTagRepository) GetAllByUserId(
	userId uuid.UUID,
	preloads []schemas.RoutineTagRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTag, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	var routineTags []schemas.RoutineTag
	result := parsedOptions.DB.
		Model(&schemas.RoutineTag{}).
		Select(`"RoutineTagTable".*`).
		Joins(`INNER JOIN "UsersToRoutineTagsTable" utrt ON utrt.tag_id = "RoutineTagTable".id`).
		Where("utrt.user_id = ? AND utrt.permission IN ?", userId, allowedPermissions).
		Scopes(r.routineTagScope.IncludePreloads(preloads)).
		Order(`"RoutineTagTable".created_at ASC`).
		Order(`"RoutineTagTable".id ASC`).
		Find(&routineTags)
	if result.Error != nil {
		return nil, exceptions.RoutineTag.NotFound().WithOrigin(result.Error)
	}

	return routineTags, nil
}

func (r *RoutineTagRepository) CreateOne(
	userId uuid.UUID,
	input inputs.CreateRoutineTagInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
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

func (r *RoutineTagRepository) CreateMany(
	userId uuid.UUID,
	input []inputs.CreateRoutineTagInput,
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
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
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
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
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
		Where(`"RoutineTagTable".id = ?`, id).
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

func (r *RoutineTagRepository) UpdateManyByIds(
	userId uuid.UUID,
	input []inputs.UpdateRoutineTagByIdInput,
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
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
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

	isRoutineTagValid := make(map[uuid.UUID]bool, len(validRoutineTags))
	for _, validRoutineTag := range validRoutineTags {
		isRoutineTagValid[validRoutineTag.Id] = true
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !isRoutineTagValid[in.Id] {
			continue
		}

		setIconNull := util.CheckSetNull(in.PartialUpdateInput.SetNull, "Icon")

		valuePlaceholders = append(valuePlaceholders, `(?::uuid, ?::text, ?::text, ?::"SupportedIcon", ?::boolean)`)
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Name,
			in.PartialUpdateInput.Values.Color,
			in.PartialUpdateInput.Values.Icon,
			setIconNull,
		)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.RoutineTag.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "RoutineTagTable" AS rt
		SET
			name = COALESCE(v.name::text, rt.name),
			color = COALESCE(v.color::text, rt.color),
			icon = CASE
				WHEN v.set_icon_null::boolean THEN NULL
				ELSE COALESCE(v.icon::"SupportedIcon", rt.icon)
			END,
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, name, color, icon, set_icon_null)
		WHERE rt.id = v.id::uuid
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
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
		Where(`"RoutineTagTable".id = ?`, id).
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
		Where(`"RoutineTagTable".id IN ?`, ids).
		Delete(&schemas.RoutineTag{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTag.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTag.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
