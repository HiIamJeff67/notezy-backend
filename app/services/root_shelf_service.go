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
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	searchcursor "notezy-backend/shared/lib/searchcursor"
	"notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type RootShelfServiceInterface interface {
	GetMyRootShelfById(ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception)
	SearchRecentRootShelves(ctx context.Context, reqDto *dtos.SearchRecentRootShelvesReqDto) (*dtos.SearchRecentRootShelvesResDto, *exceptions.Exception)
	CreateRootShelf(ctx context.Context, reqDto *dtos.CreateRootShelfReqDto) (*dtos.CreateRootShelfResDto, *exceptions.Exception)
	UpdateMyRootShelfById(ctx context.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto) (*dtos.UpdateMyRootShelfByIdResDto, *exceptions.Exception)
	RestoreMyRootShelfById(ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception)
	RestoreMyRootShelvesByIds(ctx context.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto) (*dtos.RestoreMyRootShelvesByIdsResDto, *exceptions.Exception)
	DeleteMyRootShelfById(ctx context.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto) (*dtos.DeleteMyRootShelfByIdResDto, *exceptions.Exception)
	DeleteMyRootShelvesByIds(ctx context.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto) (*dtos.DeleteMyRootShelvesByIdsResDto, *exceptions.Exception)

	// services for private shelves
	SearchPrivateShelves(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception)
}

type RootShelfService struct {
	db                  *gorm.DB
	rootShelfRepository repositories.RootShelfRepositoryInterface
}

func NewRootShelfService(
	db *gorm.DB,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
) RootShelfServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RootShelfService{
		db:                  db,
		rootShelfRepository: rootShelfRepository,
	}
}

/* ============================== Service Methods for Shelf ============================== */

func (s *RootShelfService) GetMyRootShelfById(
	ctx context.Context, reqDto *dtos.GetMyRootShelfByIdReqDto,
) (*dtos.GetMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	shelf, exception := s.rootShelfRepository.GetOneById(
		reqDto.Param.RootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyRootShelfByIdResDto{
		Id:             shelf.Id,
		Name:           shelf.Name,
		SubShelfCount:  shelf.SubShelfCount,
		ItemCount:      shelf.ItemCount,
		LastAnalyzedAt: shelf.LastAnalyzedAt,
		DeletedAt:      shelf.DeletedAt,
		UpdatedAt:      shelf.UpdatedAt,
		CreatedAt:      shelf.CreatedAt,
	}, nil
}

func (s *RootShelfService) SearchRecentRootShelves(
	ctx context.Context, reqDto *dtos.SearchRecentRootShelvesReqDto,
) (*dtos.SearchRecentRootShelvesResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	resDto := dtos.SearchRecentRootShelvesResDto{}

	query := db.Model(&schemas.RootShelf{}).
		Where("owner_id = ? AND \"RootShelfTable\".deleted_at IS NULL",
			reqDto.ContextFields.UserId,
		)
	if len(strings.ReplaceAll(reqDto.Param.Query, " ", "")) > 0 {
		query = query.Where("name ILIKE ?", "%"+reqDto.Param.Query+"%")
	}

	result := query.Order("updated_at DESC").
		Limit(int(reqDto.Param.Limit)).
		Offset(int(reqDto.Param.Offset)).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *RootShelfService) CreateRootShelf(
	ctx context.Context, reqDto *dtos.CreateRootShelfReqDto,
) (*dtos.CreateRootShelfResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	now := time.Now()
	shelfId, exception := s.rootShelfRepository.CreateOneByOwnerId(
		reqDto.ContextFields.OwnerId,
		inputs.CreateRootShelfInput{
			Name:           reqDto.Body.Name,
			LastAnalyzedAt: &now,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRootShelfResDto{
		Id:             *shelfId,
		LastAnalyzedAt: now,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *RootShelfService) UpdateMyRootShelfById(
	ctx context.Context, reqDto *dtos.UpdateMyRootShelfByIdReqDto,
) (*dtos.UpdateMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	rootShelf, exception := s.rootShelfRepository.UpdateOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRootShelfInput{
			Values: inputs.UpdateRootShelfInput{
				Name: reqDto.Body.Values.Name,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRootShelfByIdResDto{
		UpdatedAt: rootShelf.UpdatedAt,
	}, nil
}

func (s *RootShelfService) RestoreMyRootShelfById(
	ctx context.Context, reqDto *dtos.RestoreMyRootShelfByIdReqDto,
) (*dtos.RestoreMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.RestoreSoftDeletedOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.OwnerId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyRootShelfByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) RestoreMyRootShelvesByIds(
	ctx context.Context, reqDto *dtos.RestoreMyRootShelvesByIdsReqDto,
) (*dtos.RestoreMyRootShelvesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.OwnerId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyRootShelvesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelfById(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelfByIdReqDto,
) (*dtos.DeleteMyRootShelfByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.SoftDeleteOneById(
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.OwnerId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRootShelfByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RootShelfService) DeleteMyRootShelvesByIds(
	ctx context.Context, reqDto *dtos.DeleteMyRootShelvesByIdsReqDto,
) (*dtos.DeleteMyRootShelvesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.rootShelfRepository.SoftDeleteManyByIds(
		reqDto.Body.RootShelfIds,
		reqDto.ContextFields.OwnerId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRootShelvesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for  ============================== */

func (s *RootShelfService) SearchPrivateShelves(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRootShelfInput,
) (*gqlmodels.SearchRootShelfConnection, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.WithContext(ctx).Model(&schemas.RootShelf{}).
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON \"RootShelfTable\".id = uts.root_shelf_id").
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Where("\"RootShelfTable\".deleted_at IS NULL")

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRootShelfCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithError(err)
		}

		query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRootShelfSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRootShelfSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRootShelfSortByCreatedAt:
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

	var shelves []schemas.RootShelf
	if err := query.Find(&shelves).Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	hasNextPage := len(shelves) > limit
	searchEdges := make([]*gqlmodels.SearchRootShelfEdge, len(shelves))

	for index, shelf := range shelves {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRootShelfCursorFields]{
			Fields: gqlmodels.SearchRootShelfCursorFields{
				ID: shelf.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithError(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRootShelfEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                shelf.ToPrivateRootShelf(),
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

	return &gqlmodels.SearchRootShelfConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
