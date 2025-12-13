package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ============================== Definitions ============================== */

type BlockPackRepositoryInterface interface {
	HasPermission(db *gorm.DB, id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	HasPermissions(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	CheckPermissionAndGetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]schemas.BlockPack, *exceptions.Exception)
	CheckPermissionAndGetOneWithOwnerIdById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*uuid.UUID, *schemas.BlockPack, *exceptions.Exception)
	CheckPermissionsAndGetManyWithOwnerIdsByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockPackRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]uuid.UUID, []schemas.BlockPack, *exceptions.Exception)
	GetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) (*schemas.BlockPack, *exceptions.Exception)
	CreateOneBySubShelfId(db *gorm.DB, subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockPackInput, skipPermissionCheck bool) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockPackInput) (*schemas.BlockPack, *exceptions.Exception)
	RestoreSoftDeletedOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
	RestoreSoftDeletedManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
	SoftDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
	SoftDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
	HardDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
	HardDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, skipPermissionCheck bool) *exceptions.Exception
}

type BlockPackRepository struct{}

func NewBlockPackRepository() BlockPackRepositoryInterface {
	return &BlockPackRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockPackRepository) HasPermission(
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
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("\"BlockPackTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockPackRepository) HasPermissions(
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
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id == ss.id").
		Where("\"BlockPackTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockPackRepository) CheckPermissionAndGetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*schemas.BlockPack, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPack schemas.BlockPack
	result := query.First(&blockPack)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &blockPack, nil
}

func (r *BlockPackRepository) CheckPermissionsAndGetManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]schemas.BlockPack, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPacks []schemas.BlockPack
	result := query.Find(&blockPacks)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}
	if len(blockPacks) == 0 {
		return nil, exceptions.BlockPack.NotFound()
	}

	return blockPacks, nil
}

func (r *BlockPackRepository) CheckPermissionAndGetOneWithOwnerIdById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*uuid.UUID, *schemas.BlockPack, *exceptions.Exception) { // we should also return the owner id for the block groups and blocks
	if db == nil {
		db = models.NotezyDB
	}

	// note that the subQuery is querying the permission of the current user,
	// 			 and the query is querying the data and the owner id(which may be different from the current user)
	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Select("\"BlockPackTable\".*, owner_uts.user_id AS owner_id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Joins("INNER JOIN \"UsersToShelvesTable\" owner_uts ON ss.root_shelf_id = owner_uts.root_shelf_id AND owner_uts.permission = 'Owner'").
		Where("\"BlockPackTable\".id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPackWithOwnerId struct {
		schemas.BlockPack
		OwnerId uuid.UUID `gorm:"column:owner_id;"`
	}
	result := query.First(&blockPackWithOwnerId)
	if err := result.Error; err != nil {
		return nil, nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &blockPackWithOwnerId.OwnerId, &blockPackWithOwnerId.BlockPack, nil
}

func (r *BlockPackRepository) CheckPermissionsAndGetManyWithOwnerIdsByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockPackRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]uuid.UUID, []schemas.BlockPack, *exceptions.Exception) { // we should also return the owner id for the block groups and blocks
	if db == nil {
		db = models.NotezyDB
	}

	// note that the subQuery is querying the permission of the current user,
	// 			 and the query is querying the data and the owner id(which may be different from the current user)
	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = ss.root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockPack{}).
		Select("\"BlockPackTable\".*, owner_uts.user_id AS owner_id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Joins("INNER JOIN \"UsersToShelvesTable\" owner_uts ON ss.root_shelf_id = owner_uts.root_shelf_id AND owner_uts.permission = 'Owner'").
		Where("\"BlockPackTable\".id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("\"BlockPackTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var blockPacksWithOwnerIds []struct {
		schemas.BlockPack
		ownerId uuid.UUID `gorm:"column:owner_id;"`
	}
	result := query.Find(&blockPacksWithOwnerIds)
	if err := result.Error; err != nil {
		return nil, nil, exceptions.BlockPack.NotFound().WithError(err)
	}
	if len(blockPacksWithOwnerIds) == 0 {
		return nil, nil, exceptions.BlockPack.NotFound()
	}

	ownerIds := make([]uuid.UUID, len(blockPacksWithOwnerIds))
	blockPacks := make([]schemas.BlockPack, len(blockPacksWithOwnerIds))
	for index, element := range blockPacksWithOwnerIds {
		ownerIds[index] = element.ownerId
		blockPacks[index] = element.BlockPack
	}

	return ownerIds, blockPacks, nil
}

func (r *BlockPackRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
) (*schemas.BlockPack, *exceptions.Exception) {
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
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
}

func (r *BlockPackRepository) CreateOneBySubShelfId(
	db *gorm.DB,
	subShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockPackInput,
	skipPermissionCheck bool,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		subShelfRepository := NewSubShelfRepository()

		if !subShelfRepository.HasPermission(
			db,
			subShelfId,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return nil, exceptions.Shelf.NoPermission("create a block pack under this shelf")
		}
	}

	var newBlockPack schemas.BlockPack
	if err := copier.Copy(&newBlockPack, &input); err != nil {
		return nil, exceptions.BlockPack.FailedToCreate().WithError(err)
	}
	newBlockPack.ParentSubShelfId = subShelfId

	result := db.Model(&schemas.BlockPack{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockPack)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToCreate().WithError(err)
	}

	return &newBlockPack.Id, nil
}

func (r *BlockPackRepository) UpdateOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockPackInput,
) (*schemas.BlockPack, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingBlockPack, exception := r.CheckPermissionAndGetOneById(
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
	if existingBlockPack == nil {
		return nil, exceptions.BlockPack.NotFound()
	}

	if input.Values.ParentSubShelfId != nil && (input.SetNull == nil || !(*input.SetNull)["ParentSubShelfId"]) {
		subShelfRepository := NewSubShelfRepository()
		if !subShelfRepository.HasPermission(
			db,
			*input.Values.ParentSubShelfId,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return nil, exceptions.Shelf.NoPermission("move a block pack to this shelf")
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlockPack)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlockPack,
		).WithError(err)
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockPack.NoChanges()
	}

	return &updates, nil
}

func (r *BlockPackRepository) RestoreSoftDeletedOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermission(
			db,
			id,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("restore a deleted block pack")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) RestoreSoftDeletedManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermissions(
			db,
			ids,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("restore deleted block packs")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) SoftDeleteOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermission(
			db,
			id,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("soft delete a block pack")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) SoftDeleteManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermissions(
			db,
			ids,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("soft delete block packs")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) HardDeleteOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermission(
			db,
			id,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("hard delete a block pack")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.BlockPack{})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}

func (r *BlockPackRepository) HardDeleteManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	skipPermissionCheck bool,
) *exceptions.Exception {
	if db == nil {
		db = models.NotezyDB
	}

	if !skipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		}

		if !r.HasPermissions(
			db,
			ids,
			userId,
			allowedPermissions,
			types.Ternary_Negative,
		) {
			return exceptions.BlockPack.NoPermission("hard delete block packs")
		}
	}

	result := db.Model(&schemas.BlockPack{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.BlockPack{})
	if err := result.Error; err != nil {
		return exceptions.BlockPack.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockPack.NoChanges()
	}

	return nil
}
