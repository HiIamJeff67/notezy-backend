package repositories

import (
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
)

type UsersToShelvesRepositoryInterface interface {
	GetOne(rootShelfId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	CreateOne(rootShelfId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	UpdatePermission(rootShelfId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	DeleteOne(rootShelfId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UsersToShelvesRepository struct{}

func NewUsersToShelvesRepository() UsersToShelvesRepositoryInterface {
	return &UsersToShelvesRepository{}
}

func (r *UsersToShelvesRepository) GetOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var relation schemas.UsersToShelves
	result := parsedOptions.DB.
		Model(&schemas.UsersToShelves{}).
		Preload(string(schemas.UsersToShelvesRelation_User)).
		Where("root_shelf_id = ? AND user_id = ?", rootShelfId, userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&relation)
	if result.Error != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(result.Error)
	}

	return &relation, nil
}

func (r *UsersToShelvesRepository) CreateOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	relation := schemas.UsersToShelves{RootShelfId: rootShelfId, UserId: userId, Permission: permission}
	result := parsedOptions.DB.Create(&relation)
	if result.Error != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	return r.GetOne(rootShelfId, userId, options.WithDB(parsedOptions.DB))
}

func (r *UsersToShelvesRepository) UpdatePermission(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	relation := schemas.UsersToShelves{RootShelfId: rootShelfId, UserId: userId, Permission: permission}
	result := parsedOptions.DB.
		Model(&relation).
		Select("permission").
		Updates(&relation)
	if result.Error != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Shelf.NotFound()
	}

	return r.GetOne(rootShelfId, userId, options.WithDB(parsedOptions.DB))
}

func (r *UsersToShelvesRepository) DeleteOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	result := parsedOptions.DB.
		Where("root_shelf_id = ? AND user_id = ?", rootShelfId, userId).
		Delete(&schemas.UsersToShelves{})
	if result.Error != nil {
		return exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
