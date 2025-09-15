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

type MaterialRepositoryInterface interface {
	HasPermission(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission) bool
	GetOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) (*schemas.Material, *exceptions.Exception)
	CheckPermissionAndGetOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission) (*schemas.Material, *exceptions.Exception)
	CreateOne(rootShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateMaterialInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID, matchedMaterialType *enums.MaterialType, input inputs.PartialUpdateMaterialInput) (*schemas.Material, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, rootShelfId uuid.UUID, userId uuid.UUID) *exceptions.Exception
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

func (r *MaterialRepository) HasPermission(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	allowedPermission []enums.AccessControlPermission,
) bool {
	var count int64 = 0
	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTabe\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermission,
		).
		Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return true
}

func (r *MaterialRepository) GetOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) (*schemas.Material, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	material := schemas.Material{}

	result := r.db.Model(&schemas.Material{}).
		Select("\"MaterialTable\".*").
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ?  AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermissions,
		).
		First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) (*schemas.Material, *exceptions.Exception) {
	material := schemas.Material{}

	result := r.db.Model(&schemas.Material{}).
		Select("\"MaterialTable\".*").
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermissions,
		).
		First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CreateOne(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateMaterialInput,
) (*uuid.UUID, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	shelfRepository := NewRootShelfRepository(r.db)
	hasPermission := shelfRepository.HasPermission(rootShelfId, userId, allowedPermissions)
	if !hasPermission {
		return nil, exceptions.Material.NoPermission("create")
	}

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

func (r *MaterialRepository) UpdateOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	matchedMaterialType *enums.MaterialType,
	input inputs.PartialUpdateMaterialInput,
) (*schemas.Material, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	// get and check the permission of the current user to the source shelf
	existingMaterial, exception := r.CheckPermissionAndGetOneById(id, rootShelfId, userId, allowedPermissions)
	if exception != nil || existingMaterial == nil {
		return nil, exception
	}

	// check if the material type is matched
	if matchedMaterialType != nil && existingMaterial.Type != *matchedMaterialType {
		return nil, exceptions.Material.MaterialTypeNotMatch(existingMaterial.Id.String(), existingMaterial.Type, matchedMaterialType)
	}

	// if the root shelf id is required to be updated in the database
	if input.Values.RootShelfId != nil && (input.SetNull == nil || !(*input.SetNull)["RootShelfId"]) {
		shelfRepository := NewRootShelfRepository(r.db)
		// check if the user has the enough permission to the destination shelf
		if hasPermissionOfNewShelf := shelfRepository.HasPermission(*input.Values.RootShelfId, userId, allowedPermissions); !hasPermissionOfNewShelf {
			return nil, exceptions.Shelf.NoPermission()
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingMaterial)
	}

	result := r.db.Model(&schemas.Material{}).
		Where("id = ?", id). // no need to check the permission here, since we have done that part on the above
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

func (r *MaterialRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermissions,
		).
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

func (r *MaterialRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			ids, rootShelfId, userId, allowedPermissions,
		).
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

func (r *MaterialRepository) SoftDeleteOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermissions,
		).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			ids, rootShelfId, userId, allowedPermissions,
		).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteOneById(
	id uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id = ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, rootShelfId, userId, allowedPermissions,
		).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("\"MaterialTable\".id IN ? AND \"MaterialTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			ids, rootShelfId, userId, allowedPermissions,
		).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}
