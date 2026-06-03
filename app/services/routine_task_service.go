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

type RoutineTaskServiceInterface interface {
	GetMyRoutineTaskById(ctx context.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto) (*dtos.GetMyRoutineTaskByIdResDto, *exceptions.Exception)
	CreateRoutineTaskByStationId(ctx context.Context, reqDto *dtos.CreateRoutineTaskByStationIdReqDto) (*dtos.CreateRoutineTaskByStationIdResDto, *exceptions.Exception)
	UpdateMyRoutineTaskById(ctx context.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto) (*dtos.UpdateMyRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTaskById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto) (*dtos.HardDeleteMyRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTasksByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto) (*dtos.HardDeleteMyRoutineTasksByIdsResDto, *exceptions.Exception)

	// services for private routine tasks
	SearchPrivateRoutineTasks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTaskInput) (*gqlmodels.SearchRoutineTaskConnection, *exceptions.Exception)
}

type RoutineTaskService struct {
	db                    *gorm.DB
	routineTaskRepository repositories.RoutineTaskRepositoryInterface
}

func NewRoutineTaskService(
	db *gorm.DB,
	routineTaskRepository repositories.RoutineTaskRepositoryInterface,
) RoutineTaskServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RoutineTaskService{
		db:                    db,
		routineTaskRepository: routineTaskRepository,
	}
}

/* ============================== Service Methods for RoutineTask ============================== */

func (s *RoutineTaskService) GetMyRoutineTaskById(
	ctx context.Context,
	reqDto *dtos.GetMyRoutineTaskByIdReqDto,
) (*dtos.GetMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	routineTask, exception := s.routineTaskRepository.GetOneById(
		reqDto.Param.RoutineTaskId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyRoutineTaskByIdResDto{
		Id:              routineTask.Id,
		StationId:       routineTask.StationId,
		Title:           routineTask.Title,
		Purpose:         routineTask.Purpose,
		Payload:         routineTask.Payload,
		Priority:        routineTask.Priority,
		Status:          routineTask.Status,
		Attempts:        routineTask.Attempts,
		MaxAttempts:     routineTask.MaxAttempts,
		ScheduledAt:     routineTask.ScheduledAt,
		ActualStartedAt: routineTask.ActualStartedAt,
		ActualEndedAt:   routineTask.ActualEndedAt,
		UpdatedAt:       routineTask.UpdatedAt,
		CreatedAt:       routineTask.CreatedAt,
	}, nil
}

func (s *RoutineTaskService) CreateRoutineTaskByStationId(
	ctx context.Context,
	reqDto *dtos.CreateRoutineTaskByStationIdReqDto,
) (*dtos.CreateRoutineTaskByStationIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	newRoutineTaskId, exception := s.routineTaskRepository.CreateOneByStationId(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		inputs.CreateRoutineTaskInput{
			Title:       reqDto.Body.Title,
			Purpose:     reqDto.Body.Purpose,
			Payload:     reqDto.Body.Payload,
			Priority:    reqDto.Body.Priority,
			MaxAttempts: reqDto.Body.MaxAttempts,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutineTaskByStationIdResDto{
		Id:        *newRoutineTaskId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) UpdateMyRoutineTaskById(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutineTaskByIdReqDto,
) (*dtos.UpdateMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	updatedRoutineTask, exception := s.routineTaskRepository.UpdateOneById(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRoutineTaskInput{
			Values: inputs.UpdateRoutineTaskInput{
				StationId:   reqDto.Body.Values.StationId,
				Title:       reqDto.Body.Values.Title,
				Purpose:     reqDto.Body.Values.Purpose,
				Payload:     reqDto.Body.Values.Payload,
				Priority:    reqDto.Body.Values.Priority,
				MaxAttempts: reqDto.Body.Values.MaxAttempts,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRoutineTaskByIdResDto{
		UpdatedAt: updatedRoutineTask.UpdatedAt,
	}, nil
}

func (s *RoutineTaskService) HardDeleteMyRoutineTaskById(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto,
) (*dtos.HardDeleteMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.routineTaskRepository.HardDeleteOneById(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutineTaskByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) HardDeleteMyRoutineTasksByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto,
) (*dtos.HardDeleteMyRoutineTasksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.routineTaskRepository.HardDeleteManyByIds(
		reqDto.Body.RoutineTaskIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutineTasksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL RoutineTask ============================== */

func (s *RoutineTaskService) SearchPrivateRoutineTasks(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTaskInput,
) (*gqlmodels.SearchRoutineTaskConnection, *exceptions.Exception) {
	type PrivateRoutineTask struct {
		schemas.RoutineTask
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.RoutineTask{}).
		Select("\"RoutineTaskTable\".*, uts.permission AS permission").
		Joins("LEFT JOIN \"UsersToStationsTable\" uts ON \"RoutineTaskTable\".station_id = uts.station_id").
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)

	if gqlInput.StationID != nil {
		query = query.Where(
			"\"RoutineTaskTable\".station_id = ?",
			*gqlInput.StationID,
		)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"title ILIKE ? OR purpose::text ILIKE ? OR payload::text ILIKE ?",
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRoutineTaskCursorFields](*gqlInput.After)
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
		case gqlmodels.SearchRoutineTaskSortByTitle:
			query.Order("title " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByPurpose:
			query.Order("purpose " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByPriority:
			query.Order("priority " + cending).
				Order("scheduled_at " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByStatus:
			query.Order("status " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByAttempts:
			query.Order("attempts " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByMaxAttempts:
			query.Order("max_attempts " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByScheduledAt:
			query.Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByActualStartedAt:
			query.Order("actual_started_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByActualEndedAt:
			query.Order("actual_ended_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("scheduled_at " + cending).
				Order("priority " + cending).
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

	var routineTasks []PrivateRoutineTask
	if err := query.Find(&routineTasks).Error; err != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	hasNextPage := len(routineTasks) > limit
	searchEdges := make([]*gqlmodels.SearchRoutineTaskEdge, len(routineTasks))

	for index, routineTask := range routineTasks {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRoutineTaskCursorFields]{
			Fields: gqlmodels.SearchRoutineTaskCursorFields{
				ID: routineTask.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRoutineTaskEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                routineTask.RoutineTask.ToPrivateRoutineTask(),
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

	return &gqlmodels.SearchRoutineTaskConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
