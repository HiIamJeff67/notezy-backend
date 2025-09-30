package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type SubShelfServiceInterface interface {
	GetMySubShelfById(reqDto *dtos.GetMySubShelfByIdReqDto) (*dtos.GetMySubShelfByIdResDto, *exceptions.Exception)
	GetMySubShelvesByPrevSubShelfId(reqDto *dtos.GetMySubShelvesByPrevSubShelfIdReqDto) (*dtos.GetMySubShelvesByPrevSubShelfIdResDto, *exceptions.Exception)
	GetAllMySubShelvesByRootShelfId(reqDto *dtos.GetAllMySubShelvesByRootShelfIdReqDto) (*dtos.GetAllMySubShelvesByRootShelfIdResDto, *exceptions.Exception)
	CreateSubShelfByRootShelfId(reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) (*dtos.CreateSubShelfByRootShelfIdResDto, *exceptions.Exception)
	UpdateMySubShelfById(reqDto *dtos.UpdateMySubShelfByIdReqDto) (*dtos.UpdateMySubShelfByIdResDto, *exceptions.Exception)
	MoveMySubShelf(reqDto *dtos.MoveMySubShelfReqDto) (*dtos.MoveMySubShelfResDto, *exceptions.Exception)
	MoveMySubShelves(reqDto *dtos.MoveMySubShelvesReqDto) (*dtos.MoveMySubShelvesResDto, *exceptions.Exception)
	RestoreMySubShelfById(reqDto *dtos.RestoreMySubShelfByIdReqDto) (*dtos.RestoreMySubShelfByIdResDto, *exceptions.Exception)
	RestoreMySubShelvesByIds(reqDto *dtos.RestoreMySubShelvesByIdsReqDto) (*dtos.RestoreMySubShelvesByIdsResDto, *exceptions.Exception)
	DeleteMySubShelfById(reqDto *dtos.DeleteMySubShelfByIdReqDto) (*dtos.DeleteMySubShelfByIdResDto, *exceptions.Exception)
	DeleteMySubShelvesByIds(reqDto *dtos.DeleteMySubShelvesByIdsReqDto) (*dtos.DeleteMySubShelvesByIdsResDto, *exceptions.Exception)
}

type SubShelfService struct {
	db *gorm.DB
}

func NewSubShelfService(db *gorm.DB) SubShelfServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &SubShelfService{db: db}
}

/* ============================== Service Methods for SubShelf ============================== */

func (s *SubShelfService) GetMySubShelfById(reqDto *dtos.GetMySubShelfByIdReqDto) (
	*dtos.GetMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	subShelf, exception := subShelfRepository.GetOneById(
		reqDto.Param.SubShelfId,
		reqDto.ContextFields.UserId,
		nil,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMySubShelfByIdResDto{
		Id:             subShelf.Id,
		Name:           subShelf.Name,
		RootShelfId:    subShelf.RootShelfId,
		PrevSubShelfId: subShelf.PrevSubShelfId,
		Path:           subShelf.Path,
		DeletedAt:      subShelf.DeletedAt,
		UpdatedAt:      subShelf.UpdatedAt,
		CreatedAt:      subShelf.CreatedAt,
	}, nil
}

func (s *SubShelfService) GetMySubShelvesByPrevSubShelfId(reqDto *dtos.GetMySubShelvesByPrevSubShelfIdReqDto) (
	*dtos.GetMySubShelvesByPrevSubShelfIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	resDto := dtos.GetMySubShelvesByPrevSubShelfIdResDto{}

	subQuery := s.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := s.db.Model(&schemas.SubShelf{}).
		Where("prev_sub_shelf_id = ? AND EXISTS (?)", reqDto.Param.PrevSubShelfId, subQuery).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *SubShelfService) GetAllMySubShelvesByRootShelfId(reqDto *dtos.GetAllMySubShelvesByRootShelfIdReqDto) (
	*dtos.GetAllMySubShelvesByRootShelfIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	resDto := dtos.GetAllMySubShelvesByRootShelfIdResDto{}

	subQuery := s.db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := s.db.Model(&schemas.SubShelf{}).
		Where("root_shelf_id = ? AND EXISTS (?) AND deleted_at IS NULL", reqDto.Param.RootShelfId, subQuery).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *SubShelfService) CreateSubShelfByRootShelfId(reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) (
	*dtos.CreateSubShelfByRootShelfIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	subShelfId, exception := subShelfRepository.CreateOneByUserId(
		reqDto.ContextFields.UserId,
		inputs.CreateSubShelfInput{
			Name:           reqDto.Body.Name,
			RootShelfId:    reqDto.Body.RootShelfId,
			PrevSubShelfId: reqDto.Body.PrevSubShelfId,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateSubShelfByRootShelfIdResDto{
		Id:        *subShelfId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) UpdateMySubShelfById(reqDto *dtos.UpdateMySubShelfByIdReqDto) (
	*dtos.UpdateMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	subShelf, exception := subShelfRepository.UpdateOneById(
		reqDto.Body.SubShelfId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateSubShelfInput{
			Values: inputs.UpdateSubShelfInput{
				Name: reqDto.Body.Values.Name,
			},
			SetNull: reqDto.Body.SetNull,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMySubShelfByIdResDto{
		UpdatedAt: subShelf.UpdatedAt,
	}, nil
}

func (s *SubShelfService) MoveMySubShelf(reqDto *dtos.MoveMySubShelfReqDto) (
	*dtos.MoveMySubShelfResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	if reqDto.Body.DestinationSubShelfId != nil &&
		reqDto.Body.SourceSubShelfId == *reqDto.Body.DestinationSubShelfId {
		return nil, exceptions.Shelf.NoChanges()
	}

	tx := s.db.Begin()
	subShelfRepository := repositories.NewSubShelfRepository(tx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	from, exception := subShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Body.SourceSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if from.RootShelfId != reqDto.Body.SourceRootShelfId {
		tx.Rollback()
		return nil, exceptions.Shelf.NotFound()
	}

	if reqDto.Body.DestinationSubShelfId != nil {
		to, exception := subShelfRepository.CheckPermissionAndGetOneById(
			*reqDto.Body.DestinationSubShelfId,
			reqDto.ContextFields.UserId,
			nil,
			allowedPermissions,
			types.Ternary_Negative,
		)
		if exception != nil {
			return nil, exception
		}
		if to.RootShelfId != reqDto.Body.DestinationRootShelfId {
			return nil, exceptions.Shelf.NotFound()
		}

		if len(from.Path)+len(to.Path) > int(constants.MaxSubShelvesOfRootShelf) {
			return nil, exceptions.Shelf.MaximumDepthExceeded(
				int32(len(from.Path)+len(to.Path)),
				constants.MaxSubShelvesOfRootShelf,
			)
		}

		// check if to.Path contain any from.Id, if it's true, then it means the user is trying to move the sub shelf to its child
		for _, parent := range to.Path {
			if parent == reqDto.Body.SourceSubShelfId {
				return nil, exceptions.Shelf.InsertParentIntoItsChildren(
					reqDto.Body.DestinationSubShelfId,
					reqDto.Body.SourceSubShelfId,
				)
			}
		}

		to.Path = append(to.Path, to.Id)
		result := tx.Exec(`
			UPDATE "SubShelfTable" 
			SET "root_shelf_id" = ?, "prev_sub_shelf_id" = ?, "path" = ?, "updated_at" = NOW() 
			WHERE id = ? AND deleted_at IS NULL`,
			reqDto.Body.DestinationRootShelfId, reqDto.Body.DestinationSubShelfId, pg.Array(to.Path),
			reqDto.Body.SourceSubShelfId,
		)
		if err := result.Error; err != nil {
			return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
		}
	} else {
		result := tx.Exec(`
			UPDATE "SubShelfTable" 
			SET "root_shelf_id" = ?, "prev_sub_shelf_id" = ?, "path" = ?, "updated_at" = NOW() 
			WHERE id = ? AND deleted_at IS NULL`,
			reqDto.Body.DestinationRootShelfId, nil, pg.Array([]uuid.UUID{}), reqDto.Body.SourceSubShelfId,
		)
		if err := result.Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
		}
	}

	// cascading update child sub shelves
	if reqDto.Body.SourceRootShelfId != reqDto.Body.DestinationRootShelfId {
		var childSubShelves []struct {
			Id   uuid.UUID       `json:"id"`
			Path types.UUIDArray `json:"path"`
		}

		err := tx.Model(&schemas.SubShelf{}).
			Select("id, path").
			Where("root_shelf_id = ? AND path @> ARRAY[?]::uuid[] AND deleted_at IS NULL",
				reqDto.Body.SourceRootShelfId,
				reqDto.Body.SourceSubShelfId).
			Find(&childSubShelves).Error

		if err != nil {
			tx.Rollback()
			return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
		}

		if len(childSubShelves) > 0 {
			var caseStatements []string
			var args []interface{}
			var ids []uuid.UUID

			for _, child := range childSubShelves {
				sourceIndex := -1
				for i, pathItem := range child.Path {
					if pathItem == reqDto.Body.SourceSubShelfId {
						sourceIndex = i
						break
					}
				}

				if sourceIndex >= 0 {
					// 新的 path 應該是：[SourceSubShelfId] + 原本在 SourceSubShelfId 之後的部分
					var newPath types.UUIDArray
					newPath = append(newPath, reqDto.Body.SourceSubShelfId)
					if sourceIndex+1 < len(child.Path) {
						newPath = append(newPath, child.Path[sourceIndex+1:]...)
					}

					caseStatements = append(caseStatements, "WHEN id = ? THEN ?::uuid[]")
					args = append(args, child.Id, pg.Array(newPath))
					ids = append(ids, child.Id)
				}
			}

			if len(caseStatements) > 0 {
				query := fmt.Sprintf(`
                    UPDATE "SubShelfTable" 
                    SET "root_shelf_id" = ?, 
                        "path" = CASE %s END,
                        "updated_at" = NOW() 
                    WHERE id IN ? AND deleted_at IS NULL`,
					strings.Join(caseStatements, " "))

				queryArgs := []interface{}{reqDto.Body.DestinationRootShelfId}
				queryArgs = append(queryArgs, args...)
				queryArgs = append(queryArgs, ids)

				result := tx.Exec(query, queryArgs...)
				if err := result.Error; err != nil {
					tx.Rollback()
					return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Shelf.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.MoveMySubShelfResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) MoveMySubShelves(reqDto *dtos.MoveMySubShelvesReqDto) (
	*dtos.MoveMySubShelvesResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	froms, exception := subShelfRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Body.SourceSubShelfIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		types.Ternary_Negative,
	)
	if exception != nil {
		return nil, exception
	}
	for _, from := range froms {
		if from.RootShelfId != reqDto.Body.SourceRootShelfId {
			return nil, exceptions.Shelf.NotFound()
		}
	}

	if reqDto.Body.DestinationSubShelfId != nil {
		to, exception := subShelfRepository.CheckPermissionAndGetOneById(
			*reqDto.Body.DestinationSubShelfId,
			reqDto.ContextFields.UserId,
			nil,
			allowedPermissions,
			types.Ternary_Negative,
		)
		if exception != nil {
			return nil, exception
		}
		if to.RootShelfId != reqDto.Body.DestinationRootShelfId {
			return nil, exceptions.Shelf.NotFound()
		}
		if to.Path == nil {
			to.Path = []uuid.UUID{}
		}

		sourceSubShelfIdMap := make(map[uuid.UUID]bool, 0)
		for _, from := range froms {
			if len(from.Path)+len(to.Path) > int(constants.MaxSubShelvesOfRootShelf) {
				exceptions.Shelf.MaximumDepthExceeded(
					int32(len(from.Path)+len(to.Path)),
					constants.MaxSubShelvesOfRootShelf,
				).Log()
				// sourceSubShelfIdMap[from.Id] = false
			} else if from.Id == to.Id { // handling inserting node to itself here
				exceptions.Shelf.InsertParentIntoItsChildren(to.Id, from.Id).Log()
				// sourceSubShelfIdMap[from.Id] = false
			} else {
				sourceSubShelfIdMap[from.Id] = true
			}
		}

		for _, parent := range to.Path { // handling inserting node to its parent here
			if sourceSubShelfIdMap[parent] {
				exceptions.Shelf.InsertParentIntoItsChildren(
					reqDto.Body.DestinationSubShelfId,
					parent,
				).Log()
				sourceSubShelfIdMap[parent] = false // has to invalid the sub shelf
			}
		}

		validSourceSubShelfIds := []uuid.UUID{}
		for sourceSubShelfId, exist := range sourceSubShelfIdMap {
			if exist {
				validSourceSubShelfIds = append(validSourceSubShelfIds, sourceSubShelfId)
			}
		}

		to.Path = append(to.Path, to.Id)
		result := s.db.Exec(`
			UPDATE "SubShelfTable" 
			SET "root_shelf_id" = ?, "prev_sub_shelf_id" = ?, "path" = ?, "updated_at" = NOW() 
			WHERE id IN ? AND deleted_at IS NULL`,
			reqDto.Body.DestinationRootShelfId, reqDto.Body.DestinationSubShelfId, pg.Array(to.Path), validSourceSubShelfIds,
		)
		if err := result.Error; err != nil {
			return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
		}
	} else {
		validSourceSubShelfIds := []uuid.UUID{}
		for _, from := range froms {
			validSourceSubShelfIds = append(validSourceSubShelfIds, from.Id)
		}

		result := s.db.Exec(`
			UPDATE "SubShelfTable" 
			SET "root_shelf_id" = ?, "prev_sub_shelf_id" = ?, "path" = ?, "updated_at" = NOW() 
			WHERE id IN ? AND deleted_at IS NULL`,
			reqDto.Body.DestinationRootShelfId, nil, pg.Array([]uuid.UUID{}), validSourceSubShelfIds,
		)
		if err := result.Error; err != nil {
			return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
		}
	}

	return &dtos.MoveMySubShelvesResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) RestoreMySubShelfById(reqDto *dtos.RestoreMySubShelfByIdReqDto) (
	*dtos.RestoreMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.RestoreSoftDeletedOneById(reqDto.Body.SubShelfId, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMySubShelfByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) RestoreMySubShelvesByIds(reqDto *dtos.RestoreMySubShelvesByIdsReqDto) (
	*dtos.RestoreMySubShelvesByIdsResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.RestoreSoftDeletedManyByIds(reqDto.Body.SubShelfIds, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMySubShelvesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) DeleteMySubShelfById(reqDto *dtos.DeleteMySubShelfByIdReqDto) (
	*dtos.DeleteMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.SoftDeleteOneById(reqDto.Body.SubShelfId, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMySubShelfByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) DeleteMySubShelvesByIds(reqDto *dtos.DeleteMySubShelvesByIdsReqDto) (
	*dtos.DeleteMySubShelvesByIdsResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.SoftDeleteManyByIds(reqDto.Body.SubShelfIds, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMySubShelvesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
