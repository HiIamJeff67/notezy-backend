package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
)

type UsersToShelvesRepositoryInterface interface {
	GetOne(rootShelfId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	GetMany(rootShelfId uuid.UUID, userIds []uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.UsersToShelves, *exceptions.Exception)
	GetManyByRootShelfIdsAndUserId(rootShelfIds []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.UsersToShelves, *exceptions.Exception)
	CreateOne(rootShelfId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	UpsertMany(rootShelfId uuid.UUID, userIds []uuid.UUID, permissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.UsersToShelves, *exceptions.Exception)
	UpdateOne(rootShelfId uuid.UUID, userId uuid.UUID, permission enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.UsersToShelves, *exceptions.Exception)
	DeleteOne(rootShelfId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	DeleteMany(rootShelfId uuid.UUID, userIds []uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	DeleteManyByRootShelfIdsAndUserId(rootShelfIds []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
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

func (r *UsersToShelvesRepository) GetMany(
	rootShelfId uuid.UUID,
	userIds []uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var relations []schemas.UsersToShelves
	result := parsedOptions.DB.
		Model(&schemas.UsersToShelves{}).
		Where("root_shelf_id = ? AND user_id IN ?", rootShelfId, userIds).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&relations)
	if result.Error != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(result.Error)
	}

	return relations, nil
}

func (r *UsersToShelvesRepository) GetManyByRootShelfIdsAndUserId(
	rootShelfIds []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var relations []schemas.UsersToShelves
	result := parsedOptions.DB.
		Model(&schemas.UsersToShelves{}).
		Where("root_shelf_id IN ? AND user_id = ?", rootShelfIds, userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&relations)
	if result.Error != nil {
		return nil, exceptions.Shelf.NotFound().WithOrigin(result.Error)
	}

	return relations, nil
}

func (r *UsersToShelvesRepository) CreateOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	relation := schemas.UsersToShelves{
		RootShelfId: rootShelfId,
		UserId:      userId,
		Permission:  permission,
	}
	result := parsedOptions.DB.Create(&relation)
	if result.Error != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	return r.GetOne(
		rootShelfId,
		userId,
		options.WithDB(parsedOptions.DB),
	)
}

func (r *UsersToShelvesRepository) UpsertMany(
	rootShelfId uuid.UUID,
	userIds []uuid.UUID,
	permissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.UsersToShelves, *exceptions.Exception) {
	if len(userIds) != len(permissions) {
		return nil, exceptions.Shelf.InvalidInput("userIds and permissions must have equal lengths")
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	relations := make([]schemas.UsersToShelves, len(userIds))
	for index, userId := range userIds {
		relations[index] = schemas.UsersToShelves{
			RootShelfId: rootShelfId,
			UserId:      userId,
			Permission:  permissions[index],
		}
	}

	result := parsedOptions.DB.
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "user_id",
					},
					{
						Name: "root_shelf_id",
					},
				},
				DoUpdates: clause.AssignmentColumns([]string{"permission", "updated_at"}),
			},
			clause.Returning{},
		).
		CreateInBatches(&relations, parsedOptions.BatchSize)
	if result.Error != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)
	}

	return relations, nil
}

func (r *UsersToShelvesRepository) UpdateOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToShelves, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	relation := schemas.UsersToShelves{
		RootShelfId: rootShelfId,
		UserId:      userId,
		Permission:  permission,
	}
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

	return r.GetOne(
		rootShelfId,
		userId,
		options.WithDB(parsedOptions.DB),
	)
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

func (r *UsersToShelvesRepository) DeleteMany(
	rootShelfId uuid.UUID,
	userIds []uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.
		Where("root_shelf_id = ? AND user_id IN ?", rootShelfId, userIds).
		Delete(&schemas.UsersToShelves{})
	if result.Error != nil {
		return exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected != int64(len(userIds)) {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *UsersToShelvesRepository) DeleteManyByRootShelfIdsAndUserId(
	rootShelfIds []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.
		Where("root_shelf_id IN ? AND user_id = ?", rootShelfIds, userId).
		Delete(&schemas.UsersToShelves{})
	if result.Error != nil {
		return exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected != int64(len(rootShelfIds)) {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
