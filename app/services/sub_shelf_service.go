package services

import (
	"time"

	"github.com/google/uuid"
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
)

/* ============================== Interface & Instance ============================== */

type SubShelfServiceInterface interface {
	GetMySubShelfById(reqDto *dtos.GetMySubShelfByIdReqDto) (*dtos.GetMySubShelfByIdResDto, *exceptions.Exception)
	GetAllSubShelvesByRootShelfId(reqDto *dtos.GetAllSubShelvesByRootShelfIdReqDto) (*dtos.GetAllSubShelvesByRootShelfIdResDto, *exceptions.Exception)
	CreateSubShelfByRootShelfId(reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) (*dtos.CreateSubShelfByRootShelfIdResDto, *exceptions.Exception)
	RenameMySubShelfById(reqDto *dtos.RenameMySubShelfByIdReqDto) (*dtos.RenameMySubShelfByIdResDto, *exceptions.Exception)
	MoveMySubShelf(reqDto *dtos.MoveMySubShelfReqDto) (*dtos.MoveMySubShelfResDto, *exceptions.Exception)
	MoveMySubShelves(reqDto *dtos.MoveMySubShelvesReqDto) (*dtos.MoveMySubShelvesResDto, *exceptions.Exception)
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
		return nil, exceptions.User.InvalidInput().WithError(err)
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
		UpdatedAt:      subShelf.UpdatedAt,
		CreatedAt:      subShelf.CreatedAt,
	}, nil
}

func (s *SubShelfService) GetAllSubShelvesByRootShelfId(reqDto *dtos.GetAllSubShelvesByRootShelfIdReqDto) (
	*dtos.GetAllSubShelvesByRootShelfIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	resDto := dtos.GetAllSubShelvesByRootShelfIdResDto{}
	result := s.db.Model(&schemas.SubShelf{}).
		Joins("LEFT JOIN \"RootShelfTable\" rs ON \"SubShelfTable\".root_shelf_id = rs.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON rs.id = uts.root_shelf_id").
		Where("\"SubShelfTable\".root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *SubShelfService) CreateSubShelfByRootShelfId(reqDto *dtos.CreateSubShelfByRootShelfIdReqDto) (
	*dtos.CreateSubShelfByRootShelfIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	_, exception := subShelfRepository.CreateOneByUserId(
		reqDto.ContextFields.UserId,
		inputs.CreateSubShelfInput{
			Name:        reqDto.Body.Name,
			RootShelfId: reqDto.Body.RootShelfId,
			Path:        reqDto.Body.Path,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateSubShelfByRootShelfIdResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) RenameMySubShelfById(reqDto *dtos.RenameMySubShelfByIdReqDto) (
	*dtos.RenameMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	subShelf, exception := subShelfRepository.UpdateOneById(
		reqDto.Body.RootShelfId,
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

	return &dtos.RenameMySubShelfByIdResDto{
		UpdatedAt: subShelf.UpdatedAt,
	}, nil
}

func (s *SubShelfService) MoveMySubShelf(reqDto *dtos.MoveMySubShelfReqDto) (
	*dtos.MoveMySubShelfResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	from, exception := subShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Body.SourceSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}
	to, exception := subShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Body.DestinationSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}

	if len(from.Path)+len(to.Path) > constants.MaxShelfTreeDepth {
		return nil, exceptions.Shelf.MaximumDepthExceeded(
			int32(len(from.Path)+len(to.Path)),
			constants.MaxShelfTreeDepth,
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
	update := schemas.SubShelf{
		RootShelfId:    to.RootShelfId,
		PrevSubShelfId: &reqDto.Body.DestinationSubShelfId,
		Path:           to.Path,
	}
	result := s.db.Model(&schemas.SubShelf{}).
		Where("id = ?", reqDto.Body.SourceSubShelfId).
		Updates(&update)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
	}

	return &dtos.MoveMySubShelfResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) MoveMySubShelves(reqDto *dtos.MoveMySubShelvesReqDto) (
	*dtos.MoveMySubShelvesResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	froms, exception := subShelfRepository.CheckPermissionAndGetManyByIds(
		reqDto.Body.SourceSubShelfIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}
	to, exception := subShelfRepository.CheckPermissionAndGetOneById(
		reqDto.Body.DestinationSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}

	sourceSubShelfIdMap := make(map[uuid.UUID]bool, 0)
	for _, from := range *froms {
		if len(from.Path)+len(to.Path) > constants.MaxShelfTreeDepth {
			return nil, exceptions.Shelf.MaximumDepthExceeded(
				int32(len(from.Path)+len(to.Path)),
				constants.MaxShelfTreeDepth,
			)
		}
		sourceSubShelfIdMap[from.Id] = true
	}

	for _, parent := range to.Path {
		if sourceSubShelfIdMap[parent] {
			return nil, exceptions.Shelf.InsertParentIntoItsChildren(
				reqDto.Body.DestinationSubShelfId,
				parent,
			)
		}
	}

	to.Path = append(to.Path, to.Id)
	update := schemas.SubShelf{
		RootShelfId:    to.RootShelfId,
		PrevSubShelfId: &reqDto.Body.DestinationSubShelfId,
		Path:           to.Path,
	}
	result := s.db.Model(&schemas.SubShelf{}).
		Where("id IN ?", reqDto.Body.SourceSubShelfIds).
		Updates(&update)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
	}

	return &dtos.MoveMySubShelvesResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *SubShelfService) DeleteMySubShelfById(reqDto *dtos.DeleteMySubShelfByIdReqDto) (
	*dtos.DeleteMySubShelfByIdResDto, *exceptions.Exception,
) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.DeleteOneById(reqDto.Body.SubShelfId, reqDto.ContextFields.UserId)
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
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)

	exception := subShelfRepository.DeleteManyByIds(reqDto.Body.SubShelfIds, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMySubShelvesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
