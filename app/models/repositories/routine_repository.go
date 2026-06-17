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

type RoutineRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Routine, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Routine, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineRelation, opts ...options.RepositoryOptions) (*schemas.Routine, *exceptions.Exception)
	GetAllByTimeRange(from time.Time, to time.Time, stationIds []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineRelation, opts ...options.RepositoryOptions) ([]schemas.Routine, *exceptions.Exception)
	CreateOneByStationId(stationId uuid.UUID, userId uuid.UUID, input inputs.CreateRoutineInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyByStationIds(userId uuid.UUID, input []inputs.BulkCreateRoutineInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRoutineInput, opts ...options.RepositoryOptions) (*schemas.Routine, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateRoutineInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Routine, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Routine, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type RoutineRepository struct {
	routineScope scopes.RoutineScopeInterface
}

func NewRoutineRepository(routineScope scopes.RoutineScopeInterface) RoutineRepositoryInterface {
	return &RoutineRepository{
		routineScope: routineScope,
	}
}

func (r *RoutineRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Select("1").
		Scopes(r.routineScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *RoutineRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Select("DISTINCT \"RoutineTable\".id").
		Scopes(r.routineScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *RoutineRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Routine, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routine schemas.Routine
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.routineScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&routine)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.NotFound().WithOrigin(result.Error)},
		{First: routine.Id == uuid.Nil, Second: exceptions.Routine.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &routine, nil
}

func (r *RoutineRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.Routine, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routines []schemas.Routine
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.routineScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&routines)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.NotFound().WithOrigin(result.Error)},
		{First: len(routines) == 0, Second: exceptions.Routine.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return routines, nil
}

func (r *RoutineRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineRelation,
	opts ...options.RepositoryOptions,
) (*schemas.Routine, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(id, userId, preloads, allowedPermissions, opts...)
}

func (r *RoutineRepository) GetAllByTimeRange(
	from time.Time,
	to time.Time,
	stationIds []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RoutineRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.Routine, *exceptions.Exception) {
	if len(stationIds) == 0 {
		return []schemas.Routine{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var routines []schemas.Routine
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Select("\"RoutineTable\".*").
		Joins("INNER JOIN \"UsersToStationsTable\" uts ON uts.station_id = \"RoutineTable\".station_id").
		Joins("INNER JOIN \"StationTable\" station ON station.id = \"RoutineTable\".station_id AND station.deleted_at IS NULL").
		Where("\"RoutineTable\".station_id IN ?", stationIds).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Where("\"RoutineTable\".scheduled_start_at < ? AND \"RoutineTable\".scheduled_end_at > ?", to, from).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.routineScope.IncludePreloads(preloads)).
		Order("\"RoutineTable\".scheduled_start_at ASC").
		Order("\"RoutineTable\".scheduled_end_at ASC").
		Order("\"RoutineTable\".id ASC").
		Find(&routines)
	if result.Error != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(result.Error)
	}

	return routines, nil
}

func (r *RoutineRepository) CreateOneByStationId(
	stationId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateRoutineInput,
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
	stationRepository := NewStationRepository(scopes.NewStationScope())
	if !stationRepository.HasPermission(stationId, userId, allowedPermissions, opts...) {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Routine.NoPermission("create a routine under this station")
	}

	startAt := time.Now().Truncate(time.Minute)
	newRoutine := schemas.Routine{
		Id:               uuid.New(),
		StationId:        stationId,
		Status:           enums.RoutineStatus_Scheduled,
		ScheduledStartAt: startAt,
		ScheduledEndAt:   startAt.Add(time.Hour),
		Timezone:         "UTC",
	}
	if err := copier.Copy(&newRoutine, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Routine.InvalidInput().WithOrigin(err)
	}
	newRoutine.StationId = stationId

	result := parsedOptions.DB.Model(&schemas.Routine{}).
		Create(&newRoutine)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newRoutine.Id, nil
}

func (r *RoutineRepository) BulkCreateManyByStationIds(
	userId uuid.UUID,
	input []inputs.BulkCreateRoutineInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.Routine.NoChanges()
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

	newRoutines := make([]schemas.Routine, 0, len(input))
	for _, in := range input {
		if !isStationValid[in.StationId] {
			continue
		}
		startAt := time.Now().Truncate(time.Minute)
		newRoutine := schemas.Routine{
			Id:               uuid.New(),
			StationId:        in.StationId,
			Status:           enums.RoutineStatus_Scheduled,
			ScheduledStartAt: startAt,
			ScheduledEndAt:   startAt.Add(time.Hour),
			Timezone:         "UTC",
		}
		if err := copier.Copy(&newRoutine, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Routine.InvalidInput().WithOrigin(err)
		}
		newRoutine.StationId = in.StationId
		newRoutines = append(newRoutines, newRoutine)
	}
	if len(newRoutines) == 0 {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	result := parsedOptions.DB.Model(&schemas.Routine{}).
		CreateInBatches(&newRoutines, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newRoutineIds := make([]uuid.UUID, len(newRoutines))
	for index, newRoutine := range newRoutines {
		newRoutineIds[index] = newRoutine.Id
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newRoutineIds, nil
}

func (r *RoutineRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRoutineInput,
	opts ...options.RepositoryOptions,
) (*schemas.Routine, *exceptions.Exception) {
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
	existingRoutine, exception := r.CheckPermissionAndGetOneById(id, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	if input.Values.StationId != nil && (input.SetNull == nil || !(*input.SetNull)["StationId"]) {
		stationRepository := NewStationRepository(scopes.NewStationScope())
		if !stationRepository.HasPermission(*input.Values.StationId, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Routine.NoPermission("move a routine to this station")
		}
	}
	if input.Values.ScheduledStartAt != nil {
		truncatedScheduledStartAt := input.Values.ScheduledStartAt.Truncate(time.Minute)
		input.Values.ScheduledStartAt = &truncatedScheduledStartAt
	}
	if input.Values.ScheduledEndAt != nil {
		truncatedScheduledEndAt := input.Values.ScheduledEndAt.Truncate(time.Minute)
		input.Values.ScheduledEndAt = &truncatedScheduledEndAt
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRoutine)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingRoutine).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.Routine{}).
		Where("\"RoutineTable\".id = ? AND \"RoutineTable\".deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *RoutineRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateRoutineInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(input) == 0 {
		return exceptions.Routine.NoChanges()
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
	ids := make([]uuid.UUID, len(input))
	for index, in := range input {
		ids[index] = in.Id
	}
	validRoutines, exception := r.CheckPermissionsAndGetManyByIds(ids, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exceptions.Routine.NoPermission("update these routines")
	}
	isRoutineValid := make(map[uuid.UUID]bool, len(validRoutines))
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
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
			return exceptions.Routine.NoPermission("move these routines to the given stations")
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !isRoutineValid[in.Id] {
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

		scheduledStartAt := in.PartialUpdateInput.Values.ScheduledStartAt
		if scheduledStartAt != nil {
			truncatedScheduledStartAt := scheduledStartAt.Truncate(time.Minute)
			scheduledStartAt = &truncatedScheduledStartAt
		}
		scheduledEndAt := in.PartialUpdateInput.Values.ScheduledEndAt
		if scheduledEndAt != nil {
			truncatedScheduledEndAt := scheduledEndAt.Truncate(time.Minute)
			scheduledEndAt = &truncatedScheduledEndAt
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::text, ?::text, ?::\"RoutineStatus\", ?::boolean, ?::timestamptz, ?::timestamptz, ?::\"RoutinePeriod\", ?::text, ?::boolean)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.StationId,
			in.PartialUpdateInput.Values.Title,
			in.PartialUpdateInput.Values.Description,
			in.PartialUpdateInput.Values.Status,
			in.PartialUpdateInput.Values.IsPinned,
			scheduledStartAt,
			scheduledEndAt,
			in.PartialUpdateInput.Values.Period,
			in.PartialUpdateInput.Values.Timezone,
			setPeriodNull,
		)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.Routine.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "RoutineTable" AS r
		SET
			station_id = COALESCE(v.station_id::uuid, r.station_id),
			title = COALESCE(v.title::text, r.title),
			description = COALESCE(v.description::text, r.description),
			status = COALESCE(v.status::"RoutineStatus", r.status),
			is_pinned = COALESCE(v.is_pinned::boolean, r.is_pinned),
			scheduled_start_at = COALESCE(v.scheduled_start_at::timestamptz, r.scheduled_start_at),
			scheduled_end_at = COALESCE(v.scheduled_end_at::timestamptz, r.scheduled_end_at),
			period = CASE
				WHEN v.set_period_null::boolean THEN NULL
				ELSE COALESCE(v.period::"RoutinePeriod", r.period)
			END,
			timezone = COALESCE(v.timezone::text, r.timezone),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, station_id, title, description, status, is_pinned, scheduled_start_at, scheduled_end_at, period, timezone, set_period_null)
		WHERE r.id = v.id::uuid AND r.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *RoutineRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Routine, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredRoutine schemas.Routine
	result := parsedOptions.DB.
		Model(&restoredRoutine).
		Scopes(r.routineScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where("\"RoutineTable\".id = ?", id).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredRoutine.Id == uuid.Nil, Second: exceptions.Routine.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredRoutine, nil
}

func (r *RoutineRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Routine, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.Routine.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredRoutines []schemas.Routine
	result := parsedOptions.DB.
		Model(&restoredRoutines).
		Scopes(r.routineScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where("\"RoutineTable\".id IN ?", ids).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredRoutines) == 0, Second: exceptions.Routine.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredRoutines, nil
}

func (r *RoutineRepository) SoftDeleteOneById(
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
	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RoutineTable\".id = ?", id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Routine.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RoutineTable\".id IN ?", ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineRepository) HardDeleteOneById(
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

	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RoutineTable\".id = ?", id).
		Delete(&schemas.Routine{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *RoutineRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Routine.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.
		Model(&schemas.Routine{}).
		Scopes(r.routineScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.routineScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("\"RoutineTable\".id IN ?", ids).
		Delete(&schemas.Routine{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
