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
	types "notezy-backend/shared/types"
)

/* ============================== Definitions ============================== */

type MaterialRepositoryInterface interface {
	HasPermission(db *gorm.DB, id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	HasPermissions(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	CheckPermissionAndGetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*schemas.Material, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]schemas.Material, *exceptions.Exception)
	GetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) (*schemas.Material, *exceptions.Exception)
	CreateOne(db *gorm.DB, subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateMaterialInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, matchedMaterialType *enums.MaterialType, input inputs.PartialUpdateMaterialInput) (*schemas.Material, *exceptions.Exception)

	RestoreSoftDeletedOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type MaterialRepository struct{}

func NewMaterialRepository() MaterialRepositoryInterface {
	return &MaterialRepository{}
}

/* ============================== CRUD operations ============================== */

func (r *MaterialRepository) HasPermission(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) bool {
	if db == nil {
		db = models.NotezyDB
	}

	var count int64 = 0

	query := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, userId, allowedPermissions,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"MaterialTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"MaterialTable\".deleted_at IS NULL")
	}

	result := query.Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return count > 0
}

func (r *MaterialRepository) HasPermissions(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) bool {
	if db == nil {
		db = models.NotezyDB
	}

	var count int64 = 0

	query := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("\"MaterialTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?",
			ids, userId, allowedPermissions,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"MaterialTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"MaterialTable\".deleted_at IS NULL")
	}

	result := query.Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return count > 0
}

func (r *MaterialRepository) CheckPermissionAndGetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*schemas.Material, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	material := schemas.Material{}

	query := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("\"MaterialTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?",
			id, userId, allowedPermissions,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"MaterialTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"MaterialTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CheckPermissionsAndGetManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]schemas.Material, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	materials := []schemas.Material{}

	query := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("\"MaterialTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?",
			ids, userId, allowedPermissions,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"MaterialTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"MaterialTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}
	if len(materials) == 0 {
		return nil, exceptions.Material.NotFound()
	}

	return materials, nil
}

func (r *MaterialRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) (*schemas.Material, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	return r.CheckPermissionAndGetOneById(
		db,
		id,
		userId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
}

func (r *MaterialRepository) CreateOne(
	db *gorm.DB,
	subShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateMaterialInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelfRepository := NewSubShelfRepository()

	if hasPermission := subShelfRepository.HasPermission(
		db,
		subShelfId,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	); !hasPermission {
		return nil, exceptions.Material.NoPermission("create")
	}

	var newMaterial schemas.Material
	if err := copier.Copy(&newMaterial, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := db.Model(&schemas.Material{}).
		Create(&newMaterial)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToCreate().WithError(err)
	}

	return &newMaterial.Id, nil
}

func (r *MaterialRepository) UpdateOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	matchedMaterialType *enums.MaterialType,
	input inputs.PartialUpdateMaterialInput,
) (*schemas.Material, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	// get and check the permission of the current user to the source shelf
	existingMaterial, exception := r.CheckPermissionAndGetOneById(
		db,
		id,
		userId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
	if exception != nil || existingMaterial == nil {
		return nil, exception
	}

	// check if the material type is matched
	if matchedMaterialType != nil && existingMaterial.Type != *matchedMaterialType {
		return nil, exceptions.Material.MaterialTypeNotMatch(
			existingMaterial.Id.String(),
			existingMaterial.Type,
			matchedMaterialType,
		)
	}

	// if the root shelf id is required to be updated in the database
	if input.Values.ParentSubShelfId != nil && (input.SetNull == nil || !(*input.SetNull)["ParentSubShelfId"]) {
		subShelfRepository := NewSubShelfRepository()
		// check if the user has the enough permission to the destination shelf
		if hasPermissionOfNewSubShelf := subShelfRepository.HasPermission(
			db,
			*input.Values.ParentSubShelfId,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		); !hasPermissionOfNewSubShelf {
			return nil, exceptions.Shelf.NoPermission()
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingMaterial,
		)
	}

	result := db.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NULL", id). // no need to check the permission here, since we have done that part on the above
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
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Positive,
	); !hasPermission {
		return exceptions.Material.NoPermission("restore")
	}

	result := db.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NoChanges()
	}

	return nil
}

func (r *MaterialRepository) RestoreSoftDeletedManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Positive,
	); !hasPermission {
		return exceptions.Material.NoPermission("restore")
	}

	result := db.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
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
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	); !hasPermission {
		return exceptions.Material.NoPermission("soft delete")
	}

	result := db.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NULL", id).
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
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	); !hasPermission {
		return exceptions.Material.NoPermission("soft delete")
	}

	result := db.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
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
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Positive,
	); !hasPermission {
		return exceptions.Material.NoPermission("hard delete")
	}

	result := db.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
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
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Positive,
	); !hasPermission {
		return exceptions.Material.NoPermission("hard delete")
	}

	result := db.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}
