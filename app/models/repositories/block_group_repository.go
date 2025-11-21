package repositories

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	"notezy-backend/app/models/schemas/enums"
	"notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ============================== Definitions ============================== */

type BlockGroupRepositoryInterface interface {
}

type BlockGroupRepository struct{}

func NewBlockGroupRepository() BlockGroupRepositoryInterface {
	return &BlockGroupRepository{}
}

/* ============================== Implementations ============================== */

func (r *BlockGroupRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.BlockGroupRelation,
	onlyDeleted types.Ternary,
) (*schemas.BlockGroup, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	query := db.Model(&schemas.BlockGroup{}).
		Where("id = ? AND owner_id = ?",
			id, userId,
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

func (r *BlockGroupRepository) CreateOneByBlockPackId(
	db *gorm.DB,
	blockPackId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateBlockGroupInput,
) (*schemas.BlockGroup, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	allowedPermissions := []enums.AccessControlPermission{
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
		Create(&newBlockGroup)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithError(err)
	}

	return &newBlockGroup, nil
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

	existingBlockGroup, exception := r.GetOneById(
		db,
		id,
		userId,
		nil,
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
