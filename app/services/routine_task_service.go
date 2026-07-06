package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
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
)

type RoutineTaskServiceInterface interface {
	GetMyRoutineTaskById(ctx context.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto) (*dtos.GetMyRoutineTaskByIdResDto, *exceptions.Exception)
	GetAllMyRoutineTasksByRoutineIds(ctx context.Context, reqDto *dtos.GetAllMyRoutineTasksByRoutineIdsReqDto) (*dtos.GetAllMyRoutineTasksByRoutineIdsResDto, *exceptions.Exception)
	GetAllMyRoutineTasks(ctx context.Context, reqDto *dtos.GetAllMyRoutineTasksReqDto) (*dtos.GetAllMyRoutineTasksResDto, *exceptions.Exception)
	CreateRoutineTaskByRoutineId(ctx context.Context, reqDto *dtos.CreateRoutineTaskByRoutineIdReqDto) (*dtos.CreateRoutineTaskByRoutineIdResDto, *exceptions.Exception)
	UpdateMyRoutineTaskById(ctx context.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto) (*dtos.UpdateMyRoutineTaskByIdResDto, *exceptions.Exception)
	PauseMyRoutineTaskById(ctx context.Context, reqDto *dtos.PauseMyRoutineTaskByIdReqDto) (*dtos.PauseMyRoutineTaskByIdResDto, *exceptions.Exception)
	ResumeMyRoutineTaskById(ctx context.Context, reqDto *dtos.ResumeMyRoutineTaskByIdReqDto) (*dtos.ResumeMyRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTaskById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto) (*dtos.HardDeleteMyRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTasksByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto) (*dtos.HardDeleteMyRoutineTasksByIdsResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskStatusCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskStatusCountReqDto) (*dtos.VisualizeMyRoutineTaskStatusCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskPurposeCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskPurposeCountReqDto) (*dtos.VisualizeMyRoutineTaskPurposeCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskScheduledAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto) (*dtos.VisualizeMyRoutineTaskScheduledAtCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto) (*dtos.VisualizeMyRoutineTaskActualStartedAtCountResDto, *exceptions.Exception)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto) (*dtos.VisualizeMyRoutineTaskActualEndedAtCountResDto, *exceptions.Exception)

	SearchPrivateRoutineTasks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineTaskInput) (*gqlmodels.SearchRoutineTaskConnection, *exceptions.Exception)
}

type RoutineTaskService struct {
	db                        *gorm.DB
	routineTaskScope          scopes.RoutineTaskScopeInterface
	routineTaskRepository     repositories.RoutineTaskRepositoryInterface
	routineTaskPayloadAdapter adapters.RoutineTaskPayloadAdapterInterface
}

func NewRoutineTaskService(
	db *gorm.DB,
	routineTaskScope scopes.RoutineTaskScopeInterface,
	routineTaskRepository repositories.RoutineTaskRepositoryInterface,
	routineTaskPayloadAdapter adapters.RoutineTaskPayloadAdapterInterface,
) RoutineTaskServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	if routineTaskScope == nil {
		routineTaskScope = scopes.NewRoutineTaskScope()
	}
	if routineTaskPayloadAdapter == nil {
		routineTaskPayloadAdapter = adapters.NewRoutineTaskPayloadAdapter(nil)
	}
	return &RoutineTaskService{
		db:                        db,
		routineTaskScope:          routineTaskScope,
		routineTaskRepository:     routineTaskRepository,
		routineTaskPayloadAdapter: routineTaskPayloadAdapter,
	}
}

/* ============================== Helper function ============================== */

func (s *RoutineTaskService) visualizeMyRoutineTaskTimeCount(
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
		BucketStart      time.Time `gorm:"column:bucket_start;"`
		RoutineTaskCount int64     `gorm:"column:routine_task_count;"`
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
			COUNT(uts.station_id) AS routine_task_count
		`).
		Joins(
			`LEFT JOIN "RoutineTaskTable" routine_task
				ON routine_task.`+columnName+` >= buckets.bucket_start
				AND routine_task.`+columnName+` < buckets.bucket_start + ?::integer * interval '1 hour'`,
			timeHourUnit,
		).
		Joins(
			`LEFT JOIN "RoutineTable" routine
				ON routine.id = routine_task.routine_id
				AND routine.deleted_at IS NULL`,
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
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(err)
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
			Value: bucket.RoutineTaskCount,
			Meta:  meta,
		}
	}

	return data, nil
}

/* ============================== Service Methods for RoutineTask ============================== */

func (s *RoutineTaskService) GetMyRoutineTaskById(
	ctx context.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto,
) (*dtos.GetMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.IsDeleted != nil && *reqDto.Param.IsDeleted {
		return nil, exceptions.RoutineTask.NotFound()
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
		RoutineId:       routineTask.RoutineId,
		Title:           routineTask.Title,
		Purpose:         routineTask.Purpose,
		Payload:         routineTask.Payload,
		CostUnit:        routineTask.CostUnit,
		Priority:        routineTask.Priority,
		Status:          routineTask.Status,
		Attempts:        routineTask.Attempts,
		MaxAttempts:     routineTask.MaxAttempts,
		Period:          routineTask.Period,
		NextScheduledAt: routineTask.NextScheduledAt,
		ScheduledAt:     routineTask.ScheduledAt,
		ActualStartedAt: routineTask.ActualStartedAt,
		ActualEndedAt:   routineTask.ActualEndedAt,
		UpdatedAt:       routineTask.UpdatedAt,
		CreatedAt:       routineTask.CreatedAt,
	}, nil
}

func (s *RoutineTaskService) GetAllMyRoutineTasksByRoutineIds(
	ctx context.Context, reqDto *dtos.GetAllMyRoutineTasksByRoutineIdsReqDto,
) (*dtos.GetAllMyRoutineTasksByRoutineIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.AreDeleted != nil && *reqDto.Param.AreDeleted {
		resDto := dtos.GetAllMyRoutineTasksByRoutineIdsResDto{}
		return &resDto, nil
	}

	db := s.db.WithContext(ctx)

	routineTasks, exception := s.routineTaskRepository.GetAllByRoutineIds(
		reqDto.Param.RoutineIds,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutineTasksByRoutineIdsResDto, len(routineTasks))
	for index, routineTask := range routineTasks {
		resDto[index] = struct {
			Id              uuid.UUID                "json:\"id\""
			RoutineId       uuid.UUID                "json:\"routineId\""
			Title           string                   "json:\"title\""
			Purpose         enums.RoutineTaskPurpose "json:\"purpose\""
			CostUnit        int64                    "json:\"costUnit\""
			Priority        int32                    "json:\"priority\""
			Status          enums.RoutineTaskStatus  "json:\"status\""
			Attempts        int32                    "json:\"attempts\""
			MaxAttempts     int32                    "json:\"maxAttempts\""
			Period          *enums.RoutinePeriod     "json:\"period\""
			NextScheduledAt time.Time                "json:\"nextScheduledAt\""
			ScheduledAt     time.Time                "json:\"scheduledAt\""
			ActualStartedAt *time.Time               "json:\"actualStartedAt\""
			ActualEndedAt   *time.Time               "json:\"actualEndedAt\""
			UpdatedAt       time.Time                "json:\"updatedAt\""
			CreatedAt       time.Time                "json:\"createdAt\""
		}{
			Id:              routineTask.Id,
			RoutineId:       routineTask.RoutineId,
			Title:           routineTask.Title,
			Purpose:         routineTask.Purpose,
			CostUnit:        routineTask.CostUnit,
			Priority:        routineTask.Priority,
			Status:          routineTask.Status,
			Attempts:        routineTask.Attempts,
			MaxAttempts:     routineTask.MaxAttempts,
			Period:          routineTask.Period,
			NextScheduledAt: routineTask.NextScheduledAt,
			ScheduledAt:     routineTask.ScheduledAt,
			ActualStartedAt: routineTask.ActualStartedAt,
			ActualEndedAt:   routineTask.ActualEndedAt,
			UpdatedAt:       routineTask.UpdatedAt,
			CreatedAt:       routineTask.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *RoutineTaskService) GetAllMyRoutineTasks(
	ctx context.Context, reqDto *dtos.GetAllMyRoutineTasksReqDto,
) (*dtos.GetAllMyRoutineTasksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.AreDeleted != nil && *reqDto.Param.AreDeleted {
		resDto := dtos.GetAllMyRoutineTasksResDto{}
		return &resDto, nil
	}

	db := s.db.WithContext(ctx)

	routineTasks, exception := s.routineTaskRepository.GetAllByUserId(
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutineTasksResDto, len(routineTasks))
	for index, routineTask := range routineTasks {
		resDto[index] = dtos.GetMyRoutineTaskByIdResDto{
			Id:              routineTask.Id,
			RoutineId:       routineTask.RoutineId,
			Title:           routineTask.Title,
			Purpose:         routineTask.Purpose,
			Payload:         routineTask.Payload,
			CostUnit:        routineTask.CostUnit,
			Priority:        routineTask.Priority,
			Status:          routineTask.Status,
			Attempts:        routineTask.Attempts,
			MaxAttempts:     routineTask.MaxAttempts,
			Period:          routineTask.Period,
			NextScheduledAt: routineTask.NextScheduledAt,
			ScheduledAt:     routineTask.ScheduledAt,
			ActualStartedAt: routineTask.ActualStartedAt,
			ActualEndedAt:   routineTask.ActualEndedAt,
			UpdatedAt:       routineTask.UpdatedAt,
			CreatedAt:       routineTask.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *RoutineTaskService) CreateRoutineTaskByRoutineId(
	ctx context.Context, reqDto *dtos.CreateRoutineTaskByRoutineIdReqDto,
) (*dtos.CreateRoutineTaskByRoutineIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if exception := s.routineTaskPayloadAdapter.Parse(reqDto.Body.Purpose, reqDto.Body.Payload); exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	newRoutineTaskId, exception := s.routineTaskRepository.CreateOneByRoutineId(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		inputs.CreateRoutineTaskInput{
			Title:           reqDto.Body.Title,
			Purpose:         reqDto.Body.Purpose,
			Payload:         reqDto.Body.Payload,
			Priority:        reqDto.Body.Priority,
			MaxAttempts:     reqDto.Body.MaxAttempts,
			Period:          reqDto.Body.Period,
			NextScheduledAt: reqDto.Body.NextScheduledAt,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateRoutineTaskByRoutineIdResDto{
		Id:        *newRoutineTaskId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) UpdateMyRoutineTaskById(
	ctx context.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto,
) (*dtos.UpdateMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if reqDto.Body.Values.Purpose != nil || reqDto.Body.Values.Payload != nil {
		finalPurpose := reqDto.Body.Values.Purpose
		finalPayload := reqDto.Body.Values.Payload
		if finalPurpose == nil || finalPayload == nil {
			existingRoutineTask, exception := s.routineTaskRepository.GetOneById(
				reqDto.Body.RoutineTaskId,
				reqDto.ContextFields.UserId,
				nil,
				options.WithDB(db),
			)
			if exception != nil {
				return nil, exception
			}
			if finalPurpose == nil {
				finalPurpose = &existingRoutineTask.Purpose
			}
			if finalPayload == nil {
				finalPayload = &existingRoutineTask.Payload
			}
		}
		if exception := s.routineTaskPayloadAdapter.Parse(*finalPurpose, *finalPayload); exception != nil {
			return nil, exception
		}
	}

	updatedRoutineTask, exception := s.routineTaskRepository.UpdateOneById(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateRoutineTaskInput{
			Values: inputs.UpdateRoutineTaskInput{
				RoutineId:       reqDto.Body.Values.RoutineId,
				Title:           reqDto.Body.Values.Title,
				Purpose:         reqDto.Body.Values.Purpose,
				Payload:         reqDto.Body.Values.Payload,
				Priority:        reqDto.Body.Values.Priority,
				MaxAttempts:     reqDto.Body.Values.MaxAttempts,
				Period:          reqDto.Body.Values.Period,
				NextScheduledAt: reqDto.Body.Values.NextScheduledAt,
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

func (s *RoutineTaskService) PauseMyRoutineTaskById(
	ctx context.Context, reqDto *dtos.PauseMyRoutineTaskByIdReqDto,
) (*dtos.PauseMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}
	routineTask, exception := s.routineTaskRepository.CheckPermissionAndGetOneById(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if routineTask.Status != enums.RoutineTaskStatus_Idle {
		tx.Rollback()
		return nil, exceptions.RoutineTask.InvalidInput("only idle routine tasks can be paused")
	}

	now := time.Now()
	result := tx.Model(&schemas.RoutineTask{}).
		Where("id = ? AND status = ?", reqDto.Body.RoutineTaskId, enums.RoutineTaskStatus_Idle).
		Updates(map[string]any{
			"status":     enums.RoutineTaskStatus_Pause,
			"updated_at": now,
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, exceptions.RoutineTask.NoChanges()
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.PauseMyRoutineTaskByIdResDto{UpdatedAt: now}, nil
}

func (s *RoutineTaskService) ResumeMyRoutineTaskById(
	ctx context.Context, reqDto *dtos.ResumeMyRoutineTaskByIdReqDto,
) (*dtos.ResumeMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}
	routineTask, exception := s.routineTaskRepository.CheckPermissionAndGetOneById(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if routineTask.Status != enums.RoutineTaskStatus_Pause {
		tx.Rollback()
		return nil, exceptions.RoutineTask.InvalidInput("only paused routine tasks can be resumed")
	}

	now := time.Now()
	result := tx.Model(&schemas.RoutineTask{}).
		Where("id = ? AND status = ?", reqDto.Body.RoutineTaskId, enums.RoutineTaskStatus_Pause).
		Updates(map[string]any{
			"status":     enums.RoutineTaskStatus_Idle,
			"updated_at": now,
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, exceptions.RoutineTask.NoChanges()
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.ResumeMyRoutineTaskByIdResDto{UpdatedAt: now}, nil
}

func (s *RoutineTaskService) HardDeleteMyRoutineTaskById(
	ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto,
) (*dtos.HardDeleteMyRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
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
	ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto,
) (*dtos.HardDeleteMyRoutineTasksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
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

/* ============================== Service Methods for Charts ============================== */

func (s *RoutineTaskService) VisualizeMyRoutineTaskStatusCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskStatusCountReqDto,
) (*dtos.VisualizeMyRoutineTaskStatusCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	var rows []struct {
		Status           enums.RoutineTaskStatus `gorm:"column:status;"`
		RoutineTaskCount int64                   `gorm:"column:routine_task_count;"`
	}
	result := db.Model(&schemas.RoutineTask{}).
		Select(`"RoutineTaskTable".status AS status, COUNT(*) AS routine_task_count`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Group(`"RoutineTaskTable".status`).
		Scan(&rows)
	if err := result.Error; err != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	counts := make(map[enums.RoutineTaskStatus]int64, len(rows))
	for _, row := range rows {
		counts[row.Status] = row.RoutineTaskCount
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(enums.AllRoutineTaskStatuses))
	for index, status := range enums.AllRoutineTaskStatuses {
		metadata := map[string]string{"status": status.String()}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    status.String() + "-routine-task-count",
			X:     status.String() + " Routine Task Count",
			Value: counts[status],
			Meta:  meta,
		}
	}

	return &dtos.VisualizeMyRoutineTaskStatusCountResDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskPurposeCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskPurposeCountReqDto,
) (*dtos.VisualizeMyRoutineTaskPurposeCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	var rows []struct {
		Purpose          enums.RoutineTaskPurpose `gorm:"column:purpose;"`
		RoutineTaskCount int64                    `gorm:"column:routine_task_count;"`
	}
	result := db.Model(&schemas.RoutineTask{}).
		Select(`"RoutineTaskTable".purpose AS purpose, COUNT(*) AS routine_task_count`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id AND routine.deleted_at IS NULL`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Group(`"RoutineTaskTable".purpose`).
		Scan(&rows)
	if err := result.Error; err != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	counts := make(map[enums.RoutineTaskPurpose]int64, len(rows))
	for _, row := range rows {
		counts[row.Purpose] = row.RoutineTaskCount
	}

	data := make([]dtos.TwoDimensionalDatum[int64], len(enums.AllRoutineTaskPurposes))
	for index, purpose := range enums.AllRoutineTaskPurposes {
		metadata := map[string]string{"purpose": purpose.String()}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, exceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = dtos.TwoDimensionalDatum[int64]{
			Id:    purpose.String() + "-routine-task-count",
			X:     purpose.String() + " Routine Task Count",
			Value: counts[purpose],
			Meta:  meta,
		}
	}

	return &dtos.VisualizeMyRoutineTaskPurposeCountResDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskScheduledAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskScheduledAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"scheduled_at",
		"scheduledAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskScheduledAtCountResDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskActualStartedAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskActualStartedAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_started_at",
		"actualStartedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskActualStartedAtCountResDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskActualEndedAtCount(
	ctx context.Context, reqDto *dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto,
) (*dtos.VisualizeMyRoutineTaskActualEndedAtCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !util.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, exceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		reqDto.ContextFields.UserId,
		reqDto.Param.Permission,
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_ended_at",
		"actualEndedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.VisualizeMyRoutineTaskActualEndedAtCountResDto{
		Data: data,
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

	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	query := db.Model(&schemas.RoutineTask{}).
		Select(`"RoutineTaskTable".*, uts.permission AS permission`).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id AND routine.deleted_at IS NULL`).
		Joins(`LEFT JOIN "UsersToStationsTable" uts ON routine.station_id = uts.station_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)

	if len(gqlInput.RoutineIds) > 0 {
		query = query.Where(
			`"RoutineTaskTable".routine_id IN ?`,
			gqlInput.RoutineIds,
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
			query = query.Order("title " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByPurpose:
			query = query.Order("purpose " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByPriority:
			query = query.Order("priority " + cending).
				Order("scheduled_at " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByStatus:
			query = query.Order("status " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByAttempts:
			query = query.Order("attempts " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByMaxAttempts:
			query = query.Order("max_attempts " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByScheduledAt:
			query = query.Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByActualStartedAt:
			query = query.Order("actual_started_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByActualEndedAt:
			query = query.Order("actual_ended_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByLastUpdate:
			query = query.Order("updated_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineTaskSortByCreatedAt:
			query = query.Order("created_at " + cending).
				Order("scheduled_at " + cending).
				Order("priority " + cending).
				Order("updated_at " + cending)
		default:
			query = query.Order("scheduled_at " + cending).
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
	if err := query.Scopes(s.routineTaskScope.IncludePreloads(
		nil,
	)).Find(&routineTasks).Error; err != nil {
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
