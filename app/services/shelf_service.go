package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	logs "notezy-backend/app/logs"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	"notezy-backend/app/models/schemas/enums"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	lib "notezy-backend/shared/lib"
)

/* ============================== Interface & Instance ============================== */

type ShelfServiceInterface interface {
	GetMyShelfById(reqDto *dtos.GetMyShelfByIdReqDto) (*dtos.GetMyShelfByIdResDto, *exceptions.Exception)
	GetRecentShelves(reqDto *dtos.GetRecentShelvesReqDto) (*[]dtos.GetRecentShelvesResDto, *exceptions.Exception)
	CreateShelf(reqDto *dtos.CreateShelfReqDto) (*dtos.CreateShelfResDto, *exceptions.Exception)
	SynchronizeShelves(reqDto *dtos.SynchronizeShelvesReqDto) (*dtos.SynchronizeShelvesResDto, *exceptions.Exception)
	RestoreMyShelf(reqDto *dtos.RestoreMyShelfReqDto) (*dtos.RestoreMyShelfResDto, *exceptions.Exception)
	RestoreMyShelves(reqDto *dtos.RestoreMyShelvesReqDto) (*dtos.RestoreMyShelvesResDto, *exceptions.Exception)
	DeleteMyShelf(reqDto *dtos.DeleteMyShelfReqDto) (*dtos.DeleteMyShelfResDto, *exceptions.Exception)
	DeleteMyShelves(reqDto *dtos.DeleteMyShelvesReqDto) (*dtos.DeleteMyShelvesResDto, *exceptions.Exception)

	// services for private shelves
	SearchPrivateShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchShelfInput) (*gqlmodels.SearchShelfConnection, *exceptions.Exception)
}

type ShelfService struct {
	db *gorm.DB
}

func NewShelfService(db *gorm.DB) ShelfServiceInterface {
	return &ShelfService{
		db: db,
	}
}

/* ============================== Service Methods for Shelves ============================== */

func (s *ShelfService) GetMyShelfById(reqDto *dtos.GetMyShelfByIdReqDto) (*dtos.GetMyShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	shelf, exception := shelfRepository.GetOneById(reqDto.Body.ShelfId, reqDto.ContextFields.OwnerId, nil)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyShelfByIdResDto{
		Id:                       shelf.Id,
		Name:                     shelf.Name,
		EncodedStructure:         shelf.EncodedStructure,
		EncodedStructureByteSize: shelf.EncodedStructureByteSize,
		TotalShelfNodes:          shelf.TotalShelfNodes,
		TotalMaterials:           shelf.TotalMaterials,
		MaxWidth:                 shelf.MaxWidth,
		MaxDepth:                 shelf.MaxDepth,
		LastAnalyzedAt:           shelf.LastAnalyzedAt,
		DeletedAt:                shelf.DeletedAt,
		UpdatedAt:                shelf.UpdatedAt,
		CreatedAt:                shelf.CreatedAt,
	}, nil
}

func (s *ShelfService) GetRecentShelves(reqDto *dtos.GetRecentShelvesReqDto) (*[]dtos.GetRecentShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	resDto := []dtos.GetRecentShelvesResDto{}
	db := s.db.Model(&schemas.Shelf{})

	query := db.Where("owner_id = ?", reqDto.ContextFields.OwnerId)
	if len(strings.ReplaceAll(reqDto.Body.Query, " ", "")) > 0 {
		query = query.Where("name ILIKE ?", "%"+reqDto.Body.Query+"%")
	}
	logs.Info(reqDto.ContextFields.OwnerId)

	result := query.Order("updated_at DESC").
		Limit(int(reqDto.Body.Limit)).
		Offset(int(reqDto.Body.Offset)).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *ShelfService) CreateShelf(reqDto *dtos.CreateShelfReqDto) (*dtos.CreateShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	rootNode, exception := lib.NewShelfNode(reqDto.ContextFields.OwnerId, reqDto.Body.Name)
	if exception != nil {
		return nil, exception
	}
	encodedStructure, exception := lib.EncodeShelfNode(rootNode)
	if exception != nil {
		return nil, exception
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	now := time.Now()
	shelfId, exception := shelfRepository.CreateOneByOwnerId(
		reqDto.ContextFields.OwnerId,
		inputs.CreateShelfInput{
			Id:               rootNode.Id,
			Name:             reqDto.Body.Name,
			EncodedStructure: encodedStructure,
			LastAnalyzedAt:   &now,
		})
	if exception != nil {
		return nil, exception
	}
	if shelfId == nil || *shelfId != rootNode.Id {
		return nil, exceptions.Shelf.FailedToCreate("got nil shelf id")
	}

	return &dtos.CreateShelfResDto{
		Id:               *shelfId,
		EncodedStructure: encodedStructure,
		LastAnalyzedAt:   now,
		CreatedAt:        time.Now(),
	}, nil
}

func (s *ShelfService) SynchronizeShelves(reqDto *dtos.SynchronizeShelvesReqDto) (*dtos.SynchronizeShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	var updateInputs []inputs.PartialUpdateShelfInput
	for _, partialUpdate := range reqDto.Body.PartialUpdates {
		shelfRootNode, exception := lib.DecodeShelfNode(*partialUpdate.Values.EncodedStructure)
		if exception != nil {
			return nil, exception
		}
		summary, exception := shelfRootNode.AnalysisAndGenerateSummary()
		if exception != nil {
			exception.Log() // the system should not run into this section, may due to malicious attack
			return nil, exception
		}
		if summary == nil {
			return nil, exceptions.Shelf.CircularChildrenDetectedInShelfNode()
		}
		now := time.Now()
		updateInputs = append(updateInputs, inputs.PartialUpdateShelfInput{
			Values: inputs.UpdateShelfInput{
				Name:                     partialUpdate.Values.Name,
				EncodedStructure:         partialUpdate.Values.EncodedStructure,
				EncodedStructureByteSize: &summary.EncodedStructureByteSize,
				TotalShelfNodes:          &summary.TotalShelfNodes,
				TotalMaterials:           &summary.TotalMaterials,
				MaxWidth:                 &summary.MaxWidth,
				MaxDepth:                 &summary.MaxDepth,
				LastAnalyzedAt:           &now,
			},
			SetNull: partialUpdate.SetNull,
		})
	}

	exception := shelfRepository.DirectlyUpdateManyByIds(reqDto.Body.ShelfIds, reqDto.ContextFields.OwnerId, updateInputs)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SynchronizeShelvesResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *ShelfService) RestoreMyShelf(reqDto *dtos.RestoreMyShelfReqDto) (*dtos.RestoreMyShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	updatedAt, exception := shelfRepository.RestoreSoftDeletedOneById(reqDto.Body.ShelfId, reqDto.ContextFields.OwnerId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyShelfResDto{
		UpdatedAt: updatedAt,
	}, nil
}

func (s *ShelfService) RestoreMyShelves(reqDto *dtos.RestoreMyShelvesReqDto) (*dtos.RestoreMyShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	updatedAt, exception := shelfRepository.RestoreSoftDeletedManyByIds(reqDto.Body.ShelfIds, reqDto.ContextFields.OwnerId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyShelvesResDto{
		UpdatedAt: updatedAt,
	}, nil
}

func (s *ShelfService) DeleteMyShelf(reqDto *dtos.DeleteMyShelfReqDto) (*dtos.DeleteMyShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	deletedAt, exception := shelfRepository.SoftDeleteOneById(reqDto.Body.ShelfId, reqDto.ContextFields.OwnerId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyShelfResDto{
		DeletedAt: deletedAt,
	}, nil
}

func (s *ShelfService) DeleteMyShelves(reqDto *dtos.DeleteMyShelvesReqDto) (*dtos.DeleteMyShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	deletedAt, exception := shelfRepository.SoftDeleteManyByIds(reqDto.Body.ShelfIds, reqDto.ContextFields.OwnerId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyShelvesResDto{
		DeletedAt: deletedAt,
	}, nil
}

/* ============================== Service Methods for  ============================== */

func (s *ShelfService) SearchPrivateShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchShelfInput) (*gqlmodels.SearchShelfConnection, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	startTime := time.Now()

	query := s.db.WithContext(ctx).Model(&schemas.Shelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"ShelfTable\".id = uts.shelf_id").
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, exception := lib.DecodeSearchCursor[gqlmodels.SearchShelfCursorFields](*gqlInput.After)
		if exception != nil {
			return nil, exception
		}

		query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchShelfSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchShelfSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchShelfSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = max(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var shelves []schemas.Shelf
	if err := query.Find(&shelves).Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	hasNextPage := len(shelves) > limit
	searchEdges := make([]*gqlmodels.SearchShelfEdge, len(shelves))

	for index, shelf := range shelves {
		searchCursor := lib.SearchCursor[gqlmodels.SearchShelfCursorFields]{
			Fields: gqlmodels.SearchShelfCursorFields{
				ID: shelf.Id,
			},
		}
		encodedSearchCursor, exception := searchCursor.EncodeSearchCursor()
		if exception != nil {
			return nil, exception
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchShelfEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                shelf.ToPrivateShelf(),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchShelfConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
