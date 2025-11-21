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

type SubShelfRepositoryInterface interface {
	HasPermission(db *gorm.DB, id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	HasPermissions(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	CheckPermissionAndGetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*schemas.SubShelf, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]schemas.SubShelf, *exceptions.Exception)
	GetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation) (*schemas.SubShelf, *exceptions.Exception)
	GetAllByRootShelfId(db *gorm.DB, rootShelfId uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation) ([]schemas.SubShelf, *exceptions.Exception)
	CreateOneByRootShelfId(db *gorm.DB, rootShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateSubShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateSubShelfInput) (*schemas.SubShelf, *exceptions.Exception)
	RestoreSoftDeletedOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type SubShelfRepository struct{}

func NewSubShelfRepository() SubShelfRepositoryInterface {
	return &SubShelfRepository{}
}

/* ============================== Implementations ============================== */

func (r *SubShelfRepository) HasPermission(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) bool {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return count > 0
}

func (r *SubShelfRepository) HasPermissions(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) bool {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	query := db.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return count > 0
}

func (r *SubShelfRepository) CheckPermissionAndGetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*schemas.SubShelf, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subShelf := schemas.SubShelf{}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"SubShelfTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"SubShelfTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.First(&subShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &subShelf, nil
}

func (r *SubShelfRepository) CheckPermissionsAndGetManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]schemas.SubShelf, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subShelves := []schemas.SubShelf{}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"SubShelfTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"SubShelfTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&subShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}
	if len(subShelves) == 0 {
		return nil, exceptions.Shelf.NotFound()
	}

	return subShelves, nil
}

func (r *SubShelfRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
) (*schemas.SubShelf, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(
		db,
		id,
		userId,
		preloads,
		allowedPermissions,
		types.Ternary_Negative,
	)
}

func (r *SubShelfRepository) GetAllByRootShelfId(
	db *gorm.DB,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
) ([]schemas.SubShelf, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelves := []schemas.SubShelf{}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.SubShelf{}).
		Where("root_shelf_id = ? AND EXISTS (?)", rootShelfId, subQuery)
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&subShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return subShelves, nil
}

func (r *SubShelfRepository) CreateOneByRootShelfId(
	db *gorm.DB,
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateSubShelfInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	var newSubShelf schemas.SubShelf
	if input.PrevSubShelfId != nil {
		prevSubShelf, exception := r.CheckPermissionAndGetOneById(
			db,
			*input.PrevSubShelfId,
			userId,
			nil,
			allowedPermissions,
			types.Ternary_Negative,
		)
		if exception != nil {
			return nil, exception
		}
		prevSubShelf.Path = append(prevSubShelf.Path, prevSubShelf.Id)
		newSubShelf.Path = prevSubShelf.Path
	}

	if err := copier.Copy(&newSubShelf, &input); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}
	newSubShelf.RootShelfId = rootShelfId

	result := db.Model(&schemas.SubShelf{}).
		Create(&newSubShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	return &newSubShelf.Id, nil
}

func (r *SubShelfRepository) UpdateOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateSubShelfInput,
) (*schemas.SubShelf, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingSubShelf, exception := r.CheckPermissionAndGetOneById(
		db,
		id,
		userId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingSubShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingSubShelf).WithError(err)
	}

	result := db.Model(&schemas.SubShelf{}).
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

func (r *SubShelfRepository) RestoreSoftDeletedOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NOT NULL", id, subQuery).
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

func (r *SubShelfRepository) RestoreSoftDeletedManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NOT NULL", ids, subQuery).
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

func (r *SubShelfRepository) SoftDeleteOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NULL", id, subQuery).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NoChanges()
	}

	return nil
}

func (r *SubShelfRepository) SoftDeleteManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NULL", ids, subQuery).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NOT NULL", id, subQuery).
		Delete(&schemas.SubShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := db.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NOT NULL", ids, subQuery).
		Delete(&schemas.SubShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
