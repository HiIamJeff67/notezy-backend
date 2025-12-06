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

type BlockGroupRepositoryInterface interface {
	HasPermission(db *gorm.DB, id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	HasPermissions(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	CheckPermissionAndGetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*schemas.BlockGroup, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]schemas.BlockGroup, *exceptions.Exception)
	GetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockGroupRelation) (*schemas.BlockGroup, *exceptions.Exception)
	CreateOneByBlockPackId(db *gorm.DB, blockPackId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockGroupInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockGroupInput) (*schemas.BlockGroup, *exceptions.Exception)
	RestoreSoftDeletedOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type BlockGroupRepository struct{}

func NewBlockGroupRepository() BlockGroupRepositoryInterface {
	return &BlockGroupRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockGroupRepository) HasPermission(
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
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockGroupRepository) HasPermissions(
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
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("id IN ? AND EXISTS (?)",
			ids, subQuery,
		)

	switch onlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *BlockGroupRepository) CheckPermissionAndGetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*schemas.BlockGroup, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
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

	var blockGroup schemas.BlockGroup
	result := query.First(&blockGroup)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &blockGroup, nil
}

func (r *BlockGroupRepository) CheckPermissionsAndGetManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]schemas.BlockGroup, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
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

	var blockGroups []schemas.BlockGroup
	result := query.Find(&blockGroups)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}
	if len(blockGroups) == 0 {
		return nil, exceptions.BlockGroup.NotFound()
	}

	return blockGroups, nil
}

func (r *BlockGroupRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
) (*schemas.BlockGroup, *exceptions.Exception) {
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

func (r *BlockGroupRepository) CreateOneByBlockPackId(
	db *gorm.DB,
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockGroupInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackRepository := NewBlockPackRepository()

	ownerId, blockPack, exception := blockPackRepository.CheckPermissionAndGetOneWithOwnerIdById(
		db,
		blockPackId,
		userId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
	if exception != nil {
		return nil, exception
	}
	if ownerId == nil || blockPack == nil {
		return nil, exceptions.BlockPack.NoPermission("get owner's block pack")
	}

	var newBlockGroup schemas.BlockGroup
	if err := copier.Copy(&newBlockGroup, &input); err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}
	newBlockGroup.OwnerId = *ownerId // get the owner id from the CheckPermissionAndGetOneById
	newBlockGroup.BlockPackId = blockPackId

	result := db.Model(&schemas.BlockGroup{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newBlockGroup)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}

	return &newBlockGroup.Id, nil
}

func (r *BlockGroupRepository) UpdateOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockGroupInput,
) (*schemas.BlockGroup, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingBlockGroup, exception := r.CheckPermissionAndGetOneById(
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlockGroup)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlockGroup,
		).WithError(err)
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	return &updates, nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedOneById(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("restore a deleted block group")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) RestoreSoftDeletedManyByIds(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("restore deleted block groups")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) SoftDeleteOneById(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("soft delete a block group")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) SoftDeleteManyByIds(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("soft delete block groups")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) HardDeleteOneById(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermission(
		db,
		id,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("hard delete a block group")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.BlockGroup{})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}

func (r *BlockGroupRepository) HardDeleteManyByIds(
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
		enums.AccessControlPermission_Write,
	}

	if !r.HasPermissions(
		db,
		ids,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return exceptions.BlockGroup.NoPermission("hard delete block groups")
	}

	result := db.Model(&schemas.BlockGroup{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.BlockGroup{})
	if err := result.Error; err != nil {
		return exceptions.BlockGroup.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	return nil
}
