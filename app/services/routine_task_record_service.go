package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
)

type RoutineTaskRecordServiceInterface interface {
	GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx context.Context, reqDto *dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto) (*dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskRecordStatusCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto) (*dtos.VisualizeMyRoutineTaskRecordStatusCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskRecordPurposeCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto) (*dtos.VisualizeMyRoutineTaskRecordPurposeCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskRecordScheduledAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto) (*dtos.VisualizeMyRoutineTaskRecordScheduledAtCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto) (*dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto) (*dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountResDto, *exceptions.Exception)
	SearchPrivateRoutineTaskRecords(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTaskRecordInput) (*gqlmodels.SearchRoutineTaskRecordConnection, *exceptions.Exception)
}

type RoutineTaskRecordService struct {
	db                          *gorm.DB
	routineTaskRecordRepository repositories.RoutineTaskRecordRepositoryInterface
}

func NewRoutineTaskRecordService(
	db *gorm.DB,
	routineTaskRecordRepository repositories.RoutineTaskRecordRepositoryInterface,
) RoutineTaskRecordServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	if routineTaskRecordRepository == nil {
		routineTaskRecordRepository = repositories.NewRoutineTaskRecordRepository(scopes.NewRoutineTaskRecordScope())
	}

	return &RoutineTaskRecordService{
		db:                          db,
		routineTaskRecordRepository: routineTaskRecordRepository,
	}
}

func (s *RoutineTaskRecordService) visualizeMyRoutineTaskRecordTimeCount(
	ctx context.Context,
	userId uuid.UUID,
	permission enums.AccessControlPermission,
	routineTaskIds []uuid.UUID,
	timeHourUnit int,
	queryRangeStartedAt time.Time,
	queryRangeEndedAt time.Time,
	columnName string,
	fieldName string,
) ([]dtos.TwoDimensionalDatum[int64], *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	var buckets []struct {
		BucketStart            time.Time `gorm:"column:bucket_start;"`
		RoutineTaskRecordCount int64     `gorm:"column:routine_task_record_count;"`
	}

	recordJoin := `LEFT JOIN "RoutineTaskRecordTable" routine_task_record
		ON routine_task_record.` + columnName + ` >= buckets.bucket_start
		AND routine_task_record.` + columnName + ` < buckets.bucket_start + ?::integer * interval '1 hour'`
	recordJoinArgs := []any{timeHourUnit}
	if len(routineTaskIds) > 0 {
		recordJoin += ` AND routine_task_record.routine_task_id IN ?`
		recordJoinArgs = append(recordJoinArgs, routineTaskIds)
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
			COUNT(uts.station_id) AS routine_task_record_count
		`).
		Joins(recordJoin, recordJoinArgs...).
		Joins(`LEFT JOIN "RoutineTaskTable" routine_task ON routine_task.id = routine_task_record.routine_task_id`).
		Joins(`LEFT JOIN "RoutineTable" routine ON routine.id = routine_task.routine_id AND routine.deleted_at IS NULL`).
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
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(buckets))
	for index, bucket := range buckets {
		bucketEnd := bucket.BucketStart.Add(time.Duration(timeHourUnit) * time.Hour)
		x := bucket.BucketStart.Format(time.DateOnly)
		if timeHourUnit < 24 {
			x = bucket.BucketStart.Format("2006-01-02 15:04")
		}

		metadata := map[string]any{
			"bucketStart":    bucket.BucketStart,
			"bucketEnd":      bucketEnd,
			"timeHourUnit":   timeHourUnit,
			"field":          fieldName,
			"routineTaskIds": routineTaskIds,
		}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    bucket.BucketStart.Format(time.RFC3339),
			X:     x,
			Value: bucket.RoutineTaskRecordCount,
			Meta:  meta,
		}
	}

	return data, nil
}

func (s *RoutineTaskRecordService) GetAllMyRoutineTaskRecordsByRoutineTaskId(
	ctx context.Context, reqDto *dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto,
) (*dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	routineTaskRecords, exception := s.routineTaskRecordRepository.GetAllByRoutineTaskId(
		reqDto.Param.RoutineTaskId,
		reqDto.ContextFields.UserId,
		reqDto.Param.Limit,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdResDto, len(routineTaskRecords))
	for index, routineTaskRecord := range routineTaskRecords {
		resDto[index] = struct {
			Id              uuid.UUID                         "json:\"id\""
			RoutineTaskId   uuid.UUID                         "json:\"routineTaskId\""
			Purpose         enums.RoutineTaskPurpose          "json:\"purpose\""
			Status          enums.RoutineTaskRecordStatus     "json:\"status\""
			ErrorCode       *enums.RoutineTaskRecordErrorCode "json:\"errorCode\""
			ErrorReason     *string                           "json:\"errorReason\""
			CostUnit        int64                             "json:\"costUnit\""
			TotalAttempts   int64                             "json:\"totalAttempts\""
			ScheduledAt     time.Time                         "json:\"scheduledAt\""
			ActualStartedAt *time.Time                        "json:\"actualStartedAt\""
			ActualEndedAt   *time.Time                        "json:\"actualEndedAt\""
			UpdatedAt       time.Time                         "json:\"updatedAt\""
			CreatedAt       time.Time                         "json:\"createdAt\""
		}{
			Id:              routineTaskRecord.Id,
			RoutineTaskId:   routineTaskRecord.RoutineTaskId,
			Purpose:         routineTaskRecord.Purpose,
			Status:          routineTaskRecord.Status,
			ErrorCode:       routineTaskRecord.ErrorCode,
			ErrorReason:     routineTaskRecord.ErrorReason,
			CostUnit:        routineTaskRecord.CostUnit,
			TotalAttempts:   routineTaskRecord.TotalAttempts,
			ScheduledAt:     routineTaskRecord.ScheduledAt,
			ActualStartedAt: routineTaskRecord.ActualStartedAt,
			ActualEndedAt:   routineTaskRecord.ActualEndedAt,
			UpdatedAt:       routineTaskRecord.UpdatedAt,
			CreatedAt:       routineTaskRecord.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *RoutineTaskRecordService) VisualizeMyRoutineTaskRecordStatusCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto,
) (*dtos.VisualizeMyRoutineTaskRecordStatusCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	var rows []struct {
		Status                 enums.RoutineTaskRecordStatus `gorm:"column:status;"`
		RoutineTaskRecordCount int64                         `gorm:"column:routine_task_record_count;"`
	}

	query := db.Model(&schemas.RoutineTaskRecord{}).
		Select(`"RoutineTaskRecordTable".status AS status, COUNT(*) AS routine_task_record_count`).
		Joins(`INNER JOIN "RoutineTaskTable" routine_task ON routine_task.id = "RoutineTaskRecordTable".routine_task_id`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = routine_task.routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission)
	if len(reqDto.Param.RoutineTaskIds) > 0 {
		query = query.Where(`"RoutineTaskRecordTable".routine_task_id IN ?`, reqDto.Param.RoutineTaskIds)
	}

	result := query.Group(`"RoutineTaskRecordTable".status`).Scan(&rows)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	counts := make(map[enums.RoutineTaskRecordStatus]int64, len(rows))
	for _, row := range rows {
		counts[row.Status] = row.RoutineTaskRecordCount
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(enums.AllRoutineTaskRecordStatuses))
	for index, status := range enums.AllRoutineTaskRecordStatuses {
		metadata := map[string]any{"status": status.String(), "routineTaskIds": reqDto.Param.RoutineTaskIds}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    status.String() + "-routine-task-record-count",
			X:     status.String() + " Routine Task Record Count",
			Value: counts[status],
			Meta:  meta,
		}
	}

	return &dtos.VisualizeMyRoutineTaskRecordStatusCountResDto{Data: data}, nil
}

func (s *RoutineTaskRecordService) VisualizeMyRoutineTaskRecordPurposeCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto,
) (*dtos.VisualizeMyRoutineTaskRecordPurposeCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	var rows []struct {
		Purpose                enums.RoutineTaskPurpose `gorm:"column:purpose;"`
		RoutineTaskRecordCount int64                    `gorm:"column:routine_task_record_count;"`
	}

	query := db.Model(&schemas.RoutineTaskRecord{}).
		Select(`"RoutineTaskRecordTable".purpose AS purpose, COUNT(*) AS routine_task_record_count`).
		Joins(`INNER JOIN "RoutineTaskTable" routine_task ON routine_task.id = "RoutineTaskRecordTable".routine_task_id`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = routine_task.routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission)
	if len(reqDto.Param.RoutineTaskIds) > 0 {
		query = query.Where(`"RoutineTaskRecordTable".routine_task_id IN ?`, reqDto.Param.RoutineTaskIds)
	}

	result := query.Group(`"RoutineTaskRecordTable".purpose`).Scan(&rows)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	counts := make(map[enums.RoutineTaskPurpose]int64, len(rows))
	for _, row := range rows {
		counts[row.Purpose] = row.RoutineTaskRecordCount
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(enums.AllRoutineTaskPurposes))
	for index, purpose := range enums.AllRoutineTaskPurposes {
		metadata := map[string]any{"purpose": purpose.String(), "routineTaskIds": reqDto.Param.RoutineTaskIds}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    purpose.String() + "-routine-task-record-count",
			X:     purpose.String() + " Routine Task Record Count",
			Value: counts[purpose],
			Meta:  meta,
		}
	}

	return &dtos.VisualizeMyRoutineTaskRecordPurposeCountResDto{Data: data}, nil
}

func (s *RoutineTaskRecordService) VisualizeMyRoutineTaskRecordScheduledAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskRecordScheduledAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskRecordTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.RoutineTaskIds,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"scheduled_at",
		"scheduledAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskRecordScheduledAtCountResDto{Data: data}, nil
}

func (s *RoutineTaskRecordService) VisualizeMyRoutineTaskRecordActualStartedAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskRecordTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.RoutineTaskIds,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_started_at",
		"actualStartedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountResDto{Data: data}, nil
}

func (s *RoutineTaskRecordService) VisualizeMyRoutineTaskRecordActualEndedAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskRecordTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.RoutineTaskIds,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_ended_at",
		"actualEndedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountResDto{Data: data}, nil
}

func (s *RoutineTaskRecordService) SearchPrivateRoutineTaskRecords(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTaskRecordInput,
) (*gqlmodels.SearchRoutineTaskRecordConnection, *exceptions.Exception) {
	type PrivateRoutineTaskRecord struct {
		schemas.RoutineTaskRecord
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	query := db.Model(&schemas.RoutineTaskRecord{}).
		Select(`"RoutineTaskRecordTable".*, uts.permission AS permission`).
		Joins(`INNER JOIN "RoutineTaskTable" routine_task ON routine_task.id = "RoutineTaskRecordTable".routine_task_id`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = routine_task.routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)

	if len(gqlInput.RoutineTaskIds) > 0 {
		query = query.Where(`"RoutineTaskRecordTable".routine_task_id IN ?`, gqlInput.RoutineTaskIds)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			`"RoutineTaskRecordTable".purpose::text ILIKE ?
				OR "RoutineTaskRecordTable".status::text ILIKE ?
				OR "RoutineTaskRecordTable".error_code::text ILIKE ?
				OR "RoutineTaskRecordTable".error_reason ILIKE ?`,
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchRoutineTaskRecordCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"RoutineTaskRecordTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRoutineTaskRecordSortByPurpose:
			query = query.Order(`"RoutineTaskRecordTable".purpose ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByStatus:
			query = query.Order(`"RoutineTaskRecordTable".status ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByCostUnit:
			query = query.Order(`"RoutineTaskRecordTable".cost_unit ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByTotalAttempts:
			query = query.Order(`"RoutineTaskRecordTable".total_attempts ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByScheduledAt:
			query = query.Order(`"RoutineTaskRecordTable".scheduled_at ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByActualStartedAt:
			query = query.Order(`"RoutineTaskRecordTable".actual_started_at ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByActualEndedAt:
			query = query.Order(`"RoutineTaskRecordTable".actual_ended_at ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByLastUpdate:
			query = query.Order(`"RoutineTaskRecordTable".updated_at ` + cending).
				Order(`"RoutineTaskRecordTable".created_at ` + cending)
		case gqlmodels.SearchRoutineTaskRecordSortByCreatedAt:
			query = query.Order(`"RoutineTaskRecordTable".created_at ` + cending)
		default:
			query = query.Order(`"RoutineTaskRecordTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var routineTaskRecords []PrivateRoutineTaskRecord
	if err := query.Find(&routineTaskRecords).Error; err != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	hasNextPage := len(routineTaskRecords) > limit
	searchEdges := make([]*gqlmodels.SearchRoutineTaskRecordEdge, len(routineTaskRecords))

	for index, routineTaskRecord := range routineTaskRecords {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRoutineTaskRecordCursorFields]{
			Fields: gqlmodels.SearchRoutineTaskRecordCursorFields{
				ID: routineTaskRecord.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRoutineTaskRecordEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                routineTaskRecord.RoutineTaskRecord.ToPrivateRoutineTaskRecord(),
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

	return &gqlmodels.SearchRoutineTaskRecordConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
