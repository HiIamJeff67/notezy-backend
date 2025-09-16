package repositories

import (
	"fmt"
	"strings"
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
	constants "notezy-backend/shared/constants"
)

/* ============================== Definitions ============================== */

type RootShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission) bool
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads *[]schemas.RootShelfRelation) (*schemas.RootShelf, *exceptions.Exception)
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads *[]schemas.RootShelfRelation, allowedPermissions []enums.AccessControlPermission) (*schemas.RootShelf, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateRootShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRootShelfInput) (*schemas.RootShelf, *exceptions.Exception)
	DirectlyUpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateRootShelfInput) *exceptions.Exception
	DirectlyUpdateManyByIds(ids []uuid.UUID, userId uuid.UUID, inputs []inputs.PartialUpdateRootShelfInput) *exceptions.Exception
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

/* ============================== Helper functions ============================== */

func (r *RootShelfRepository) getPartialUpdatePlaceholderUnit() string {
	return "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
}

func (r *RootShelfRepository) getNumOfPartialUpdateArguments(numOfRows int) int {
	// include values.ShelfId, values.ShelfName, values.ShelfEncodedStructure, setNull.ShelfName, setNull.ShelfEncodedStructure
	// and reserved 1 for the ownerId.
	// Note that if its for batch updates, the ownerIds could be more than one, hence this doesn't work
	if numOfRows > (constants.MaxDatabaseUpdateParameters-1)/5 {
		return 0
	}
	return numOfRows*5 + 1
}

/* ============================== CRUD operations ============================== */

func (r *RootShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermission []enums.AccessControlPermission,
) bool {
	var count int64 = 0
	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermission).
		Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return true
}

func (r *RootShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.RootShelfRelation,
) (*schemas.RootShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	rootShelf := schemas.RootShelf{}
	db := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		First(&rootShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &rootShelf, nil
}

func (r *RootShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads *[]schemas.RootShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
) (*schemas.RootShelf, *exceptions.Exception) {
	rootShelf := schemas.RootShelf{}
	db := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		First(&rootShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &rootShelf, nil
}

func (r *RootShelfRepository) CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateRootShelfInput) (*uuid.UUID, *exceptions.Exception) {
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

	existingRootShelf, exception := r.CheckPermissionAndGetOneById(id, userId, nil, allowedPermissions)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingRootShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingRootShelf).WithError(err)
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Shelf.NoChanges()
	}

	return &updates, nil
}

func (r *RootShelfRepository) DirectlyUpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateRootShelfInput,
) *exceptions.Exception {
	return r.DirectlyUpdateManyByIds([]uuid.UUID{id}, userId, []inputs.PartialUpdateRootShelfInput{input})
}

func (r *RootShelfRepository) DirectlyUpdateManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	inputs []inputs.PartialUpdateRootShelfInput,
) *exceptions.Exception {
	if len(ids) != len(inputs) || len(ids) == 0 {
		return exceptions.Shelf.NoChanges()
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	placeholders := make([]string, 0, len(ids))
	args := make([]interface{}, 0, r.getNumOfPartialUpdateArguments(len(ids))) // the number of the arguments depends on the number of columns in partial update dto

	for index, id := range ids {
		placeholders = append(placeholders, r.getPartialUpdatePlaceholderUnit())
		args = append(args, id)

		// safetly dereference all the values
		args = append(args, util.DerefOrNil(inputs[index].Values.Name))
		args = append(args, util.DerefOrNil(inputs[index].Values.TotalShelfNodes))
		args = append(args, util.DerefOrNil(inputs[index].Values.TotalMaterials))

		// safetly dereference all the setNulls
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "Name"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "TotalShelfNodes"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "TotalMaterials"))
	}

	sql := fmt.Sprintf(`
			UPDATE "%s" AS s
			SET
				name = CASE
					WHEN v.set_null_name THEN NULL
					ELSE COALESCE(v.name, s.name)
				END,
				total_shelf_nodes = CASE
					WHEN v.total_shelf_nodes THEN NULL
					ELSE COALESCE(v.total_shelf_nodes, s.total_shelf_nodes)
				END, 
				total_materials = CASE
					WHEN v.total_materials THEN NULL
					ELSE COALESCE(v.total_materials, s.total_materials)
				END, 
				max_width = CASE
					WHEN v.max_width THEN NULL
					ELSE COALESCE(v.max_width, s.max_width)
				END, 
				max_depth = CASE
					WHEN v.max_depth THEN NULL
					ELSE COALESCE(v.max_depth, s.max_depth)
				END, 
				updated_at = NOW()
			FROM (VALUES %s) AS v(id, name, set_null_name)
			LEFT JOIN "UsersToShelvesTable" AS uts ON s.id = uts.root_shelf_id
			WHERE s.id = v.id AND uts.user_id = ? AND uts.permission IN ?;
		`, schemas.RootShelf{}.TableName(), strings.Join(placeholders, ","))

	args = append(args, userId)
	args = append(args, allowedPermissions)

	result := r.db.Raw(sql, args...)
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NoChanges()
	}

	return nil
}

func (r *RootShelfRepository) RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
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

func (r *RootShelfRepository) RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
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

func (r *RootShelfRepository) SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Delete(&schemas.RootShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *RootShelfRepository) HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("\"RootShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Delete(&schemas.RootShelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
