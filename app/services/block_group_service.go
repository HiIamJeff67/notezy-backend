package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	"notezy-backend/app/models/schemas/enums"
	blockgroupsql "notezy-backend/app/models/sql/block_group"
	validation "notezy-backend/app/validation"
	"notezy-backend/shared/constants"
	"notezy-backend/shared/lib/queue"
	"notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type BlockGroupServiceInterface interface {
}

type BlockGroupService struct {
	db                        *gorm.DB
	blockGroupGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository           repositories.BlockRepositoryInterface
}

func NewBlockGroupService(
	db *gorm.DB,
	blockGroupGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockGroupServiceInterface {
	return &BlockGroupService{
		db:                        db,
		blockGroupGroupRepository: blockGroupGroupRepository,
		blockRepository:           blockRepository,
	}
}

/* ============================== Implementations ============================== */

func (s *BlockGroupService) GetMyBlockGroupById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto,
) (*dtos.GetMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	blockGroup, exception := s.blockGroupGroupRepository.GetOneById(
		db,
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlockGroupByIdResDto{
		Id:               blockGroup.Id,
		BlockPackId:      blockGroup.BlockPackId,
		PrevBlockGroupId: blockGroup.PrevBlockGroupId,
		SyncBlockGroupId: blockGroup.SyncBlockGroupId,
		MegaByteSize:     blockGroup.MegaByteSize,
		DeletedAt:        blockGroup.DeletedAt,
		UpdatedAt:        blockGroup.UpdatedAt,
		CreatedAt:        blockGroup.CreatedAt,
	}, nil
}

func (s *BlockGroupService) GetMyBlockGroupAndItsBlocksById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto,
) (*dtos.GetMyBlockGroupAndItsBlocksByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	onlyDeleted := types.Ternary_Negative

	type outputType struct {
		// block group data
		BlockGroupId        uuid.UUID  `gorm:"column:block_group_id;"`
		BlockPackId         uuid.UUID  `gorm:"column:block_pack_id;"`
		PrevBlockGroupId    *uuid.UUID `gorm:"column:prev_block_group_id;"`
		SyncBlockGroupId    *uuid.UUID `gorm:"column:sync_block_group_id;"`
		MegaByteSize        float64    `gorm:"column:mega_byte_size;"`
		BlockGroupDeletedAt *time.Time `gorm:"column:block_group_deleted_at;"`
		BlockGroupUpdatedAt time.Time  `gorm:"column:block_group_updated_at;"`
		BlockGroupCreatedAt time.Time  `gorm:"column:block_group_created_at;"`
		// block data
		BlockId        uuid.UUID       `gorm:"column:block_id;"`
		ParentBlockId  *uuid.UUID      `gorm:"column:parent_block_id;"`
		BlockType      enums.BlockType `gorm:"column:block_type;"`
		BlockProps     datatypes.JSON  `gorm:"column:block_props;"`
		BlockContent   datatypes.JSON  `gorm:"column:block_content;"`
		BlockDeletedAt *time.Time      `gorm:"column:block_deleted_at;"`
		BlockUpdatedAt time.Time       `gorm:"column:block_updated_at;"`
		BlockCreatedAt time.Time       `gorm:"column:block_created_at;"`
	}
	var output []outputType
	result := db.Raw(blockgroupsql.GetMyBlockGroupAndItsBlocksByIdSQL,
		reqDto.Param.BlockGroupId, reqDto.ContextFields.UserId, pg.Array(allowedPermissions), onlyDeleted,
	).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}
	if len(output) == 0 {
		return nil, exceptions.BlockGroup.NotFound()
	}

	var root *dtos.EditableRawBlockContent = nil
	parentToChildrenMap := make(map[uuid.UUID][]outputType, len(output))
	for _, row := range output {
		if row.BlockGroupId != output[0].BlockGroupId {
			return nil, exceptions.BlockGroup.MoreThanOneBlockGroupDetected(output[0].BlockGroupId, row.BlockGroupId)
		}

		if row.ParentBlockId == nil {
			if root != nil {
				// duplicate root block detected
				return nil, exceptions.BlockGroup.RepeatedRootBlockInBlockGroupDetected(output[0].BlockGroupId, row.BlockId)
			}
			root = &dtos.EditableRawBlockContent{
				Id:       row.BlockId,
				Type:     row.BlockType,
				Props:    row.BlockProps,
				Content:  row.BlockContent,
				Children: []dtos.EditableRawBlockContent{},
			}
		} else {
			parentToChildrenMap[*row.ParentBlockId] = append(parentToChildrenMap[*row.ParentBlockId], row)
		}
	}

	if root == nil {
		return nil, exceptions.BlockGroup.NoRootBlockInBlockGroup(output[0].BlockGroupId)
	}

	q := queue.NewQueue[*dtos.EditableRawBlockContent](len(output))
	q.Enqueue(root)
	visited := make(map[uuid.UUID]bool, len(output))
	visited[root.Id] = true
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
		}
		// at this point, current cannot be nil, because the root is not nil, and the below new element enqueued to the queue is alos not nil

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		children, exist := parentToChildrenMap[current.Id]
		if !exist {
			// no children under the current
			continue
		}
		current.Children = make([]dtos.EditableRawBlockContent, 0, len(children))
		for _, child := range children {
			current.Children = append(current.Children, dtos.EditableRawBlockContent{
				Id:       child.BlockId,
				Type:     child.BlockType,
				Props:    child.BlockProps,
				Content:  child.BlockContent,
				Children: []dtos.EditableRawBlockContent{}, // the children of the child should be initialize here
			})
			currentChildPtr := &current.Children[len(current.Children)-1] // get the pointer to the child in current.Children
			q.Enqueue(currentChildPtr)                                    // make sure we passing the pointer of the editable child to the queue, so that we can modify its children field later
		}
	}

	return &dtos.GetMyBlockGroupAndItsBlocksByIdResDto{
		Id:                   output[0].BlockGroupId,
		BlockPackId:          output[0].BlockPackId,
		PrevBlockGroupId:     output[0].PrevBlockGroupId,
		SyncBlockGroupId:     output[0].SyncBlockGroupId,
		DeletedAt:            output[0].BlockGroupDeletedAt,
		UpdatedAt:            output[0].BlockGroupUpdatedAt,
		CreatedAt:            output[0].BlockGroupCreatedAt,
		EditableBlockContent: *root,
	}, nil
}

func (s *BlockGroupService) GetMyBlockGroupsByPrevBlockGroupId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto,
) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var resDto dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("prev_block_group_id = ? AND EXISTS (?) AND deleted_at IS NULL",
			reqDto.Param.PrevBlockGroupId, subQuery,
		).Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetAllMyBlockGroupsByBlockPackId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto,
) (*dtos.GetAllMyBlockGroupsByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var resDto dtos.GetAllMyBlockGroupsByBlockPackIdResDto

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("block_pack_id = ? AND EXISTS (?) AND deleted_at IS NULL",
			reqDto.Param.BlockPackId, subQuery,
		).Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetMyBlocksGroupsAndTheirBlocksByBlockPackId() {}

func (s *BlockGroupService) CreateBlockGroupByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	newBlockGroupId, exception := s.blockGroupGroupRepository.CreateOneByBlockPackId(
		db,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
	)
	if exception != nil {
		return nil, exception
	}
	if newBlockGroupId == nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("got nil block group id")
	}

	return &dtos.CreateBlockGroupByBlockPackIdResDto{
		Id:        *newBlockGroupId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) CreateBlockGroupAndItsBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	newBlockGroupId, exception := s.blockGroupGroupRepository.CreateOneByBlockPackId(
		tx,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if newBlockGroupId == nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("got nil block group id")
	}

	input := []inputs.CreateBlockInput{}
	rootProps, err := json.Marshal(reqDto.Body.EditableBlockContent.Props)
	if err != nil {
		tx.Rollback()
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}
	rootContent, err := json.Marshal(reqDto.Body.EditableBlockContent.Content)
	if err != nil {
		tx.Rollback()
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}
	input = append(input, inputs.CreateBlockInput{
		PrevBlockId: nil,
		Type:        reqDto.Body.EditableBlockContent.Type,
		Props:       datatypes.JSON(rootProps),
		Content:     datatypes.JSON(rootContent),
	})
	// the capacity is just a pointer, so it is no performance issue to intialize a queue with capacity of max integer
	q := queue.NewQueue[*dtos.EditableBlockContent](constants.MAX_INT)
	q.Enqueue(&reqDto.Body.EditableBlockContent) // enqueue the root block which is just under the reqDto.Body along with the block group data
	visited := make(map[uuid.UUID]bool)
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			tx.Rollback()
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for _, child := range current.Children {
			props, err := json.Marshal(child.Props)
			if err != nil {
				tx.Rollback()
				return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
			}
			content, err := json.Marshal(child.Content)
			if err != nil {
				tx.Rollback()
				return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
			}

			input = append(input, inputs.CreateBlockInput{
				Id:          child.Id,
				PrevBlockId: &current.Id,
				Type:        child.Type,
				Props:       props,
				Content:     content,
			})

			q.Enqueue(&child)
		}
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupId(
		tx,
		*newBlockGroupId,
		reqDto.ContextFields.UserId,
		constants.MaxBatchCreateBlockSize,
		input,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) CreateBlockGroupsAndTheirBlocksByBlockPackId() {}

func (s *BlockGroupService) MoveMyBlockGroupById() {}

func (s *BlockGroupService) MoveMyBlockGroupsByIds() {}

func (s *BlockGroupService) RestoreMyBlockGroupById() {}

func (s *BlockGroupService) RestoreMyBlockGroupsByIds() {}

func (s *BlockGroupService) DeleteMyBlockGroupById() {}

func (s *BlockGroupService) DeleteMyBlockGroupsByIds() {}
