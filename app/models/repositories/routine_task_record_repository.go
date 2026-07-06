package repositories

import (
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineTaskRecordRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRecordRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RoutineTaskRecord, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRecordRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RoutineTaskRecord, *exceptions.Exception)
	GetAllByRoutineTaskId(routineTaskId uuid.UUID, userId uuid.UUID, limit int, preloads []schemas.RoutineTaskRecordRelation, opts ...options.RepositoryOptions) ([]schemas.RoutineTaskRecord, *exceptions.Exception)
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type RoutineTaskRecordRepository struct {
	routineTaskRecordScope scopes.RoutineTaskRecordScopeInterface
}

func NewRoutineTaskRecordRepository(
	routineTaskRecordScope scopes.RoutineTaskRecordScopeInterface,
) RoutineTaskRecordRepositoryInterface {
	return &RoutineTaskRecordRepository{
		routineTaskRecordScope: routineTaskRecordScope,
	}
}

func (r *RoutineTaskRecordRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Select("1").
		Scopes(r.routineTaskRecordScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Limit(1).
		Scan(&marker)
	if result.Error != nil {
		return false
	}

	return marker == 1
}

func (r *RoutineTaskRecordRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Select(`DISTINCT "RoutineTaskRecordTable".id`).
		Scopes(r.routineTaskRecordScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&permittedIds)
	if result.Error != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *RoutineTaskRecordRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRecordRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTaskRecord, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTaskRecord schemas.RoutineTaskRecord
	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Scopes(r.routineTaskRecordScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineTaskRecordScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&routineTaskRecord)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.NotFound().WithOrigin(result.Error)},
		{First: routineTaskRecord.Id == uuid.Nil, Second: exceptions.RoutineTask.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &routineTaskRecord, nil
}

func (r *RoutineTaskRecordRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRecordRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTaskRecord, *exceptions.Exception) {
	if len(ids) == 0 {
		return []schemas.RoutineTaskRecord{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTaskRecords []schemas.RoutineTaskRecord
	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Scopes(r.routineTaskRecordScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineTaskRecordScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&routineTaskRecords)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.NotFound().WithOrigin(result.Error)},
		{First: len(routineTaskRecords) == 0, Second: exceptions.RoutineTask.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return routineTaskRecords, nil
}

func (r *RoutineTaskRecordRepository) GetAllByRoutineTaskId(
	routineTaskId uuid.UUID,
	userId uuid.UUID,
	limit int,
	preloads []schemas.RoutineTaskRecordRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTaskRecord, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	if limit <= 0 {
		limit = 100
	}

	var routineTaskRecords []schemas.RoutineTaskRecord
	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Select(`"RoutineTaskRecordTable".*`).
		Joins(`INNER JOIN "RoutineTaskTable" routine_task ON routine_task.id = "RoutineTaskRecordTable".routine_task_id`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = routine_task.routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where(`"RoutineTaskRecordTable".routine_task_id = ?`, routineTaskId).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(r.routineTaskRecordScope.IncludePreloads(preloads)).
		Order(`"RoutineTaskRecordTable".created_at DESC`).
		Limit(limit).
		Find(&routineTaskRecords)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	return routineTaskRecords, nil
}

func (r *RoutineTaskRecordRepository) HardDeleteOneById(
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
		Model(&schemas.RoutineTaskRecord{}).
		Scopes(r.routineTaskRecordScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Where(`"RoutineTaskRecordTable".id = ?`, id).
		Delete(&schemas.RoutineTaskRecord{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineTaskRecordRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.RoutineTask.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTaskRecord{}).
		Scopes(r.routineTaskRecordScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Where(`"RoutineTaskRecordTable".id IN ?`, ids).
		Delete(&schemas.RoutineTaskRecord{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
