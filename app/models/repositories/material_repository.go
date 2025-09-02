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
	"notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type MaterialRepositoryInterface interface {
	GetOneById(id uuid.UUID, userId uuid.UUID) (*schemas.Material, *exceptions.Exception)
	CreateOne(input inputs.CreateMaterialInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateMaterialInput) (*schemas.Material, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) MaterialRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &MaterialRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *MaterialRepository) GetOneById(id uuid.UUID, userId uuid.UUID) (*schemas.Material, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	material := schemas.Material{}

	result := r.db.Model(&schemas.Material{}).
		Select("\"MaterialTable\".*").
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelves\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CreateOne(input inputs.CreateMaterialInput) (*uuid.UUID, *exceptions.Exception) {
	var newMaterial schemas.Material
	if err := copier.Copy(&newMaterial, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := r.db.Model(&schemas.Material{}).
		Create(&newMaterial)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToCreate().WithError(err)
	}

	return &newMaterial.Id, nil
}

func (r *MaterialRepository) UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateMaterialInput) (*schemas.Material, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	existingMaterial, exception := r.GetOneById(id, userId) // get and check the permission of the current user
	if exception != nil || existingMaterial == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingMaterial)
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Material.NoChanges()
	}

	return &updates, nil
}

func (r *MaterialRepository) RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}
