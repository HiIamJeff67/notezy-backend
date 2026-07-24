package repositories

import (
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
)

type UsersToStationsRepositoryInterface interface {
	GetOne(stationId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UsersToStations, *exceptions.Exception)
	CreateOne(stationId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToStations, *exceptions.Exception)
	UpdatePermission(stationId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToStations, *exceptions.Exception)
	DeleteOne(stationId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UsersToStationsRepository struct{}

func NewUsersToStationsRepository() UsersToStationsRepositoryInterface {
	return &UsersToStationsRepository{}
}

func (r *UsersToStationsRepository) GetOne(
	stationId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToStations, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var relation schemas.UsersToStations
	result := parsedOptions.DB.
		Model(&schemas.UsersToStations{}).
		Preload(string(schemas.UsersToStationsRelation_User)).
		Where("station_id = ? AND user_id = ?", stationId, userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&relation)
	if result.Error != nil {
		return nil, exceptions.Station.NotFound().WithOrigin(result.Error)
	}

	return &relation, nil
}

func (r *UsersToStationsRepository) CreateOne(
	stationId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToStations, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	relation := schemas.UsersToStations{StationId: stationId, UserId: userId, Permission: permission}
	result := parsedOptions.DB.Create(&relation)
	if result.Error != nil {
		return nil, exceptions.Station.FailedToCreate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Station.NoChanges()
	}

	return r.GetOne(stationId, userId, options.WithDB(parsedOptions.DB))
}

func (r *UsersToStationsRepository) UpdatePermission(
	stationId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToStations, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	relation := schemas.UsersToStations{StationId: stationId, UserId: userId, Permission: permission}
	result := parsedOptions.DB.
		Model(&relation).
		Select("permission").
		Updates(&relation)
	if result.Error != nil {
		return nil, exceptions.Station.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Station.NotFound()
	}

	return r.GetOne(stationId, userId, options.WithDB(parsedOptions.DB))
}

func (r *UsersToStationsRepository) DeleteOne(
	stationId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	result := parsedOptions.DB.
		Where("station_id = ? AND user_id = ?", stationId, userId).
		Delete(&schemas.UsersToStations{})
	if result.Error != nil {
		return exceptions.Station.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return exceptions.Station.NotFound()
	}

	return nil
}
