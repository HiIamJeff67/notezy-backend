package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
)

type RoutineTagServiceInterface interface {
	GetMyRoutineTagById(ctx context.Context, reqDto *dtos.GetMyRoutineTagByIdReqDto) (*dtos.GetMyRoutineTagByIdResDto, *exceptions.Exception)
	GetAllMyRoutineTags(ctx context.Context, reqDto *dtos.GetAllMyRoutineTagsReqDto) (*dtos.GetAllMyRoutineTagsResDto, *exceptions.Exception)
	CreateRoutineTag(ctx context.Context, reqDto *dtos.CreateRoutineTagReqDto) (*dtos.CreateRoutineTagResDto, *exceptions.Exception)
	CreateRoutineTags(ctx context.Context, reqDto *dtos.CreateRoutineTagsReqDto) (*dtos.CreateRoutineTagsResDto, *exceptions.Exception)
	UpdateMyRoutineTagById(ctx context.Context, reqDto *dtos.UpdateMyRoutineTagByIdReqDto) (*dtos.UpdateMyRoutineTagByIdResDto, *exceptions.Exception)
	UpdateMyRoutineTagsByIds(ctx context.Context, reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto) (*dtos.UpdateMyRoutineTagsByIdsResDto, *exceptions.Exception)
	HardDeleteMyRoutineTagById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto) (*dtos.HardDeleteMyRoutineTagByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTagsByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto) (*dtos.HardDeleteMyRoutineTagsByIdsResDto, *exceptions.Exception)

	SearchPrivateRoutineTags(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTagInput) (*gqlmodels.SearchRoutineTagConnection, *exceptions.Exception)
}

type RoutineTagService struct {
	db                   *gorm.DB
	routineTagRepository repositories.RoutineTagRepositoryInterface
}

func NewRoutineTagService(
	db *gorm.DB,
	routineTagRepository repositories.RoutineTagRepositoryInterface,
) RoutineTagServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RoutineTagService{
		db:                   db,
		routineTagRepository: routineTagRepository,
	}
}

/* ============================== Service Methods for RoutineTag ============================== */

func (s *RoutineTagService) GetMyRoutineTagById(
	ctx context.Context,
	reqDto *dtos.GetMyRoutineTagByIdReqDto,
) (*dtos.GetMyRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.IsDeleted != nil && *reqDto.Param.IsDeleted {
		return nil, exceptions.RoutineTag.NotFound()
	}

	db := s.db.WithContext(ctx)

	routineTag, exception := s.routineTagRepository.GetOneById(
		reqDto.Param.RoutineTagId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyRoutineTagByIdResDto{
		Id:        routineTag.Id,
		Name:      routineTag.Name,
		Color:     routineTag.Color,
		Icon:      routineTag.Icon,
		UpdatedAt: routineTag.UpdatedAt,
		CreatedAt: routineTag.CreatedAt,
	}, nil
}

func (s *RoutineTagService) GetAllMyRoutineTags(
	ctx context.Context,
	reqDto *dtos.GetAllMyRoutineTagsReqDto,
) (*dtos.GetAllMyRoutineTagsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.AreDeleted != nil && *reqDto.Param.AreDeleted {
		resDto := dtos.GetAllMyRoutineTagsResDto{}
		return &resDto, nil
	}

	db := s.db.WithContext(ctx)

	routineTags, exception := s.routineTagRepository.GetAllByUserId(
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutineTagsResDto, len(routineTags))
	for index, routineTag := range routineTags {
		resDto[index] = dtos.GetMyRoutineTagByIdResDto{
			Id:        routineTag.Id,
			Name:      routineTag.Name,
			Color:     routineTag.Color,
			Icon:      routineTag.Icon,
			UpdatedAt: routineTag.UpdatedAt,
			CreatedAt: routineTag.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *RoutineTagService) CreateRoutineTag(
	ctx context.Context,
	reqDto *dtos.CreateRoutineTagReqDto,
) (*dtos.CreateRoutineTagResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	newRoutineTagId, exception := s.routineTagRepository.CreateOneByUserId(
		reqDto.ContextFields.UserId,
		inputs.CreateRoutineTagInput{
			Id:    reqDto.Body.Id,
			Name:  reqDto.Body.Name,
			Color: reqDto.Body.Color,
			Icon:  reqDto.Body.Icon,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutineTagResDto{
		Id:        *newRoutineTagId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) CreateRoutineTags(
	ctx context.Context,
	reqDto *dtos.CreateRoutineTagsReqDto,
) (*dtos.CreateRoutineTagsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkCreateRoutineTagInput, len(reqDto.Body.CreatedRoutineTags))
	for index, createdRoutineTag := range reqDto.Body.CreatedRoutineTags {
		input[index] = inputs.BulkCreateRoutineTagInput{
			Id:    createdRoutineTag.Id,
			Name:  createdRoutineTag.Name,
			Color: createdRoutineTag.Color,
			Icon:  createdRoutineTag.Icon,
		}
	}
	newRoutineTagIds, exception := s.routineTagRepository.BulkCreateManyByUserId(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutineTagsResDto{
		Ids:       newRoutineTagIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) UpdateMyRoutineTagById(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutineTagByIdReqDto,
) (*dtos.UpdateMyRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	updatedRoutineTag, exception := s.routineTagRepository.UpdateOneById(
		reqDto.Body.RoutineTagId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRoutineTagInput{
			Values: inputs.UpdateRoutineTagInput{
				Name:  reqDto.Body.Values.Name,
				Color: reqDto.Body.Values.Color,
				Icon:  reqDto.Body.Values.Icon,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRoutineTagByIdResDto{
		UpdatedAt: updatedRoutineTag.UpdatedAt,
	}, nil
}

func (s *RoutineTagService) UpdateMyRoutineTagsByIds(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto,
) (*dtos.UpdateMyRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkUpdateRoutineTagInput, len(reqDto.Body.UpdatedRoutineTags))
	for index, updatedRoutineTag := range reqDto.Body.UpdatedRoutineTags {
		input[index] = inputs.BulkUpdateRoutineTagInput{
			Id: updatedRoutineTag.RoutineTagId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateRoutineTagInput]{
				Values: inputs.UpdateRoutineTagInput{
					Name:  updatedRoutineTag.Values.Name,
					Color: updatedRoutineTag.Values.Color,
					Icon:  updatedRoutineTag.Values.Icon,
				},
				SetNull: updatedRoutineTag.SetNull,
			},
		}
	}
	exception := s.routineTagRepository.BulkUpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRoutineTagsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteMyRoutineTagById(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto,
) (*dtos.HardDeleteMyRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineTagRepository.HardDeleteOneById(
		reqDto.Body.RoutineTagId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutineTagByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteMyRoutineTagsByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto,
) (*dtos.HardDeleteMyRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTag.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineTagRepository.HardDeleteManyByIds(
		reqDto.Body.RoutineTagIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutineTagsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL RoutineTag ============================== */

func (s *RoutineTagService) SearchPrivateRoutineTags(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTagInput,
) (*gqlmodels.SearchRoutineTagConnection, *exceptions.Exception) {
	type PrivateRoutineTag struct {
		schemas.RoutineTag
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	query := db.Model(&schemas.RoutineTag{}).
		Select(`"RoutineTagTable".*, utrt.permission AS permission`).
		Joins(`LEFT JOIN "UsersToRoutineTagsTable" utrt ON "RoutineTagTable".id = utrt.tag_id`).
		Where("utrt.user_id = ? AND utrt.permission IN ?", userId, allowedPermissions)

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRoutineTagCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRoutineTagSortByName:
			query = query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTagSortByLastUpdate:
			query = query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTagSortByCreatedAt:
			query = query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query = query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var routineTags []PrivateRoutineTag
	if err := query.Find(&routineTags).Error; err != nil {
		return nil, exceptions.RoutineTag.NotFound().WithOrigin(err)
	}

	hasNextPage := len(routineTags) > limit
	searchEdges := make([]*gqlmodels.SearchRoutineTagEdge, len(routineTags))

	for index, routineTag := range routineTags {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRoutineTagCursorFields]{
			Fields: gqlmodels.SearchRoutineTagCursorFields{
				ID: routineTag.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRoutineTagEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                routineTag.RoutineTag.ToPrivateRoutineTag(),
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

	return &gqlmodels.SearchRoutineTagConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
