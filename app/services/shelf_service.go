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
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	lib "notezy-backend/shared/lib"
)

/* ============================== Interface & Instance ============================== */

type ShelfServiceInterface interface {
	CreateShelf(reqDto *dtos.CreateShelfReqDto) (*dtos.CreateShelfResDto, *exceptions.Exception)
	SynchronizeShelves(reqDto *dtos.SynchronizeShelvesReqDto) (*dtos.SynchronizeShelvesResDto, *exceptions.Exception)

	// services for private shelves
	SearchPrivateShelves(ctx context.Context, ownerId uuid.UUID, gqlInput gqlmodels.SearchShelfInput) (*gqlmodels.SearchShelfConnection, *exceptions.Exception)
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

func (s *ShelfService) CreateShelf(reqDto *dtos.CreateShelfReqDto) (*dtos.CreateShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	_, exception := shelfRepository.CreateOneByOwnerId(reqDto.OwnerId, inputs.CreateShelfInput{Name: reqDto.Name})
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateShelfResDto{CreatedAt: time.Now()}, nil
}

func (s *ShelfService) SynchronizeShelves(reqDto *dtos.SynchronizeShelvesReqDto) (*dtos.SynchronizeShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	shelfRepository := repositories.NewShelfRepository(s.db)

	var updateInputs []inputs.PartialUpdateShelfInput
	for _, partialUpdate := range reqDto.PartialUpdates {
		shelfRootNode, exception := lib.DecodeShelfNode(*partialUpdate.Values.EncodedStructure)
		if exception != nil {
			return nil, exception
		}
		isSimple, exception := shelfRootNode.IsChildrenSimple()
		if exception != nil {
			exception.Log() // the system should not run into this section, may due to malicious attack
			return nil, exception
		}
		if !isSimple {
			return nil, exceptions.Shelf.CircularChildrenDetectedInShelfNode()
		}
		updateInputs = append(updateInputs, inputs.PartialUpdateShelfInput{
			Values: inputs.UpdateShelfInput{
				Name:             partialUpdate.Values.Name,
				EncodedStructure: partialUpdate.Values.EncodedStructure,
			},
			SetNull: partialUpdate.SetNull,
		})
	}

	exception := shelfRepository.DirectlyUpdateManyByIds(reqDto.ShelfIds, reqDto.OwnerId, updateInputs)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SynchronizeShelvesResDto{UpdatedAt: time.Now()}, nil
}

/* ============================== Service Methods for  ============================== */

func (s *ShelfService) SearchPrivateShelves(ctx context.Context, ownerId uuid.UUID, gqlInput gqlmodels.SearchShelfInput) (*gqlmodels.SearchShelfConnection, *exceptions.Exception) {
	startTime := time.Now()

	query := s.db.WithContext(ctx).Model(&schemas.Shelf{}).
		Where("owner_id = ?", ownerId)

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
