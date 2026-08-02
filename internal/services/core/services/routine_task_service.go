package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-tasks"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"
	adapters "github.com/HiIamJeff67/notezy-backend/internal/adapters"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/searchcursor"
	timeutil "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/timeutil"
)

type RoutineTaskServiceInterface interface {
	GetMyRoutineTaskById(ctx context.Context, reqDto *routinetasksdto.GetMyRoutineTaskByIdRequestDto) (*routinetasksdto.GetMyRoutineTaskByIdResponseDto, *exceptions.Exception)
	GetAllMyRoutineTasksByRoutineIds(ctx context.Context, reqDto *routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto) (*routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto, *exceptions.Exception)
	GetAllMyRoutineTasks(ctx context.Context, reqDto *routinetasksdto.GetAllMyRoutineTasksRequestDto) (*routinetasksdto.GetAllMyRoutineTasksResponseDto, *exceptions.Exception)
	CreateRoutineTaskByRoutineId(ctx context.Context, reqDto *routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto) (*routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto, *exceptions.Exception)
	UpdateMyRoutineTaskById(ctx context.Context, reqDto *routinetasksdto.UpdateMyRoutineTaskByIdRequestDto) (*routinetasksdto.UpdateMyRoutineTaskByIdResponseDto, *exceptions.Exception)
	PauseMyRoutineTaskById(ctx context.Context, reqDto *routinetasksdto.PauseMyRoutineTaskByIdRequestDto) (*routinetasksdto.PauseMyRoutineTaskByIdResponseDto, *exceptions.Exception)
	ResumeMyRoutineTaskById(ctx context.Context, reqDto *routinetasksdto.ResumeMyRoutineTaskByIdRequestDto) (*routinetasksdto.ResumeMyRoutineTaskByIdResponseDto, *exceptions.Exception)
	HardDeleteMyRoutineTaskById(ctx context.Context, reqDto *routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto) (*routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto, *exceptions.Exception)
	HardDeleteMyRoutineTasksByIds(ctx context.Context, reqDto *routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto) (*routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto, *exceptions.Exception)
	VisualizeMyRoutineTaskStatusCount(ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto) (*routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto, *exceptions.Exception)
	VisualizeMyRoutineTaskPurposeCount(ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto) (*routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto, *exceptions.Exception)
	VisualizeMyRoutineTaskScheduledAtCount(ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto) (*routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto, *exceptions.Exception)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto) (*routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto, *exceptions.Exception)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto) (*routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto, *exceptions.Exception)

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
		db = data.NotezyDB
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
) ([]routinetasksdto.RoutineTaskCountDatum, *exceptions.Exception) {
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
		return nil, apiexceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	data := make([]routinetasksdto.RoutineTaskCountDatum, len(buckets))
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
			return nil, apiexceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = routinetasksdto.RoutineTaskCountDatum{
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
	ctx context.Context, reqDto *routinetasksdto.GetMyRoutineTaskByIdRequestDto,
) (*routinetasksdto.GetMyRoutineTaskByIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.IsDeleted != nil && *reqDto.Param.IsDeleted {
		return nil, apiexceptions.RoutineTask.NotFound()
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	routineTask, exception := s.routineTaskRepository.GetOneById(
		reqDto.Param.RoutineTaskId,
		actorUserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.GetMyRoutineTaskByIdResponseDto{
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
	ctx context.Context, reqDto *routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto,
) (*routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.AreDeleted != nil && *reqDto.Param.AreDeleted {
		resDto := routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto{}
		return &resDto, nil
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	routineTasks, exception := s.routineTaskRepository.GetAllByRoutineIds(
		reqDto.Param.RoutineIds,
		actorUserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto, len(routineTasks))
	for index, routineTask := range routineTasks {
		resDto[index] = routinetasksdto.RoutineTaskResponseDto{
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
	ctx context.Context, reqDto *routinetasksdto.GetAllMyRoutineTasksRequestDto,
) (*routinetasksdto.GetAllMyRoutineTasksResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if reqDto.Param.AreDeleted != nil && *reqDto.Param.AreDeleted {
		resDto := routinetasksdto.GetAllMyRoutineTasksResponseDto{}
		return &resDto, nil
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	routineTasks, exception := s.routineTaskRepository.GetAllByUserId(
		actorUserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(routinetasksdto.GetAllMyRoutineTasksResponseDto, len(routineTasks))
	for index, routineTask := range routineTasks {
		resDto[index] = routinetasksdto.GetMyRoutineTaskByIdResponseDto{
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
	ctx context.Context, reqDto *routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto,
) (*routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if exception := s.routineTaskPayloadAdapter.Parse(reqDto.Body.Purpose, reqDto.Body.Payload); exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	newRoutineTaskId, exception := s.routineTaskRepository.CreateOneByRoutineId(
		reqDto.Body.RoutineId,
		actorUserId,
		inputs.CreateRoutineTaskInput{
			ActorUserId:     actorUserId,
			Title:           reqDto.Body.Title,
			Purpose:         reqDto.Body.Purpose,
			Payload:         reqDto.Body.Payload,
			Priority:        reqDto.Body.Priority,
			MaxAttempts:     reqDto.Body.MaxAttempts,
			Period:          reqDto.Body.Period,
			NextScheduledAt: reqDto.Body.NextScheduledAt,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto{
		Id:        *newRoutineTaskId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) UpdateMyRoutineTaskById(
	ctx context.Context, reqDto *routinetasksdto.UpdateMyRoutineTaskByIdRequestDto,
) (*routinetasksdto.UpdateMyRoutineTaskByIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	if reqDto.Body.Values.Purpose != nil || reqDto.Body.Values.Payload != nil {
		finalPurpose := reqDto.Body.Values.Purpose
		finalPayload := reqDto.Body.Values.Payload
		if finalPurpose == nil || finalPayload == nil {
			existingRoutineTask, exception := s.routineTaskRepository.GetOneById(
				reqDto.Body.RoutineTaskId,
				actorUserId,
				nil,
				options.WithDB(db),
				options.WithAllowedPermissions(allowedPermissions),
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
		actorUserId,
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
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.UpdateMyRoutineTaskByIdResponseDto{
		UpdatedAt: updatedRoutineTask.UpdatedAt,
	}, nil
}

func (s *RoutineTaskService) PauseMyRoutineTaskById(
	ctx context.Context, reqDto *routinetasksdto.PauseMyRoutineTaskByIdRequestDto,
) (*routinetasksdto.PauseMyRoutineTaskByIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()
	routineTask, exception := s.routineTaskRepository.CheckPermissionAndGetOneById(
		reqDto.Body.RoutineTaskId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if routineTask.Status != enums.RoutineTaskStatus_Idle {
		tx.Rollback()
		return nil, apiexceptions.RoutineTask.InvalidInput("only idle routine tasks can be paused")
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
		return nil, apiexceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, apiexceptions.RoutineTask.NoChanges()
	}

	if err := tx.Commit().Error; err != nil {
		return nil, apiexceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return &routinetasksdto.PauseMyRoutineTaskByIdResponseDto{UpdatedAt: now}, nil
}

func (s *RoutineTaskService) ResumeMyRoutineTaskById(
	ctx context.Context, reqDto *routinetasksdto.ResumeMyRoutineTaskByIdRequestDto,
) (*routinetasksdto.ResumeMyRoutineTaskByIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()
	routineTask, exception := s.routineTaskRepository.CheckPermissionAndGetOneById(
		reqDto.Body.RoutineTaskId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if routineTask.Status != enums.RoutineTaskStatus_Pause {
		tx.Rollback()
		return nil, apiexceptions.RoutineTask.InvalidInput("only paused routine tasks can be resumed")
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
		return nil, apiexceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, apiexceptions.RoutineTask.NoChanges()
	}

	if err := tx.Commit().Error; err != nil {
		return nil, apiexceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return &routinetasksdto.ResumeMyRoutineTaskByIdResponseDto{UpdatedAt: now}, nil
}

func (s *RoutineTaskService) HardDeleteMyRoutineTaskById(
	ctx context.Context, reqDto *routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto,
) (*routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.routineTaskRepository.HardDeleteOneById(
		reqDto.Body.RoutineTaskId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) HardDeleteMyRoutineTasksByIds(
	ctx context.Context, reqDto *routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto,
) (*routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.routineTaskRepository.HardDeleteManyByIds(
		reqDto.Body.RoutineTaskIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for Charts ============================== */

func (s *RoutineTaskService) VisualizeMyRoutineTaskStatusCount(
	ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto,
) (*routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
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
		Where("uts.user_id = ? AND uts.permission = ?", actorUserId, reqDto.Param.Permission).
		Group(`"RoutineTaskTable".status`).
		Scan(&rows)
	if err := result.Error; err != nil {
		return nil, apiexceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	counts := make(map[enums.RoutineTaskStatus]int64, len(rows))
	for _, row := range rows {
		counts[row.Status] = row.RoutineTaskCount
	}

	data := make([]routinetasksdto.RoutineTaskCountDatum, len(enums.AllRoutineTaskStatuses))
	for index, status := range enums.AllRoutineTaskStatuses {
		metadata := map[string]string{"status": status.String()}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, apiexceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = routinetasksdto.RoutineTaskCountDatum{
			Id:    status.String() + "-routine-task-count",
			X:     status.String() + " Routine Task Count",
			Value: counts[status],
			Meta:  meta,
		}
	}

	return &routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskPurposeCount(
	ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto,
) (*routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
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
		Where("uts.user_id = ? AND uts.permission = ?", actorUserId, reqDto.Param.Permission).
		Group(`"RoutineTaskTable".purpose`).
		Scan(&rows)
	if err := result.Error; err != nil {
		return nil, apiexceptions.RoutineTask.NotFound().WithOrigin(err)
	}

	counts := make(map[enums.RoutineTaskPurpose]int64, len(rows))
	for _, row := range rows {
		counts[row.Purpose] = row.RoutineTaskCount
	}

	data := make([]routinetasksdto.RoutineTaskCountDatum, len(enums.AllRoutineTaskPurposes))
	for index, purpose := range enums.AllRoutineTaskPurposes {
		metadata := map[string]string{"purpose": purpose.String()}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, apiexceptions.Routine.FailedToMarshalData(metadata)
		}

		data[index] = routinetasksdto.RoutineTaskCountDatum{
			Id:    purpose.String() + "-routine-task-count",
			X:     purpose.String() + " Routine Task Count",
			Value: counts[purpose],
			Meta:  meta,
		}
	}

	return &routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskScheduledAtCount(
	ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto,
) (*routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !timeutil.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		actorUserId,
		enums.AccessControlPermission(reqDto.Param.Permission),
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"scheduled_at",
		"scheduledAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskActualStartedAtCount(
	ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto,
) (*routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !timeutil.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		actorUserId,
		enums.AccessControlPermission(reqDto.Param.Permission),
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_started_at",
		"actualStartedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto{
		Data: data,
	}, nil
}

func (s *RoutineTaskService) VisualizeMyRoutineTaskActualEndedAtCount(
	ctx context.Context, reqDto *routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto,
) (*routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.QueryRangeStartedAt.Before(reqDto.Param.QueryRangeEndedAt) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt should be earlier then queryRangeEndedAt")
	}
	if !timeutil.IsTimeWithin(reqDto.Param.QueryRangeStartedAt, reqDto.Param.QueryRangeEndedAt, 360*24*time.Hour) {
		return nil, apiexceptions.RoutineTask.InvalidDto("queryRangeStartedAt and queryRangeEndedAt should be within 360 days")
	}

	data, exception := s.visualizeMyRoutineTaskTimeCount(
		ctx,
		actorUserId,
		enums.AccessControlPermission(reqDto.Param.Permission),
		reqDto.Param.TimeHourUnit,
		reqDto.Param.QueryRangeStartedAt,
		reqDto.Param.QueryRangeEndedAt,
		"actual_ended_at",
		"actualEndedAt",
	)
	if exception != nil {
		return nil, exception
	}

	return &routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto{
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

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
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
			return nil, apiexceptions.Search.FailedToDecode().WithOrigin(err)
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
		return nil, apiexceptions.RoutineTask.NotFound().WithOrigin(err)
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
			return nil, apiexceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, apiexceptions.Search.FailedToUnmarshalSearchCursor()
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
