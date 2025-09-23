package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type RootShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RootShelfRelation, allowedPermissions []enums.AccessControlPermission, includeDeleted bool) (*schemas.RootShelf, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.RootShelfRelation) (*schemas.RootShelf, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateRootShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRootShelfInput) (*schemas.RootShelf, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type RootShelfRepository struct {
	db *gorm.DB
}

func NewRootShelfRepository(db *gorm.DB) RootShelfRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RootShelfRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *RootShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) bool {
	var count int64 = 0

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery).
		Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return count > 0
}

func (r *RootShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RootShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	includeDeleted bool,
) (*schemas.RootShelf, *exceptions.Exception) {
	rootShelf := schemas.RootShelf{}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	db := r.db.Model(&schemas.RootShelf{}).
		Where("\"RootShelfTable\".id = ? AND EXISTS (?)", id, subQuery)

	if !includeDeleted {
		db = db.Where("\"RootShelfTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.First(&rootShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &rootShelf, nil
}

func (r *RootShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.RootShelfRelation,
) (*schemas.RootShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	return r.CheckPermissionAndGetOneById(id, userId, preloads, allowedPermissions, false)
}

func (r *RootShelfRepository) CreateOneByOwnerId(
	ownerId uuid.UUID,
	input inputs.CreateRootShelfInput,
) (*uuid.UUID, *exceptions.Exception) {
	var newRootShelf schemas.RootShelf
	newRootShelf.OwnerId = ownerId
	if err := copier.Copy(&newRootShelf, &input); err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Create(&newRootShelf)
	if err := result.Error; err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"shelf_idx_owner_id_name\" (SQLSTATE 23505)":
			return nil, exceptions.Shelf.DuplicateName(input.Name)
		default:
			return nil, exceptions.Shelf.FailedToCreate().WithError(err)
		}
	}

	// create the users to shelves relation with the permission of admin
	newUsersToShelves := schemas.UsersToShelves{
		UserId:      ownerId,
		RootShelfId: newRootShelf.Id,
		Permission:  enums.AccessControlPermission_Admin,
	}
	result = r.db.Model(&schemas.UsersToShelves{}).
		Create(&newUsersToShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	return &newRootShelf.Id, nil
}

func (r *RootShelfRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRootShelfInput,
) (*schemas.RootShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	existingRootShelf, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		false,
	)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRootShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values, input.SetNull, *existingRootShelf,
		).WithError(err)
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	return &updates, nil
}

func (r *RootShelfRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NoChanges()
	}

	return nil
}

func (r *RootShelfRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery).
		Delete(&schemas.RootShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	subQuery := r.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := r.db.Model(&schemas.RootShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery).
		Delete(&schemas.RootShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
