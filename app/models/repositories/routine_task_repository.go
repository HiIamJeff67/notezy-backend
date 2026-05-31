package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
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

type RoutineTaskRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.RoutineTask, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.RoutineTask, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RoutineTaskRelation, opts ...options.RepositoryOptions) (*schemas.RoutineTask, *exceptions.Exception)
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
		Clauses(clause.Locking{Strength: "SHARE"}).
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
		Select("DISTINCT \"RoutineTaskTable\".id").
		Scopes(r.routineTaskScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Clauses(clause.Locking{Strength: "SHARE"}).
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
		Clauses(clause.Locking{Strength: "SHARE"}).
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
		Clauses(clause.Locking{Strength: "SHARE"}).
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
		return nil, exceptions.RoutineTask.NoPermission("create a routine task under this station")
	}

	newRoutineTask := schemas.RoutineTask{
		Status:      enums.RoutineTaskStatus_Idle,
		ScheduledAt: time.Now().Truncate(time.Minute),
	}
	if err := copier.Copy(&newRoutineTask, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.RoutineTask.InvalidInput().WithOrigin(err)
	}
	newRoutineTask.StationId = stationId
	newRoutineTask.ScheduledAt = time.Now().Truncate(time.Minute)
	if len(newRoutineTask.Payload) == 0 {
		newRoutineTask.Payload = datatypes.JSON([]byte("{}"))
	}

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

	newRoutineTasks := make([]schemas.RoutineTask, 0, len(input))
	for _, in := range input {
		if !isStationValid[in.StationId] {
			continue
		}
		newRoutineTask := schemas.RoutineTask{
			Id:          uuid.New(),
			Status:      enums.RoutineTaskStatus_Idle,
			Attempts:    0,
			MaxAttempts: 1,
			ScheduledAt: time.Now().Truncate(time.Minute),
		}
		if err := copier.Copy(&newRoutineTask, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.RoutineTask.InvalidInput().WithOrigin(err)
		}
		newRoutineTask.ScheduledAt = time.Now().Truncate(time.Minute)
		if len(newRoutineTask.Payload) == 0 {
			newRoutineTask.Payload = datatypes.JSON([]byte("{}"))
		}
		if newRoutineTask.MaxAttempts <= 0 {
			newRoutineTask.MaxAttempts = 1
		}
		if newRoutineTask.Purpose == "" {
			newRoutineTask.Purpose = enums.RoutineTaskPurpose_CreateBlockPack
		}
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
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRoutineTask)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingRoutineTask).WithOrigin(err)
	}

	result := parsedOptions.DB.
		Model(&schemas.RoutineTask{}).
		Where("\"RoutineTaskTable\".id = ?", id).
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
	validRoutineTasks, exception := r.CheckPermissionsAndGetManyByIds(ids, userId, nil, allowedPermissions, opts...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exceptions.RoutineTask.NoPermission("update these routine tasks")
	}

	routineTaskById := make(map[uuid.UUID]schemas.RoutineTask, len(validRoutineTasks))
	for _, validRoutineTask := range validRoutineTasks {
		routineTaskById[validRoutineTask.Id] = validRoutineTask
	}
	stationRepository := NewStationRepository(scopes.NewStationScope())

	for _, in := range input {
		existingRoutineTask, exist := routineTaskById[in.Id]
		if !exist {
			continue
		}
		if in.PartialUpdateInput.Values.StationId != nil &&
			(in.PartialUpdateInput.SetNull == nil || !(*in.PartialUpdateInput.SetNull)["StationId"]) &&
			!stationRepository.HasPermission(*in.PartialUpdateInput.Values.StationId, userId, allowedPermissions, opts...) {
			parsedOptions.DB.Rollback()
			return exceptions.RoutineTask.NoPermission("move these routine tasks to the given stations")
		}
		updates, err := util.PartialUpdatePreprocess(in.PartialUpdateInput.Values, in.PartialUpdateInput.SetNull, existingRoutineTask)
		if err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Util.FailedToPreprocessPartialUpdate(
				in.PartialUpdateInput.Values,
				in.PartialUpdateInput.SetNull,
				existingRoutineTask,
			).WithOrigin(err)
		}
		result := parsedOptions.DB.
			Model(&schemas.RoutineTask{}).
			Where("\"RoutineTaskTable\".id = ?", in.Id).
			Select("*").
			Updates(&updates)
		if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
			{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)},
			{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
		}); exception != nil {
			parsedOptions.DB.Rollback()
			return exception
		}
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
		Where("\"RoutineTaskTable\".id = ?", id).
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
		Where("\"RoutineTaskTable\".id IN ?", ids).
		Delete(&schemas.RoutineTask{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.RoutineTask.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.RoutineTask.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
