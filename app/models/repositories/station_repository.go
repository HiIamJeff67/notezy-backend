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

type StationRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.StationRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Station, enums.AccessControlPermission, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.StationRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Station, []enums.AccessControlPermission, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.StationRelation, opts ...options.RepositoryOptions) (*schemas.Station, enums.AccessControlPermission, *exceptions.Exception)
	GetAllByUserId(userId uuid.UUID, preloads []schemas.StationRelation, opts ...options.RepositoryOptions) ([]schemas.Station, []enums.AccessControlPermission, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateStationInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	CreateManyByOwnerId(ownerId uuid.UUID, input []inputs.CreateStationInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateStationInput, opts ...options.RepositoryOptions) (*schemas.Station, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateStationInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Station, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Station, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type StationRepository struct {
	stationScope scopes.StationScopeInterface
}

func NewStationRepository(stationScope scopes.StationScopeInterface) StationRepositoryInterface {
	return &StationRepository{
		stationScope: stationScope,
	}
}

func (r *StationRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.Station{}).
		Select("1").
		Scopes(r.stationScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *StationRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.Station{}).
		Select(`DISTINCT "StationTable".id`).
		Scopes(r.stationScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *StationRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.StationRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Station, enums.AccessControlPermission, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var station schemas.Station
	result := parsedOptions.DB.
		Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.stationScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&station)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.NotFound().WithOrigin(result.Error)},
		{First: station.Id == uuid.Nil, Second: exceptions.Station.NotFound()},
	}); exception != nil {
		return nil, "", exception
	}

	var permission enums.AccessControlPermission
	result = parsedOptions.DB.
		Model(&schemas.UsersToStations{}).
		Select("permission").
		Where(
			"station_id = ? AND user_id = ? AND permission IN ?",
			station.Id,
			userId,
			allowedPermissions,
		).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&permission)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.NotFound().WithOrigin(result.Error)},
		{First: permission == "", Second: exceptions.Station.NotFound()},
	}); exception != nil {
		return nil, "", exception
	}

	return &station, permission, nil
}

func (r *StationRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.StationRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.Station, []enums.AccessControlPermission, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var stations []schemas.Station
	result := parsedOptions.DB.
		Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.stationScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&stations)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.NotFound().WithOrigin(result.Error)},
		{First: len(stations) == 0, Second: exceptions.Station.NotFound()},
	}); exception != nil {
		return nil, nil, exception
	}

	var usersToStations []schemas.UsersToStations
	result = parsedOptions.DB.
		Model(&schemas.UsersToStations{}).
		Select("station_id, permission").
		Where(
			"station_id IN ? AND user_id = ? AND permission IN ?",
			ids,
			userId,
			allowedPermissions,
		).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&usersToStations)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.NotFound().WithOrigin(result.Error)},
		{First: len(usersToStations) == 0, Second: exceptions.Station.NotFound()},
	}); exception != nil {
		return nil, nil, exception
	}

	permissionByStationId := make(map[uuid.UUID]enums.AccessControlPermission, len(usersToStations))
	for _, usersToStation := range usersToStations {
		permissionByStationId[usersToStation.StationId] = usersToStation.Permission
	}

	permissions := make([]enums.AccessControlPermission, len(stations))
	for index, station := range stations {
		permission, exist := permissionByStationId[station.Id]
		if !exist {
			return nil, nil, exceptions.Station.NotFound()
		}
		permissions[index] = permission
	}

	return stations, permissions, nil
}

func (r *StationRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.StationRelation,
	opts ...options.RepositoryOptions,
) (*schemas.Station, enums.AccessControlPermission, *exceptions.Exception) {
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

func (r *StationRepository) GetAllByUserId(
	userId uuid.UUID,
	preloads []schemas.StationRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.Station, []enums.AccessControlPermission, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	type stationWithPermission struct {
		schemas.Station
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	var stationsWithPermissions []stationWithPermission
	result := parsedOptions.DB.
		Model(&schemas.Station{}).
		Select(`"StationTable".*, uts.permission AS permission`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "StationTable".id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.stationScope.IncludePreloads(preloads)).
		Order(`"StationTable".created_at ASC`).
		Order(`"StationTable".id ASC`).
		Find(&stationsWithPermissions)
	if result.Error != nil {
		return nil, nil, exceptions.Station.NotFound().WithOrigin(result.Error)
	}

	stations := make([]schemas.Station, len(stationsWithPermissions))
	permissions := make([]enums.AccessControlPermission, len(stationsWithPermissions))
	for index, stationWithPermission := range stationsWithPermissions {
		stations[index] = stationWithPermission.Station
		permissions[index] = stationWithPermission.Permission
	}

	return stations, permissions, nil
}

func (r *StationRepository) CreateOneByOwnerId(
	ownerId uuid.UUID,
	input inputs.CreateStationInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	var newStation schemas.Station
	newStation.OwnerId = ownerId
	if err := copier.Copy(&newStation, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Station.InvalidInput().WithOrigin(err)
	}
	if newStation.Id == uuid.Nil {
		newStation.Id = uuid.New()
	}

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newStation)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToCreate().WithOrigin(result.Error)},
		{First: newStation.Id == uuid.Nil, Second: exceptions.Station.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newUsersToStations := schemas.UsersToStations{
		UserId:     ownerId,
		StationId:  newStation.Id,
		Permission: enums.AccessControlPermission_Owner,
	}
	result = parsedOptions.DB.Model(&schemas.UsersToStations{}).
		Create(&newUsersToStations)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newStation.Id, nil
}

func (r *StationRepository) CreateManyByOwnerId(
	ownerId uuid.UUID,
	input []inputs.CreateStationInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.Station.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
	}

	newStations := make([]schemas.Station, 0, len(input))
	for _, in := range input {
		var newStation schemas.Station
		newStation.OwnerId = ownerId
		if err := copier.Copy(&newStation, &in); err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Station.InvalidInput().WithOrigin(err)
		}
		if newStation.Id == uuid.Nil {
			newStation.Id = uuid.New()
		}
		newStations = append(newStations, newStation)
	}

	result := parsedOptions.DB.Model(&schemas.Station{}).
		CreateInBatches(&newStations, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	newStationIds := make([]uuid.UUID, len(newStations))
	newUsersToStations := make([]schemas.UsersToStations, len(newStations))
	for index, newStation := range newStations {
		newStationIds[index] = newStation.Id
		newUsersToStations[index] = schemas.UsersToStations{
			UserId:     ownerId,
			StationId:  newStation.Id,
			Permission: enums.AccessControlPermission_Owner,
		}
	}
	result = parsedOptions.DB.Model(&schemas.UsersToStations{}).
		CreateInBatches(&newUsersToStations, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return newStationIds, nil
}

func (r *StationRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateStationInput,
	opts ...options.RepositoryOptions,
) (*schemas.Station, *exceptions.Exception) {
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

	existingStation, _, exception := r.CheckPermissionAndGetOneById(
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingStation)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingStation).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Where(`"StationTable".id = ? AND "StationTable".deleted_at IS NULL`, id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *StationRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateStationInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(input) == 0 {
		return exceptions.Station.NoChanges()
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

	validStations, _, exception := r.CheckPermissionsAndGetManyByIds(
		ids,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return exceptions.Station.NoPermission("update these stations")
	}

	isStationValid := make(map[uuid.UUID]bool, len(validStations))
	for _, validStation := range validStations {
		isStationValid[validStation.Id] = true
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !isStationValid[in.Id] {
			continue
		}

		setIconNull := false
		setHeaderBackgroundURLNull := false
		if in.PartialUpdateInput.SetNull != nil {
			for field, setNull := range *in.PartialUpdateInput.SetNull {
				if !setNull {
					continue
				}
				switch strings.ToLower(strings.ReplaceAll(field, "_", "")) {
				case "icon":
					setIconNull = true
				case "headerbackgroundurl":
					setHeaderBackgroundURLNull = true
				}
			}
		}

		valuePlaceholders = append(valuePlaceholders, `(?::uuid, ?::text, ?::text, ?::"SupportedIcon", ?::text, ?::boolean, ?::boolean)`)
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Name,
			in.PartialUpdateInput.Values.Description,
			in.PartialUpdateInput.Values.Icon,
			in.PartialUpdateInput.Values.HeaderBackgroundURL,
			setIconNull,
			setHeaderBackgroundURLNull,
		)
	}

	if len(valuePlaceholders) == 0 {
		parsedOptions.DB.Rollback()
		return exceptions.Station.NoChanges()
	}

	sql := fmt.Sprintf(`
		UPDATE "StationTable" AS s
		SET
			name = COALESCE(v.name::text, s.name),
			description = COALESCE(v.description::text, s.description),
			icon = CASE
				WHEN v.set_icon_null::boolean THEN NULL
				ELSE COALESCE(v.icon::"SupportedIcon", s.icon)
			END,
			header_background_url = CASE
				WHEN v.set_header_background_url_null::boolean THEN NULL
				ELSE COALESCE(v.header_background_url::text, s.header_background_url)
			END,
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, name, description, icon, header_background_url, set_icon_null, set_header_background_url_null)
		WHERE s.id = v.id::uuid AND s.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return nil
}

func (r *StationRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Station, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredStation schemas.Station
	result := parsedOptions.DB.Model(&restoredStation).
		Scopes(r.stationScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where(`"StationTable".id = ?`, id).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredStation.Id == uuid.Nil, Second: exceptions.Station.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredStation, nil
}

func (r *StationRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Station, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.Station.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredStations []schemas.Station
	result := parsedOptions.DB.Model(&restoredStations).
		Scopes(r.stationScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Returning{}).
		Where(`"StationTable".id IN ?`, ids).
		Updates(map[string]interface{}{"deleted_at": nil})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredStations) == 0, Second: exceptions.Station.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredStations, nil
}

func (r *StationRepository) SoftDeleteOneById(
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

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"StationTable".id = ?`, id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *StationRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Station.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"StationTable".id IN ?`, ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *StationRepository) SoftDeleteManyByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("owner_id = ?", userId).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *StationRepository) HardDeleteOneById(
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

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"StationTable".id = ?`, id).
		Delete(&schemas.Station{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *StationRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Station.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"StationTable".id IN ?`, ids).
		Delete(&schemas.Station{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *StationRepository) HardDeleteManyByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.Station{}).
		Scopes(r.stationScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where("owner_id = ?", userId).
		Delete(&schemas.Station{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Station.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
