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

type BlockRepositoryInterface interface {
	HasPermission(db *gorm.DB, id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	HasPermissions(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) bool
	CheckPermissionAndGetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) (*schemas.Block, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, allowedPermissions []enums.AccessControlPermission, onlyDeleted types.Ternary) ([]schemas.Block, *exceptions.Exception)
	GetOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, preloads []schemas.BlockRelation, onlyDeleted types.Ternary) (*schemas.Block, *exceptions.Exception)
	CreateOneByBlockGroupId(db *gorm.DB, blockGroupId uuid.UUID, userId uuid.UUID, input inputs.CreateBlockInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateBlockInput) (*schemas.Block, *exceptions.Exception)
	RestoreSoftDeletedOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	RestoreSoftDeletedManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	SoftDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteOneById(db *gorm.DB, id uuid.UUID, userId uuid.UUID) *exceptions.Exception
	HardDeleteManyByIds(db *gorm.DB, ids []uuid.UUID, userId uuid.UUID) *exceptions.Exception
}

type BlockRepository struct{}

func NewBlockRepository() BlockRepositoryInterface {
	return &BlockRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockRepository) HasPermission(
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
	query := db.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
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

func (r *BlockRepository) HasPermissions(
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
	query := db.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
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

func (r *BlockRepository) CheckPermissionAndGetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) (*schemas.Block, *exceptions.Exception) {
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
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
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

	var block schemas.Block
	result := query.First(&block)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound().WithError(err)
	}

	return &block, nil
}

func (r *BlockRepository) CheckPermissionsAndGetManyByIds(
	db *gorm.DB,
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	allowedPermissions []enums.AccessControlPermission,
	onlyDeleted types.Ternary,
) ([]schemas.Block, *exceptions.Exception) {
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
		Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
		Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
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

	var blocks []schemas.Block
	result := query.First(&blocks)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound().WithError(err)
	}

	return blocks, nil
}

func (r *BlockRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockRelation,
	onlyDeleted types.Ternary,
) (*schemas.Block, *exceptions.Exception) {
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

func (r *BlockRepository) CreateOneByBlockGroupId(
	db *gorm.DB,
	blockGroupId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockGroupRepository := NewBlockGroupRepository()

	if !blockGroupRepository.HasPermission(
		db,
		blockGroupId,
		userId,
		allowedPermissions,
		types.Ternary_Negative,
	) {
		return nil, exceptions.Block.NoPermission("get owner's block group")
	}

	var newBlock schemas.Block
	if err := copier.Copy(&newBlock, &input); err != nil {
		return nil, exceptions.Block.FailedToCreate().WithError(err)
	}
	newBlock.BlockGroupId = blockGroupId

	result := db.Model(&schemas.Block{}).
		Create(&newBlock)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToCreate().WithError(err)
	}

	return &newBlock.Id, nil
}

func (r *BlockRepository) UpdateOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateBlockInput,
) (*schemas.Block, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	// maybe we need a more efficient way to update the field of blocks
	// since they will be used quite frequently

	existingBlock, exception := r.CheckPermissionAndGetOneById(
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

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingBlock)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingBlock,
		).WithError(err)
	}

	result := db.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Block.NoChanges()
	}

	return &updates, nil
}

func (r *BlockRepository) RestoreSoftDeletedOneById(
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
		return exceptions.Block.NoPermission("restore a deleted block")
	}

	result := db.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) RestoreSoftDeletedManyByIds(
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
		return exceptions.Block.NoPermission("restore deleted blocks")
	}

	result := db.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) SoftDeleteOneById(
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
		return exceptions.Block.NoPermission("soft delete a block")
	}

	result := db.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) SoftDeleteManyByIds(
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
		return exceptions.Block.NoPermission("soft delete blocks")
	}

	result := db.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) HardDeleteOneById(
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
		return exceptions.BlockGroup.NoPermission("hard delete a block")
	}

	result := db.Model(&schemas.Block{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.Block{})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}

func (r *BlockRepository) HardDeleteManyByIds(
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
		return exceptions.Block.NoPermission("hard delete blocks")
	}

	result := db.Model(&schemas.Block{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.Block{})
	if err := result.Error; err != nil {
		return exceptions.Block.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Block.NoChanges()
	}

	return nil
}
