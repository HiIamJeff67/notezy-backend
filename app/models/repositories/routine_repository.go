package repositories

import (
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
	routineById := make(map[uuid.UUID]schemas.Routine, len(validRoutines))
	for _, validRoutine := range validRoutines {
		routineById[validRoutine.Id] = validRoutine
	}
	stationRepository := NewStationRepository(scopes.NewStationScope())

	for _, in := range input {
		existingRoutine, exist := routineById[in.Id]
		if !exist {
			continue
		}
		if in.PartialUpdateInput.Values.StationId != nil &&
			(in.PartialUpdateInput.SetNull == nil || !(*in.PartialUpdateInput.SetNull)["StationId"]) &&
			!stationRepository.HasPermission(*in.PartialUpdateInput.Values.StationId, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return exceptions.Routine.NoPermission("move these routines to the given stations")
		}
		updates, err := util.PartialUpdatePreprocess(in.PartialUpdateInput.Values, in.PartialUpdateInput.SetNull, existingRoutine)
		if err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Util.FailedToPreprocessPartialUpdate(
				in.PartialUpdateInput.Values,
				in.PartialUpdateInput.SetNull,
				existingRoutine,
			).WithOrigin(err)
		}
		result := parsedOptions.DB.
			Model(&schemas.Routine{}).
			Where("\"RoutineTable\".id = ? AND \"RoutineTable\".deleted_at IS NULL", in.Id).
			Select("*").
			Updates(&updates)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.Routine.FailedToUpdate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return exception
		}
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
