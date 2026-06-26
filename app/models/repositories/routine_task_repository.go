package repositories

import (
	"fmt"
	"strings"
	"time"

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

type RoutineTaskRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RoutineTask, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RoutineTask, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, opts ...options.RepositoryOptions) (*schemas.RoutineTask, *exceptions.Exception)
	GetAllByStationIds(stationIds []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, opts ...options.RepositoryOptions) ([]schemas.RoutineTask, *exceptions.Exception)
	CreateOneByStationId(stationId uuid.UUID, userId uuid.UUID, input inputs.CreateRoutineTaskInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyByStationIds(userId uuid.UUID, input []inputs.BulkCreateRoutineTaskInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRoutineTaskInput, opts ...options.RepositoryOptions) (*schemas.RoutineTask, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateRoutineTaskInput, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type RoutineTaskRepository struct {
	routineTaskScope scopes.RoutineTaskScopeInterface
}

func NewRoutineTaskRepository(routineTaskScope scopes.RoutineTaskScopeInterface) RoutineTaskRepositoryInterface {
	return &RoutineTaskRepository{
		routineTaskScope: routineTaskScope,
	}
}

func (r *RoutineTaskRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Select("1").
		Scopes(r.routineTaskScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *RoutineTaskRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Select(`DISTINCT "RoutineTaskTable".id`).
		Scopes(r.routineTaskScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *RoutineTaskRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTask, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTask schemas.RoutineTask
	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Scopes(r.routineTaskScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineTaskScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&routineTask)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.NotFound().WithOrigin(result.Error)},
		{First: routineTask.Id == uuid.Nil, Second: exceptions.RoutineTask.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &routineTask, nil
}

func (r *RoutineTaskRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTask, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTasks []schemas.RoutineTask
	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Scopes(r.routineTaskScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineTaskScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&routineTasks)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.NotFound().WithOrigin(result.Error)},
		{First: len(routineTasks) == 0, Second: exceptions.RoutineTask.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return routineTasks, nil
}

func (r *RoutineTaskRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRelation,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTask, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(id, userId, preloads, allowedPermissions, opts...)
}

func (r *RoutineTaskRepository) GetAllByStationIds(
	stationIds []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineTaskRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.RoutineTask, *exceptions.Exception) {
	if len(stationIds) == 0 {
		return []schemas.RoutineTask{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var routineTasks []schemas.RoutineTask
	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Select(`"RoutineTaskTable".*`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "RoutineTaskTable".station_id`).
		Joins(`INNER JOIN "StationTable" station ON station.id = "RoutineTaskTable".station_id AND station.deleted_at IS NULL`).
		Where(`"RoutineTaskTable".station_id IN ?`, stationIds).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(r.routineTaskScope.IncludePreloads(preloads)).
		Order(`"RoutineTaskTable".scheduled_at ASC`).
		Order(`"RoutineTaskTable".priority DESC`).
		Order(`"RoutineTaskTable".id ASC`).
		Find(&routineTasks)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	return routineTasks, nil
}

func (r *RoutineTaskRepository) CreateOneByStationId(
	stationId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateRoutineTaskInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
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
	stationRepository := NewStationRepository(scopes.NewStationScope())
	if !stationRepository.HasPermission(
		stationId,
		userId,
		allowedPermissions,
		opts...,
	) {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTask.NoPermission("create a routine task under this station")
	}

	newRoutineTask := schemas.RoutineTask{}
	if err := copier.Copy(&newRoutineTask, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTask.InvalidInput().WithOrigin(err)
	}
	newRoutineTask.StationId = stationId
	newRoutineTask.ScheduledAt = input.ScheduledAt.Truncate(time.Minute)

	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Create(&newRoutineTask)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newRoutineTask.Id, nil
}

func (r *RoutineTaskRepository) BulkCreateManyByStationIds(
	userId uuid.UUID,
	input []inputs.BulkCreateRoutineTaskInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.RoutineTask.NoChanges()
	}

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
	stationIds := make([]uuid.UUID, len(input))
	for index, in := range input {
		stationIds[index] = in.StationId
	}

	stationRepository := NewStationRepository(scopes.NewStationScope())
	validStations, _, exception := stationRepository.CheckPermissionsAndGetManyByIds(stationIds, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	isStationValid := make(map[uuid.UUID]bool, len(validStations))
	for _, validStation := range validStations {
		isStationValid[validStation.Id] = true
	}

	newRoutineTasks := make([]schemas.RoutineTask, 0, len(input))
	for _, in := range input {
		if !isStationValid[in.StationId] {
			continue
		}
		if in.ScheduledAt.IsZero() {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("scheduledAt is required"))
		}
		newRoutineTask := schemas.RoutineTask{
			Id:        uuid.New(),
			StationId: in.StationId,
		}
		if err := copier.Copy(&newRoutineTask, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.InvalidInput().WithOrigin(err)
		}
		newRoutineTask.StationId = in.StationId
		newRoutineTask.ScheduledAt = in.ScheduledAt.Truncate(time.Minute)
		newRoutineTasks = append(newRoutineTasks, newRoutineTask)
	}

	if len(newRoutineTasks) == 0 {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTask.NoChanges()
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		CreateInBatches(&newRoutineTasks, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	newRoutineTaskIds := make([]uuid.UUID, len(newRoutineTasks))
	for index, newRoutineTask := range newRoutineTasks {
		newRoutineTaskIds[index] = newRoutineTask.Id
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newRoutineTaskIds, nil
}

func (r *RoutineTaskRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRoutineTaskInput,
	opts ...options.RepositoryOptions,
) (*schemas.RoutineTask, *exceptions.Exception) {
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
	existingRoutineTask, exception := r.CheckPermissionAndGetOneById(id, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	if input.Values.StationId != nil && (input.SetNull == nil || !(*input.SetNull)["StationId"]) {
		stationRepository := NewStationRepository(scopes.NewStationScope())
		if !stationRepository.HasPermission(*input.Values.StationId, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.NoPermission("move a routine task to this station")
		}
	}
	if input.Values.ScheduledAt != nil {
		truncatedScheduledAt := input.Values.ScheduledAt.Truncate(time.Minute)
		input.Values.ScheduledAt = &truncatedScheduledAt
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRoutineTask)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingRoutineTask).WithOrigin(err)
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Where(`"RoutineTaskTable".id = ?`, id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *RoutineTaskRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateRoutineTaskInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(input) == 0 {
		return exceptions.RoutineTask.NoChanges()
	}

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
	ids := make([]uuid.UUID, len(input))
	for index, in := range input {
		ids[index] = in.Id
	}
	validRoutineTasks, exception := r.CheckPermissionsAndGetManyByIds(ids, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exceptions.RoutineTask.NoPermission("update these routine tasks")
	}

	isRoutineTaskValid := make(map[uuid.UUID]bool, len(validRoutineTasks))
	for _, validRoutineTask := range validRoutineTasks {
		isRoutineTaskValid[validRoutineTask.Id] = true
	}

	targetStationIdSet := make(map[uuid.UUID]bool)
	for _, in := range input {
		if in.PartialUpdateInput.Values.StationId == nil ||
			(in.PartialUpdateInput.SetNull != nil && (*in.PartialUpdateInput.SetNull)["StationId"]) {
			continue
		}
		targetStationIdSet[*in.PartialUpdateInput.Values.StationId] = true
	}
	if len(targetStationIdSet) > 0 {
		targetStationIds := make([]uuid.UUID, 0, len(targetStationIdSet))
		for targetStationId := range targetStationIdSet {
			targetStationIds = append(targetStationIds, targetStationId)
		}
		stationRepository := NewStationRepository(scopes.NewStationScope())
		if !stationRepository.HavePermissions(targetStationIds, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return exceptions.RoutineTask.NoPermission("move these routine tasks to the given stations")
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !isRoutineTaskValid[in.Id] {
			continue
		}

		setPeriodNull := false
		if in.PartialUpdateInput.SetNull != nil {
			for field, setNull := range *in.PartialUpdateInput.SetNull {
				if setNull && strings.ToLower(strings.ReplaceAll(field, "_", "")) == "period" {
					setPeriodNull = true
					break
				}
			}
		}

		scheduledAt := in.PartialUpdateInput.Values.ScheduledAt
		if scheduledAt != nil {
			truncatedScheduledAt := scheduledAt.Truncate(time.Minute)
			scheduledAt = &truncatedScheduledAt
		}

		valuePlaceholders = append(valuePlaceholders, `(?::uuid, ?::uuid, ?::text, ?::"RoutineTaskPurpose", ?::jsonb, ?::integer, ?::integer, ?::"RoutinePeriod", ?::timestamptz, ?::boolean)`)
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.StationId,
			in.PartialUpdateInput.Values.Title,
			in.PartialUpdateInput.Values.Purpose,
			in.PartialUpdateInput.Values.Payload,
			in.PartialUpdateInput.Values.Priority,
			in.PartialUpdateInput.Values.MaxAttempts,
			in.PartialUpdateInput.Values.Period,
			scheduledAt,
			setPeriodNull,
		)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.RoutineTask.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "RoutineTaskTable" AS rt
		SET
			station_id = COALESCE(v.station_id::uuid, rt.station_id),
			title = COALESCE(v.title::text, rt.title),
			purpose = COALESCE(v.purpose::"RoutineTaskPurpose", rt.purpose),
			payload = COALESCE(v.payload::jsonb, rt.payload),
			priority = COALESCE(v.priority::integer, rt.priority),
			max_attempts = COALESCE(v.max_attempts::integer, rt.max_attempts),
			period = CASE
				WHEN v.set_period_null::boolean THEN NULL
				ELSE COALESCE(v.period::"RoutinePeriod", rt.period)
			END,
			scheduled_at = COALESCE(v.scheduled_at::timestamptz, rt.scheduled_at),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, station_id, title, purpose, payload, priority, max_attempts, period, scheduled_at, set_period_null)
		WHERE rt.id = v.id::uuid
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *RoutineTaskRepository) HardDeleteOneById(
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
		Model(&schemas.RoutineTask{}).
		Scopes(r.routineTaskScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Where(`"RoutineTaskTable".id = ?`, id).
		Delete(&schemas.RoutineTask{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineTaskRepository) HardDeleteManyByIds(
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
		Model(&schemas.RoutineTask{}).
		Scopes(r.routineTaskScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Where(`"RoutineTaskTable".id IN ?`, ids).
		Delete(&schemas.RoutineTask{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
