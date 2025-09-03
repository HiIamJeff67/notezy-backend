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

type ShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission) bool
	GetOneById(id uuid.UUID, ownerId uuid.UUID, preloads *[]schemas.ShelfRelation) (*schemas.Shelf, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) (*schemas.Shelf, *exceptions.Exception)
	DirectlyUpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) *exceptions.Exception
	DirectlyUpdateManyByIds(ids []uuid.UUID, ownerId uuid.UUID, inputs []inputs.PartialUpdateShelfInput) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type ShelfRepository struct {
	db *gorm.DB
}

func NewShelfRepository(db *gorm.DB) ShelfRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &ShelfRepository{db: db}
}

/* ============================== Helper functions ============================== */

func (r *ShelfRepository) getPartialUpdatePlaceholderUnit() string {
	return "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
}

func (r *ShelfRepository) getNumOfPartialUpdateArguments(numOfRows int) int {
	// include values.ShelfId, values.ShelfName, values.ShelfEncodedStructure, setNull.ShelfName, setNull.ShelfEncodedStructure
	// and reserved 1 for the ownerId.
	// Note that if its for batch updates, the ownerIds could be more than one, hence this doesn't work
	if numOfRows > (constants.MaxDatabaseUpdateParameters-1)/5 {
		return 0
	}
	return numOfRows*5 + 1
}

/* ============================== CRUD operations ============================== */

func (r *ShelfRepository) HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermission []enums.AccessControlPermission) bool {
	var count int64 = 0
	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermission).
		Count(&count)
	if err := result.Error; err != nil || count == 0 {
		return false
	}

	return true
}

func (r *ShelfRepository) GetOneById(id uuid.UUID, ownerId uuid.UUID, preloads *[]schemas.ShelfRelation) (*schemas.Shelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	shelf := schemas.Shelf{}
	db := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id")
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, ownerId, allowedPermissions).
		First(&shelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &shelf, nil
}

func (r *ShelfRepository) CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateShelfInput) (*uuid.UUID, *exceptions.Exception) {
	var newShelf schemas.Shelf
	newShelf.OwnerId = ownerId
	if err := copier.Copy(&newShelf, &input); err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	result := r.db.Model(&schemas.Shelf{}).
		Create(&newShelf)
	if err := result.Error; err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"shelf_idx_owner_id_name\" (SQLSTATE 23505)":
			return nil, exceptions.Shelf.DuplicateName(input.Name)
		default:
			return nil, exceptions.Shelf.FailedToCreate().WithError(err)
		}
	}
	if newShelf.Id != input.Id {
		return nil, exceptions.Shelf.FailedToCreate().WithDetails("Create different id shelf")
	}

	// create the users to shelves relation with the permission of admin
	newUsersToShelves := schemas.UsersToShelves{
		UserId:     ownerId,
		ShelfId:    input.Id,
		Permission: enums.AccessControlPermission_Admin,
	}
	result = r.db.Model(&schemas.UsersToShelves{}).
		Create(&newUsersToShelves)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	return &newShelf.Id, nil
}

func (r *ShelfRepository) UpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) (*schemas.Shelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	existingShelf, exception := r.GetOneById(id, ownerId, nil)
	if exception != nil || existingShelf == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingShelf)
	}

	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, ownerId, allowedPermissions).
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

func (r *ShelfRepository) DirectlyUpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) *exceptions.Exception {
	return r.DirectlyUpdateManyByIds([]uuid.UUID{id}, ownerId, []inputs.PartialUpdateShelfInput{input})
}

func (r *ShelfRepository) DirectlyUpdateManyByIds(ids []uuid.UUID, ownerId uuid.UUID, inputs []inputs.PartialUpdateShelfInput) *exceptions.Exception {
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
		args = append(args, util.DerefOrNil(inputs[index].Values.EncodedStructure))
		args = append(args, util.DerefOrNil(inputs[index].Values.EncodedStructureByteSize))
		args = append(args, util.DerefOrNil(inputs[index].Values.TotalShelfNodes))
		args = append(args, util.DerefOrNil(inputs[index].Values.TotalMaterials))
		args = append(args, util.DerefOrNil(inputs[index].Values.MaxWidth))
		args = append(args, util.DerefOrNil(inputs[index].Values.MaxDepth))

		// safetly dereference all the setNulls
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "Name"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "EncodedStructure"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "EncodedStructureByteSize"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "TotalShelfNodes"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "TotalMaterials"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "MaxWidth"))
		args = append(args, util.CheckSetNull(inputs[index].SetNull, "MaxDepth"))
	}

	sql := fmt.Sprintf(`
			UPDATE "%s" AS s
			SET
				name = CASE
					WHEN v.set_null_name THEN NULL
					ELSE COALESCE(v.name, s.name)
				END, 
				encoded_structure = CASE
					WHEN v.set_null_encoded_structure THEN NULL
					ELSE COALESCE(v.encoded_structure, s.encoded_structure)
				END,
				encoded_structure_byte_size = CASE
					WHEN v.encoded_structure_byte_size THEN NULL
					ELSE COALESCE(v.encoded_structure_byte_size, s.encoded_structure_byte_size)
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
			FROM (VALUES %s) AS v(id, name, encoded_structure, set_null_name, set_null_encoded_structure)
			LEFT JOIN "UsersToShelvesTable" AS uts ON s.id = uts.shelf_id
			WHERE s.id = v.id AND uts.user_id = ? AND uts.permission IN ?;
		`, schemas.Shelf{}.TableName(), strings.Join(placeholders, ","))

	args = append(args, ownerId)
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

func (r *ShelfRepository) RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *ShelfRepository) RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *ShelfRepository) SoftDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	now := time.Now()
	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Update("deleted_at", now)
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *ShelfRepository) SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	now := time.Now()
	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Update("deleted_at", now)
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *ShelfRepository) HardDeleteOneById(id uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id = ? AND uts.user_id = ? AND uts.permission IN ?", id, userId, allowedPermissions).
		Delete(&schemas.Shelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}

func (r *ShelfRepository) HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Admin,
	}

	result := r.db.Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("\"ShelfTable\".id IN ? AND uts.user_id = ? AND uts.permission IN ?", ids, userId, allowedPermissions).
		Delete(&schemas.Shelf{})
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
