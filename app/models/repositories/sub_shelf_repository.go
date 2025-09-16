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

/* ============================== Interface & Instance ============================== */

type SubShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission) bool
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads *[]schemas.SubShelfRelation) (*schemas.SubShelf, *exceptions.Exception)
	GetAllByRootShelfId(rootShelfId uuid.UUID, userId uuid.UUID, preloads *[]schemas.SubShelfRelation) (*[]schemas.SubShelf, *exceptions.Exception)
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads *[]schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission) (*schemas.SubShelf, *exceptions.Exception)
	CheckPermissionAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads *[]schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission) (*[]schemas.SubShelf, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateSubShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateSubShelfInput) (*schemas.SubShelf, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type SubShelfRepository struct {
	db *gorm.DB
}

func NewSubShelfRepository(db *gorm.DB) SubShelfRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &SubShelfRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *SubShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) bool {
	var count int64 = 0
	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rs.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return true
}

func (r *SubShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.SubShelfRelation,
) (*schemas.SubShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelf := schemas.SubShelf{}
	db := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		First(&subShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate()
	}

	return &subShelf, nil
}

func (r *SubShelfRepository) GetAllByRootShelfId(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.SubShelfRelation,
) (*[]schemas.SubShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelves := []schemas.SubShelf{}
	db := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rs.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"SubShelfTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?", rootShelfId, userId, allowedPermissions).
		Find(&subShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &subShelves, nil
}

func (r *SubShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
) (*schemas.SubShelf, *exceptions.Exception) {
	subShelf := schemas.SubShelf{}
	db := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		First(&subShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &subShelf, nil
}

func (r *SubShelfRepository) CheckPermissionAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
) (*[]schemas.SubShelf, *exceptions.Exception) {
	subShelves := []schemas.SubShelf{}
	db := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"SubShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Find(&subShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &subShelves, nil
}

func (r *SubShelfRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateSubShelfInput,
) (*uuid.UUID, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	rootShelfrepository := NewRootShelfRepository(r.db)
	if hasPermission := rootShelfrepository.HasPermission(input.RootShelfId, userId, allowedPermissions); !hasPermission {
		return nil, exceptions.Shelf.NoPermission()
	}

	var newSubShelf schemas.SubShelf
	if len(input.Path) == 0 {
		newSubShelf.PrevSubShelfId = nil
	} else {
		newSubShelf.PrevSubShelfId = &input.Path[len(input.Path)-1]
	}
	if err := copier.Copy(&newSubShelf, &input); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Create(&newSubShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	return &newSubShelf.Id, nil
}

func (r *SubShelfRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateSubShelfInput,
) (*schemas.SubShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	existingSubShelf, exception := r.CheckPermissionAndGetOneById(id, userId, nil, allowedPermissions)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingSubShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingSubShelf).WithError(err)
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
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
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
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
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
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
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *SubShelfRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
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
	id uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
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
	ids []uuid.UUID,
	userId uuid.UUID,
) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rc.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Delete(&schemas.SubShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
