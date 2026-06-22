package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineServiceInterface interface {
	GetMyRoutineById(ctx context.Context, reqDto *dtos.GetMyRoutineByIdReqDto) (*dtos.GetMyRoutineByIdResDto, *exceptions.Exception)
	GetMyRoutinesByStationId(ctx context.Context, reqDto *dtos.GetMyRoutinesByStationIdReqDto) (*dtos.GetMyRoutinesByStationIdResDto, *exceptions.Exception)
	GetAllMyRoutinesByTimeRange(ctx context.Context, reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto) (*dtos.GetAllMyRoutinesByTimeRangeResDto, *exceptions.Exception)
	CreateRoutineByStationId(ctx context.Context, reqDto *dtos.CreateRoutineByStationIdReqDto) (*dtos.CreateRoutineByStationIdResDto, *exceptions.Exception)
	CreateRoutinesByStationIds(ctx context.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto) (*dtos.CreateRoutinesByStationIdsResDto, *exceptions.Exception)
	UpdateMyRoutineById(ctx context.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto) (*dtos.UpdateMyRoutineByIdResDto, *exceptions.Exception)
	UpdateMyRoutinesByIds(ctx context.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto) (*dtos.UpdateMyRoutinesByIdsResDto, *exceptions.Exception)
	LinkRoutineTagById(ctx context.Context, reqDto *dtos.LinkRoutineTagByIdReqDto) (*dtos.LinkRoutineTagByIdResDto, *exceptions.Exception)
	BulkLinkRoutineTagsByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineTagsByIdsReqDto) (*dtos.BulkLinkRoutineTagsByIdsResDto, *exceptions.Exception)
	LinkRoutineTaskById(ctx context.Context, reqDto *dtos.LinkRoutineTaskByIdReqDto) (*dtos.LinkRoutineTaskByIdResDto, *exceptions.Exception)
	BulkLinkRoutineTasksByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineTasksByIdsReqDto) (*dtos.BulkLinkRoutineTasksByIdsResDto, *exceptions.Exception)
	LinkRoutineItemById(ctx context.Context, reqDto *dtos.LinkRoutineItemByIdReqDto) (*dtos.LinkRoutineItemByIdResDto, *exceptions.Exception)
	BulkLinkRoutineItemsByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineItemsByIdsReqDto) (*dtos.BulkLinkRoutineItemsByIdsResDto, *exceptions.Exception)
	RestoreMyRoutineById(ctx context.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto) (*dtos.RestoreMyRoutineByIdResDto, *exceptions.Exception)
	RestoreMyRoutinesByIds(ctx context.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto) (*dtos.RestoreMyRoutinesByIdsResDto, *exceptions.Exception)
	DeleteMyRoutineById(ctx context.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto) (*dtos.DeleteMyRoutineByIdResDto, *exceptions.Exception)
	DeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto) (*dtos.DeleteMyRoutinesByIdsResDto, *exceptions.Exception)
	HardDeleteMyRoutineById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto) (*dtos.HardDeleteMyRoutineByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto) (*dtos.HardDeleteMyRoutinesByIdsResDto, *exceptions.Exception)
	VisualizeMyRoutineStatusCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineStatusCountReqDto) (*dtos.VisualizeMyRoutineStatusCountResDto, *exceptions.Exception)
	VisualizeMyRoutinePeriodCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutinePeriodCountReqDto) (*dtos.VisualizeMyRoutinePeriodCountResDto, *exceptions.Exception)
	VisualizeMyRoutineScheduledStartAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineScheduledStartAtCountReqDto) (*dtos.VisualizeMyRoutineScheduledStartAtCountResDto, *exceptions.Exception)
	VisualizeMyRoutineScheduledEndAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineScheduledEndAtCountReqDto) (*dtos.VisualizeMyRoutineScheduledEndAtCountResDto, *exceptions.Exception)

	// services for graphql routines
	SearchPrivateRoutines(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineInput) (*gqlmodels.SearchRoutineConnection, *exceptions.Exception)
}

type RoutineService struct {
	db                    *gorm.DB
	routineScope          scopes.RoutineScopeInterface
	stationRepository     repositories.StationRepositoryInterface
	routineRepository     repositories.RoutineRepositoryInterface
	routineTagRepository  repositories.RoutineTagRepositoryInterface
	routineTaskRepository repositories.RoutineTaskRepositoryInterface
	itemRepository        repositories.ItemRepositoryInterface
}

func NewRoutineService(
	db *gorm.DB,
	routineScope scopes.RoutineScopeInterface,
	stationRepository repositories.StationRepositoryInterface,
	routineRepository repositories.RoutineRepositoryInterface,
	routineTagRepository repositories.RoutineTagRepositoryInterface,
	routineTaskRepository repositories.RoutineTaskRepositoryInterface,
	itemRepository repositories.ItemRepositoryInterface,
) RoutineServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RoutineService{
		db:                    db,
		routineScope:          routineScope,
		stationRepository:     stationRepository,
		routineRepository:     routineRepository,
		routineTagRepository:  routineTagRepository,
		routineTaskRepository: routineTaskRepository,
		itemRepository:        itemRepository,
	}
}

/* ============================== Helper Functions ============================== */

func (s *RoutineService) visualizeMyRoutineTimeCount(
	ctx context.Context,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	timeHourUnit int,
	queryRangeStartedAt time.Time,
	queryRangeEndedAt time.Time,
	columnName string,
	fieldName string,
) ([]dtos.TwoDimensionalDatum[int64], *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	var buckets []struct {
		BucketStart  time.Time `gorm:"column:bucket_start;"`
		RoutineCount int64     `gorm:"column:routine_count;"`
	}
	result := db.
		Table(
			`generate_series(
				date_trunc('hour', ?::timestamptz),
				date_trunc('hour', ?::timestamptz - interval '1 microsecond'),
				?::integer * interval '1 hour'
			) AS buckets(bucket_start)`,
			queryRangeStartedAt,
			queryRangeEndedAt,
			timeHourUnit,
		).
		Select(`
			buckets.bucket_start AS bucket_start,
			COUNT(uts.station_id) AS routine_count
		`).
		Joins(
			`LEFT JOIN "RoutineTable" routine
				ON routine.`+columnName+` >= buckets.bucket_start
				AND routine.`+columnName+` < buckets.bucket_start + ?::integer * interval '1 hour'
				AND routine.deleted_at IS NULL`,
			timeHourUnit,
		).
		Joins(
			`LEFT JOIN "UsersToStationsTable" uts
				ON uts.station_id = routine.station_id
				AND uts.user_id = ?
				AND uts.permission = ?`,
			userId,
			permission,
		).
		Group("buckets.bucket_start").
		Order("buckets.bucket_start ASC").
		Scan(&buckets)
	if err := result.Error; err != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(err)
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(buckets))
	for index, bucket := range buckets {
		bucketEnd := bucket.BucketStart.Add(time.Duration(timeHourUnit) * time.Hour)
		x := bucket.BucketStart.Format(time.DateOnly)
		if timeHourUnit < 24 {
			x = bucket.BucketStart.Format("2006-01-02 15:04")
		}

		metadata := map[string]any{
			"bucketStart":  bucket.BucketStart,
			"bucketEnd":    bucketEnd,
			"timeHourUnit": timeHourUnit,
			"field":        fieldName,
		}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    bucket.BucketStart.Format(time.RFC3339),
			X:     x,
			Value: bucket.RoutineCount,
			Meta:  meta,
		}
	}

	return data, nil
}

/* ============================== Service Methods for Routine ============================== */

func (s *RoutineService) GetMyRoutineById(
	ctx context.Context, reqDto *dtos.GetMyRoutineByIdReqDto,
) (*dtos.GetMyRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	routine, exception := s.routineRepository.GetOneById(
		reqDto.Param.RoutineId,
		reqDto.ContextFields.UserId,
		[]schemas.RoutineRelation{
			schemas.RoutineRelation_RoutinesToTags,
			schemas.RoutineRelation_RoutinesToTasks,
			schemas.RoutineRelation_RoutinesToItems,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	tagIds := make([]uuid.UUID, len(routine.RoutinesToTags))
	for index, routineToTag := range routine.RoutinesToTags {
		tagIds[index] = routineToTag.TagId
	}
	taskIds := make([]uuid.UUID, len(routine.RoutinesToTasks))
	for index, routineToTask := range routine.RoutinesToTasks {
		taskIds[index] = routineToTask.TaskId
	}
	itemIds := make([]uuid.UUID, len(routine.RoutinesToItems))
	for index, routineToItem := range routine.RoutinesToItems {
		itemIds[index] = routineToItem.ItemId
	}

	return &dtos.GetMyRoutineByIdResDto{
		Id:               routine.Id,
		StationId:        routine.StationId,
		Title:            routine.Title,
		Description:      routine.Description,
		Status:           routine.Status,
		IsPinned:         routine.IsPinned,
		ScheduledStartAt: routine.ScheduledStartAt,
		ScheduledEndAt:   routine.ScheduledEndAt,
		Period:           routine.Period,
		Timezone:         routine.Timezone,
		DeletedAt:        routine.DeletedAt,
		UpdatedAt:        routine.UpdatedAt,
		CreatedAt:        routine.CreatedAt,
		TagIds:           tagIds,
		TaskIds:          taskIds,
		ItemIds:          itemIds,
	}, nil
}

func (s *RoutineService) GetMyRoutinesByStationId(
	ctx context.Context, reqDto *dtos.GetMyRoutinesByStationIdReqDto,
) (*dtos.GetMyRoutinesByStationIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	var routines []schemas.Routine
	query := db.Model(&schemas.Routine{}).
		Select(`"RoutineTable".*`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "RoutineTable".station_id`).
		Joins(`INNER JOIN "StationTable" station ON station.id = "RoutineTable".station_id AND station.deleted_at IS NULL`).
		Where(`"RoutineTable".station_id = ?`, reqDto.Param.StationId).
		Where("uts.user_id = ? AND uts.permission IN ?", reqDto.ContextFields.UserId, allowedPermissions).
		Scopes(s.routineScope.IncludePreloads(
			[]schemas.RoutineRelation{
				schemas.RoutineRelation_RoutinesToTags,
				schemas.RoutineRelation_RoutinesToTasks,
				schemas.RoutineRelation_RoutinesToItems,
			},
		))

	query = query.Scopes(s.routineScope.FilterOnlyDeleted(onlyDeleted))

	result := query.Order(`"RoutineTable".scheduled_start_at ASC`).
		Order(`"RoutineTable".scheduled_end_at ASC`).
		Order(`"RoutineTable".id ASC`).
		Find(&routines)
	if result.Error != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(result.Error)
	}

	resDto := make(dtos.GetMyRoutinesByStationIdResDto, len(routines))
	for index, routine := range routines {
		tagIds := make([]uuid.UUID, len(routine.RoutinesToTags))
		for index, routineToTag := range routine.RoutinesToTags {
			tagIds[index] = routineToTag.TagId
		}
		taskIds := make([]uuid.UUID, len(routine.RoutinesToTasks))
		for index, routineToTask := range routine.RoutinesToTasks {
			taskIds[index] = routineToTask.TaskId
		}
		itemIds := make([]uuid.UUID, len(routine.RoutinesToItems))
		for index, routineToItem := range routine.RoutinesToItems {
			itemIds[index] = routineToItem.ItemId
		}
		resDto[index] = struct {
			Id               uuid.UUID            "json:\"id\""
			StationId        uuid.UUID            "json:\"stationId\""
			Title            string               "json:\"title\""
			Status           enums.RoutineStatus  "json:\"status\""
			IsPinned         bool                 "json:\"isPinned\""
			ScheduledStartAt time.Time            "json:\"scheduledStartAt\""
			ScheduledEndAt   time.Time            "json:\"scheduledEndAt\""
			Period           *enums.RoutinePeriod "json:\"period\""
			Timezone         string               "json:\"timezone\""
			DeletedAt        *time.Time           "json:\"deletedAt\""
			UpdatedAt        time.Time            "json:\"updatedAt\""
			CreatedAt        time.Time            "json:\"createdAt\""
			TagIds           []uuid.UUID          "json:\"tagIds\""
			TaskIds          []uuid.UUID          "json:\"taskIds\""
			ItemIds          []uuid.UUID          "json:\"itemIds\""
		}{
			Id:               routine.Id,
			StationId:        routine.StationId,
			Title:            routine.Title,
			Status:           routine.Status,
			IsPinned:         routine.IsPinned,
			ScheduledStartAt: routine.ScheduledStartAt,
			ScheduledEndAt:   routine.ScheduledEndAt,
			Period:           routine.Period,
			Timezone:         routine.Timezone,
			DeletedAt:        routine.DeletedAt,
			UpdatedAt:        routine.UpdatedAt,
			CreatedAt:        routine.CreatedAt,
			TagIds:           tagIds,
			TaskIds:          taskIds,
			ItemIds:          itemIds,
		}
	}

	return &resDto, nil
}

func (s *RoutineService) GetAllMyRoutinesByTimeRange(
	ctx context.Context, reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto,
) (*dtos.GetAllMyRoutinesByTimeRangeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.From.Before(reqDto.Param.To) { // make sure from is before to
		return nil, exceptions.Routine.InvalidInput().WithOrigin(fmt.Errorf("from must be before to"))
	}
	if !util.IsTimeWithin(reqDto.Param.From, reqDto.Param.To, 360*24*time.Hour) { // make sure the time range is within 360 days which is approximate 1 year
		return nil, exceptions.Routine.QueriedTimeRangeTooLarge(reqDto.Param.From, reqDto.Param.To)
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	routines, exception := s.routineRepository.GetAllByTimeRange(
		reqDto.Param.From,
		reqDto.Param.To,
		reqDto.Param.StationIds,
		reqDto.ContextFields.UserId,
		[]schemas.RoutineRelation{
			schemas.RoutineRelation_RoutinesToTags,
			schemas.RoutineRelation_RoutinesToTasks,
			schemas.RoutineRelation_RoutinesToItems,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutinesByTimeRangeResDto, len(routines))
	for index, routine := range routines {
		tagIds := make([]uuid.UUID, len(routine.RoutinesToTags))
		for index, routineToTag := range routine.RoutinesToTags {
			tagIds[index] = routineToTag.TagId
		}
		taskIds := make([]uuid.UUID, len(routine.RoutinesToTasks))
		for index, routineToTask := range routine.RoutinesToTasks {
			taskIds[index] = routineToTask.TaskId
		}
		itemIds := make([]uuid.UUID, len(routine.RoutinesToItems))
		for index, routineToItem := range routine.RoutinesToItems {
			itemIds[index] = routineToItem.ItemId
		}
		resDto[index] = struct {
			Id               uuid.UUID            "json:\"id\""
			StationId        uuid.UUID            "json:\"stationId\""
			Title            string               "json:\"title\""
			Status           enums.RoutineStatus  "json:\"status\""
			IsPinned         bool                 "json:\"isPinned\""
			ScheduledStartAt time.Time            "json:\"scheduledStartAt\""
			ScheduledEndAt   time.Time            "json:\"scheduledEndAt\""
			Period           *enums.RoutinePeriod "json:\"period\""
			Timezone         string               "json:\"timezone\""
			DeletedAt        *time.Time           "json:\"deletedAt\""
			UpdatedAt        time.Time            "json:\"updatedAt\""
			CreatedAt        time.Time            "json:\"createdAt\""
			TagIds           []uuid.UUID          "json:\"tagIds\""
			TaskIds          []uuid.UUID          "json:\"taskIds\""
			ItemIds          []uuid.UUID          "json:\"itemIds\""
		}{
			Id:               routine.Id,
			StationId:        routine.StationId,
			Title:            routine.Title,
			Status:           routine.Status,
			IsPinned:         routine.IsPinned,
			ScheduledStartAt: routine.ScheduledStartAt,
			ScheduledEndAt:   routine.ScheduledEndAt,
			Period:           routine.Period,
			Timezone:         routine.Timezone,
			DeletedAt:        routine.DeletedAt,
			UpdatedAt:        routine.UpdatedAt,
			CreatedAt:        routine.CreatedAt,
			TagIds:           tagIds,
			TaskIds:          taskIds,
			ItemIds:          itemIds,
		}
	}

	return &resDto, nil
}

func (s *RoutineService) CreateRoutineByStationId(
	ctx context.Context, reqDto *dtos.CreateRoutineByStationIdReqDto,
) (*dtos.CreateRoutineByStationIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	newRoutineId, exception := s.routineRepository.CreateOneByStationId(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		inputs.CreateRoutineInput{
			Id:               reqDto.Body.Id,
			Title:            reqDto.Body.Title,
			Description:      reqDto.Body.Description,
			Status:           reqDto.Body.Status,
			IsPinned:         reqDto.Body.IsPinned,
			ScheduledStartAt: reqDto.Body.ScheduledStartAt,
			ScheduledEndAt:   reqDto.Body.ScheduledEndAt,
			Period:           reqDto.Body.Period,
			Timezone:         reqDto.Body.Timezone,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutineByStationIdResDto{
		Id:        *newRoutineId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) CreateRoutinesByStationIds(
	ctx context.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto,
) (*dtos.CreateRoutinesByStationIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkCreateRoutineInput, len(reqDto.Body.CreatedRoutines))
	for index, createdRoutine := range reqDto.Body.CreatedRoutines {
		input[index] = inputs.BulkCreateRoutineInput{
			Id:               createdRoutine.Id,
			StationId:        createdRoutine.StationId,
			Title:            createdRoutine.Title,
			Description:      createdRoutine.Description,
			Status:           createdRoutine.Status,
			IsPinned:         createdRoutine.IsPinned,
			ScheduledStartAt: createdRoutine.ScheduledStartAt,
			ScheduledEndAt:   createdRoutine.ScheduledEndAt,
			Period:           createdRoutine.Period,
			Timezone:         createdRoutine.Timezone,
		}
	}
	newRoutineIds, exception := s.routineRepository.BulkCreateManyByStationIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutinesByStationIdsResDto{
		Ids:       newRoutineIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) UpdateMyRoutineById(
	ctx context.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto,
) (*dtos.UpdateMyRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	updatedRoutine, exception := s.routineRepository.UpdateOneById(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRoutineInput{
			Values: inputs.UpdateRoutineInput{
				StationId:        reqDto.Body.Values.StationId,
				Title:            reqDto.Body.Values.Title,
				Description:      reqDto.Body.Values.Description,
				Status:           reqDto.Body.Values.Status,
				IsPinned:         reqDto.Body.Values.IsPinned,
				ScheduledStartAt: reqDto.Body.Values.ScheduledStartAt,
				ScheduledEndAt:   reqDto.Body.Values.ScheduledEndAt,
				Period:           reqDto.Body.Values.Period,
				Timezone:         reqDto.Body.Values.Timezone,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRoutineByIdResDto{
		UpdatedAt: updatedRoutine.UpdatedAt,
	}, nil
}

func (s *RoutineService) UpdateMyRoutinesByIds(
	ctx context.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto,
) (*dtos.UpdateMyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkUpdateRoutineInput, len(reqDto.Body.UpdatedRoutines))
	for index, updatedRoutine := range reqDto.Body.UpdatedRoutines {
		input[index] = inputs.BulkUpdateRoutineInput{
			Id: updatedRoutine.RoutineId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateRoutineInput]{
				Values: inputs.UpdateRoutineInput{
					StationId:        updatedRoutine.Values.StationId,
					Title:            updatedRoutine.Values.Title,
					Description:      updatedRoutine.Values.Description,
					Status:           updatedRoutine.Values.Status,
					IsPinned:         updatedRoutine.Values.IsPinned,
					ScheduledStartAt: updatedRoutine.Values.ScheduledStartAt,
					ScheduledEndAt:   updatedRoutine.Values.ScheduledEndAt,
					Period:           updatedRoutine.Values.Period,
					Timezone:         updatedRoutine.Values.Timezone,
				},
				SetNull: updatedRoutine.SetNull,
			},
		}
	}
	exception := s.routineRepository.BulkUpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyRoutinesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) LinkRoutineTagById(
	ctx context.Context, reqDto *dtos.LinkRoutineTagByIdReqDto,
) (*dtos.LinkRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.routineTagRepository.HasPermission(
		reqDto.Body.RoutineTagId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.RoutineTag.NoPermission("get the routine tag")
	}

	var newRoutinesToTags schemas.RoutinesToTags
	newRoutinesToTags.RoutineId = reqDto.Body.RoutineId
	newRoutinesToTags.TagId = reqDto.Body.RoutineTagId

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Where("routine_id = ? AND tag_id = ?", newRoutinesToTags.RoutineId, newRoutinesToTags.TagId).
			Delete(&schemas.RoutinesToTags{})
	} else {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Create(&newRoutinesToTags)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTags().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineTagByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineTagsByIds(
	ctx context.Context, reqDto *dtos.BulkLinkRoutineTagsByIdsReqDto,
) (*dtos.BulkLinkRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	isRoutineExist := make(map[uuid.UUID]bool)
	isRoutineTagExist := make(map[uuid.UUID]bool)
	var routineIds []uuid.UUID
	var routineTagIds []uuid.UUID
	for _, linkedRoutineAndTag := range reqDto.Body.LinkedRoutinesAndTags {
		if !isRoutineExist[linkedRoutineAndTag.RoutineId] {
			isRoutineExist[linkedRoutineAndTag.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndTag.RoutineId)
		}
		if !isRoutineTagExist[linkedRoutineAndTag.RoutineTagId] {
			isRoutineTagExist[linkedRoutineAndTag.RoutineTagId] = true
			routineTagIds = append(routineTagIds, linkedRoutineAndTag.RoutineTagId)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isRoutineTagValid := make(map[uuid.UUID]bool)
	validRoutineTags, exception := s.routineTagRepository.CheckPermissionsAndGetManyByIds(
		routineTagIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutineTag := range validRoutineTags {
		isRoutineTagValid[validRoutineTag.Id] = true
	}

	var newRoutinesToTags []schemas.RoutinesToTags
	for _, linkedRoutineAndTag := range reqDto.Body.LinkedRoutinesAndTags {
		if !isRoutineValid[linkedRoutineAndTag.RoutineId] ||
			!isRoutineTagValid[linkedRoutineAndTag.RoutineTagId] {
			continue
		}
		newRoutinesToTags = append(newRoutinesToTags, schemas.RoutinesToTags{
			RoutineId: linkedRoutineAndTag.RoutineId,
			TagId:     linkedRoutineAndTag.RoutineTagId,
		})
	}
	if len(newRoutinesToTags) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToTags))
	for index, newRoutineToTag := range newRoutinesToTags {
		values[index] = []any{newRoutineToTag.RoutineId, newRoutineToTag.TagId}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Where("(routine_id, tag_id) IN ?", values).
			Delete(&schemas.RoutinesToTags{})
	} else {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Create(&newRoutinesToTags)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTags().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineTagsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) LinkRoutineTaskById(
	ctx context.Context, reqDto *dtos.LinkRoutineTaskByIdReqDto,
) (*dtos.LinkRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.routineTaskRepository.HasPermission(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.RoutineTask.NoPermission("get the routine task")
	}

	var newRoutinesToTasks schemas.RoutinesToTasks
	newRoutinesToTasks.RoutineId = reqDto.Body.RoutineId
	newRoutinesToTasks.TaskId = reqDto.Body.RoutineTaskId

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Where("routine_id = ? AND task_id = ?", newRoutinesToTasks.RoutineId, newRoutinesToTasks.TaskId).
			Delete(&schemas.RoutinesToTasks{})
	} else {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Create(&newRoutinesToTasks)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTasks().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineTaskByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineTasksByIds(
	ctx context.Context, reqDto *dtos.BulkLinkRoutineTasksByIdsReqDto,
) (*dtos.BulkLinkRoutineTasksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	isRoutineExist := make(map[uuid.UUID]bool)
	isRoutineTaskExist := make(map[uuid.UUID]bool)
	var routineIds []uuid.UUID
	var routineTaskIds []uuid.UUID
	for _, linkedRoutineAndTask := range reqDto.Body.LinkedRoutinesAndTasks {
		if !isRoutineExist[linkedRoutineAndTask.RoutineId] {
			isRoutineExist[linkedRoutineAndTask.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndTask.RoutineId)
		}
		if !isRoutineTaskExist[linkedRoutineAndTask.RoutineTaskId] {
			isRoutineTaskExist[linkedRoutineAndTask.RoutineTaskId] = true
			routineTaskIds = append(routineTaskIds, linkedRoutineAndTask.RoutineTaskId)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isRoutineTaskValid := make(map[uuid.UUID]bool)
	validRoutineTasks, exception := s.routineTaskRepository.CheckPermissionsAndGetManyByIds(
		routineTaskIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutineTag := range validRoutineTasks {
		isRoutineTaskValid[validRoutineTag.Id] = true
	}

	var newRoutinesToTasks []schemas.RoutinesToTasks
	for _, linkedRoutineAndTask := range reqDto.Body.LinkedRoutinesAndTasks {
		if !isRoutineValid[linkedRoutineAndTask.RoutineId] ||
			!isRoutineTaskValid[linkedRoutineAndTask.RoutineTaskId] {
			continue
		}
		newRoutinesToTasks = append(newRoutinesToTasks, schemas.RoutinesToTasks{
			RoutineId: linkedRoutineAndTask.RoutineId,
			TaskId:    linkedRoutineAndTask.RoutineTaskId,
		})
	}
	if len(newRoutinesToTasks) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToTasks))
	for index, newRoutineToTask := range newRoutinesToTasks {
		values[index] = []any{newRoutineToTask.RoutineId, newRoutineToTask.TaskId}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Where("(routine_id, task_id) IN ?", values).
			Delete(&schemas.RoutinesToTasks{})
	} else {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Create(&newRoutinesToTasks)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTasks().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineTasksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) LinkRoutineItemById(
	ctx context.Context, reqDto *dtos.LinkRoutineItemByIdReqDto,
) (*dtos.LinkRoutineItemByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.itemRepository.HasPermission(
		reqDto.Body.ItemId,
		reqDto.Body.ItemType,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Item.NoPermission("get the item")
	}

	var newRoutinesToItems schemas.RoutinesToItems
	newRoutinesToItems.RoutineId = reqDto.Body.RoutineId
	newRoutinesToItems.ItemId = reqDto.Body.ItemId
	newRoutinesToItems.ItemType = reqDto.Body.ItemType

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Where(
				"routine_id = ? AND item_id = ? AND item_type = ?",
				newRoutinesToItems.RoutineId,
				newRoutinesToItems.ItemId,
				newRoutinesToItems.ItemType,
			).
			Delete(&schemas.RoutinesToItems{})
	} else {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Create(&newRoutinesToItems)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkItems().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineItemByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineItemsByIds(
	ctx context.Context, reqDto *dtos.BulkLinkRoutineItemsByIdsReqDto,
) (*dtos.BulkLinkRoutineItemsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	isRoutineExist := make(map[uuid.UUID]bool)
	isItemExist := make(map[types.Pair[uuid.UUID, enums.ItemType]]bool)
	var routineIds []uuid.UUID
	var itemIdentities []types.Pair[uuid.UUID, enums.ItemType]
	for _, linkedRoutineAndItem := range reqDto.Body.LinkedRoutinesAndItems {
		if !isRoutineExist[linkedRoutineAndItem.RoutineId] {
			isRoutineExist[linkedRoutineAndItem.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndItem.RoutineId)
		}
		itemIdentity := types.Pair[uuid.UUID, enums.ItemType]{
			First:  linkedRoutineAndItem.ItemId,
			Second: linkedRoutineAndItem.ItemType,
		}
		if !isItemExist[itemIdentity] {
			isItemExist[itemIdentity] = true
			itemIdentities = append(itemIdentities, itemIdentity)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isItemValid := make(map[types.Pair[uuid.UUID, enums.ItemType]]bool)
	validItems, exception := s.itemRepository.CheckPermissionsAndGetManyByIds(
		itemIdentities,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validItem := range validItems {
		isItemValid[types.Pair[uuid.UUID, enums.ItemType]{
			First:  validItem.Id,
			Second: validItem.Type,
		}] = true
	}

	var newRoutinesToItems []schemas.RoutinesToItems
	for _, linkedRoutineAndItem := range reqDto.Body.LinkedRoutinesAndItems {
		itemIdentity := types.Pair[uuid.UUID, enums.ItemType]{
			First:  linkedRoutineAndItem.ItemId,
			Second: linkedRoutineAndItem.ItemType,
		}
		if !isRoutineValid[linkedRoutineAndItem.RoutineId] ||
			!isItemValid[itemIdentity] {
			continue
		}
		newRoutinesToItems = append(newRoutinesToItems, schemas.RoutinesToItems{
			RoutineId: linkedRoutineAndItem.RoutineId,
			ItemId:    linkedRoutineAndItem.ItemId,
			ItemType:  linkedRoutineAndItem.ItemType,
		})
	}
	if len(newRoutinesToItems) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToItems))
	for index, newRoutineToItem := range newRoutinesToItems {
		values[index] = []any{newRoutineToItem.RoutineId, newRoutineToItem.ItemId, newRoutineToItem.ItemType}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Where("(routine_id, item_id, item_type) IN ?", values).
			Delete(&schemas.RoutinesToItems{})
	} else {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Create(&newRoutinesToItems)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkItems().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineItemsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) RestoreMyRoutineById(
	ctx context.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto,
) (*dtos.RestoreMyRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredRoutine, exception := s.routineRepository.RestoreSoftDeletedOneById(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyRoutineByIdResDto{
		Id:               restoredRoutine.Id,
		StationId:        restoredRoutine.StationId,
		Title:            restoredRoutine.Title,
		Description:      restoredRoutine.Description,
		Status:           restoredRoutine.Status,
		IsPinned:         restoredRoutine.IsPinned,
		ScheduledStartAt: restoredRoutine.ScheduledStartAt,
		ScheduledEndAt:   restoredRoutine.ScheduledEndAt,
		Period:           restoredRoutine.Period,
		Timezone:         restoredRoutine.Timezone,
		DeletedAt:        restoredRoutine.DeletedAt,
		UpdatedAt:        restoredRoutine.UpdatedAt,
		CreatedAt:        restoredRoutine.CreatedAt,
	}, nil
}

func (s *RoutineService) RestoreMyRoutinesByIds(
	ctx context.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto,
) (*dtos.RestoreMyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredRoutines, exception := s.routineRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.RoutineIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := dtos.RestoreMyRoutinesByIdsResDto{}
	for _, restoredRoutine := range restoredRoutines {
		resDto = append(resDto, dtos.RestoreMyRoutineByIdResDto{
			Id:               restoredRoutine.Id,
			StationId:        restoredRoutine.StationId,
			Title:            restoredRoutine.Title,
			Description:      restoredRoutine.Description,
			Status:           restoredRoutine.Status,
			IsPinned:         restoredRoutine.IsPinned,
			ScheduledStartAt: restoredRoutine.ScheduledStartAt,
			ScheduledEndAt:   restoredRoutine.ScheduledEndAt,
			Period:           restoredRoutine.Period,
			Timezone:         restoredRoutine.Timezone,
			DeletedAt:        restoredRoutine.DeletedAt,
			UpdatedAt:        restoredRoutine.UpdatedAt,
			CreatedAt:        restoredRoutine.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *RoutineService) DeleteMyRoutineById(
	ctx context.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto,
) (*dtos.DeleteMyRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineRepository.SoftDeleteOneById(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) DeleteMyRoutinesByIds(
	ctx context.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto,
) (*dtos.DeleteMyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineRepository.SoftDeleteManyByIds(
		reqDto.Body.RoutineIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteMyRoutineById(
	ctx context.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto,
) (*dtos.HardDeleteMyRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineRepository.HardDeleteOneById(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteMyRoutinesByIds(
	ctx context.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto,
) (*dtos.HardDeleteMyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.routineRepository.HardDeleteManyByIds(
		reqDto.Body.RoutineIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for Charts ============================== */

func (s *RoutineService) VisualizeMyRoutineStatusCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineStatusCountReqDto,
) (*dtos.VisualizeMyRoutineStatusCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	var counts struct {
		ScheduledCount  int64 `gorm:"column:scheduled_count;"`
		InProgressCount int64 `gorm:"column:in_progress_count;"`
		CompletedCount  int64 `gorm:"column:completed_count;"`
		OverDueCount    int64 `gorm:"column:over_due_count;"`
	}
	result := db.Model(&schemas.Routine{}).
		Select(`
			COUNT(*) FILTER (WHERE status = ?) as scheduled_count,
			COUNT(*) FILTER (WHERE status = ?) as in_progress_count,
			COUNT(*) FILTER (WHERE status = ?) as completed_count,
			COUNT(*) FILTER (WHERE status = ?) as over_due_count
		`,
			enums.RoutineStatus_Scheduled,
			enums.RoutineStatus_InProgress,
			enums.RoutineStatus_Completed,
			enums.RoutineStatus_OverDue,
		).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "RoutineTable".station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Where(`"RoutineTable".deleted_at IS NULL`).
		Scan(&counts)
	if err := result.Error; err != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(err)
	}

	scheduledRoutineMetadata := map[string]string{"status": "scheduled"}
	scheduledRoutineMeta, err := json.Marshal(scheduledRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(scheduledRoutineMetadata)
	}

	inProgressRoutineMetadata := map[string]string{"status": "inProgress"}
	inProgressRoutineMeta, err := json.Marshal(inProgressRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(inProgressRoutineMetadata)
	}

	completedRoutineMetadata := map[string]string{"status": "completed"}
	completedRoutineMeta, err := json.Marshal(completedRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(completedRoutineMetadata)
	}

	overDueRoutineMetadata := map[string]string{"status": "overDue"}
	overDueRoutineMeta, err := json.Marshal(overDueRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(overDueRoutineMetadata)
	}

	return &dtos.VisualizeMyRoutineStatusCountResDto{
		Data: []dtos.TwoDimensionalDatum[int64]{
			dtos.TwoDimensionalDatum[int64]{
				Id:    "scheduled-routine-count",
				X:     "Scheduled Routine Count",
				Value: counts.ScheduledCount,
				Meta:  scheduledRoutineMeta,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "in-progress-routine-count",
				X:     "In Progress Routine Count",
				Value: counts.InProgressCount,
				Meta:  inProgressRoutineMeta,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "completed-routine-count",
				X:     "Completed Routine Count",
				Value: counts.CompletedCount,
				Meta:  completedRoutineMeta,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "over-due-routine-count",
				X:     "Over Due Routine Count",
				Value: counts.OverDueCount,
				Meta:  overDueRoutineMeta,
			},
		},
	}, nil
}

func (s *RoutineService) VisualizeMyRoutinePeriodCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutinePeriodCountReqDto,
) (*dtos.VisualizeMyRoutinePeriodCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	var counts struct {
		DailyCount   int64 `gorm:"column:daily_count;"`
		WeeklyCount  int64 `gorm:"column:weekly_count;"`
		MonthlyCount int64 `gorm:"column:monthly_count;"`
	}
	result := db.Model(&schemas.Routine{}).
		Select(`
			COUNT(*) FILTER (WHERE period = ?) as daily_count,
			COUNT(*) FILTER (WHERE period = ?) as weekly_count,
			COUNT(*) FILTER (WHERE period = ?) as monthly_count
		`,
			enums.RoutinePeriod_Daily,
			enums.RoutinePeriod_Weekly,
			enums.RoutinePeriod_Monthly,
		).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "RoutineTable".station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Where(`"RoutineTable".deleted_at IS NULL`).
		Scan(&counts)
	if err := result.Error; err != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(err)
	}

	dailyRoutineMetadata := map[string]string{"period": "daily"}
	dailyRoutineMeta, err := json.Marshal(dailyRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(dailyRoutineMetadata)
	}

	weeklyRoutineMetadata := map[string]string{"period": "daily"}
	weeklyRoutineMeta, err := json.Marshal(weeklyRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(weeklyRoutineMetadata)
	}

	monthlyRoutineMetadata := map[string]string{"period": "daily"}
	monthlyRoutineMeta, err := json.Marshal(monthlyRoutineMetadata)
	if err != nil {
		return nil, exceptions.Routine.FailedToMarshalData(monthlyRoutineMetadata)
	}

	return &dtos.VisualizeMyRoutinePeriodCountResDto{
		Data: []dtos.TwoDimensionalDatum[int64]{
			dtos.TwoDimensionalDatum[int64]{
				Id:    "daily-routine-count",
				X:     "Daily Routine Count",
				Value: counts.DailyCount,
				Meta:  dailyRoutineMeta,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "weekly-routine-count",
				X:     "Weekly Routine Count",
				Value: counts.WeeklyCount,
				Meta:  weeklyRoutineMeta,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "monthly-routine-count",
				X:     "Monthly Routine Count",
				Value: counts.MonthlyCount,
				Meta:  monthlyRoutineMeta,
			},
		},
	}, nil
}

func (s *RoutineService) VisualizeMyRoutineScheduledStartAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineScheduledStartAtCountReqDto,
) (*dtos.VisualizeMyRoutineScheduledStartAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.Routine.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.Routine.QueriedTimeRangeTooLarge(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt)
	}

	data, exception := s.visualizeMyRoutineTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"scheduled_start_at",
		"scheduledStartAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineScheduledStartAtCountResDto{
		Data: data,
	}, nil
}

func (s *RoutineService) VisualizeMyRoutineScheduledEndAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineScheduledEndAtCountReqDto,
) (*dtos.VisualizeMyRoutineScheduledEndAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Routine.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.Routine.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.Routine.QueriedTimeRangeTooLarge(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt)
	}

	data, exception := s.visualizeMyRoutineTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"scheduled_end_at",
		"scheduledEndAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineScheduledEndAtCountResDto{
		Data: data,
	}, nil
}

/* ============================== Service Methods for GraphQL Routine ============================== */

func (s *RoutineService) SearchPrivateRoutines(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineInput,
) (*gqlmodels.SearchRoutineConnection, *exceptions.Exception) {
	type PrivateRoutine struct {
		schemas.Routine
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

	query := db.Model(&schemas.Routine{}).
		Select(`"RoutineTable".*, uts.permission AS permission`).
		Joins(`LEFT JOIN "UsersToStationsTable" uts ON "RoutineTable".station_id = uts.station_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Where(`"RoutineTable".deleted_at IS NULL`)

	if len(gqlInput.StationIds) > 0 {
		query = query.Where(
			`"RoutineTable".station_id IN ?`,
			gqlInput.StationIds,
		)
	}

	if len(gqlInput.TagIds) > 0 {
		if !s.routineTagRepository.HavePermissions(
			gqlInput.TagIds,
			userId,
			allowedPermissions,
			options.WithDB(db),
		) {
			return nil, exceptions.RoutineTag.NoPermission("filter routines by these tags")
		}

		subQuery := db.
			Session(&gorm.Session{NewDB: true}).
			Model(&schemas.RoutinesToTags{}).
			Select("1").
			Where(`"RoutinesToTagsTable".routine_id = "RoutineTable".id`).
			Where(`"RoutinesToTagsTable".tag_id IN ?`, gqlInput.TagIds)

		query = query.Where("EXISTS (?)", subQuery)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"title ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRoutineCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(
			`"RoutineTable".id > ?`,
			searchCursor.Fields.ID,
		)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRoutineSortByTitle:
			query.Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByStatus:
			query.Order("status " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByScheduledStartAt:
			query.Order("scheduled_start_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByScheduledEndAt:
			query.Order("scheduled_end_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByPeriod:
			query.Order("period " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("title " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("title " + cending).
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

	var routines []PrivateRoutine
	if err := query.Scopes(s.routineScope.IncludePreloads(
		[]schemas.RoutineRelation{
			schemas.RoutineRelation_RoutinesToTags,
			schemas.RoutineRelation_RoutinesToTasks,
			schemas.RoutineRelation_RoutinesToItems,
		},
	)).Find(&routines).Error; err != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(err)
	}

	hasNextPage := len(routines) > limit
	searchEdges := make([]*gqlmodels.SearchRoutineEdge, len(routines))

	for index, routine := range routines {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRoutineCursorFields]{
			Fields: gqlmodels.SearchRoutineCursorFields{
				ID: routine.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		finalRoutine := routine.Routine.ToPrivateSearchableRoutine()
		for _, routineToTag := range routine.Routine.RoutinesToTags {
			finalRoutine.TagIds = append(finalRoutine.TagIds, routineToTag.TagId)
		}
		for _, routineToTask := range routine.Routine.RoutinesToTasks {
			finalRoutine.TaskIds = append(finalRoutine.TaskIds, routineToTask.TaskId)
		}
		for _, routineToItem := range routine.Routine.RoutinesToItems {
			finalRoutine.ItemIds = append(finalRoutine.ItemIds, routineToItem.ItemId)
		}
		searchEdges[index] = &gqlmodels.SearchRoutineEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                finalRoutine,
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: false,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchRoutineConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
