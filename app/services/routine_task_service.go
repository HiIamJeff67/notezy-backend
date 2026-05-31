package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
)

type RoutineTaskServiceInterface interface {
	GetOneById(ctx context.Context, reqDto *dtos.GetOneRoutineTaskByIdReqDto) (*dtos.GetOneRoutineTaskByIdResDto, *exceptions.Exception)
	CreateOneByStationId(ctx context.Context, reqDto *dtos.CreateOneRoutineTaskByStationIdReqDto) (*dtos.CreateOneRoutineTaskByStationIdResDto, *exceptions.Exception)
	UpdateOneById(ctx context.Context, reqDto *dtos.UpdateOneRoutineTaskByIdReqDto) (*dtos.UpdateOneRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteOneById(ctx context.Context, reqDto *dtos.HardDeleteOneRoutineTaskByIdReqDto) (*dtos.HardDeleteOneRoutineTaskByIdResDto, *exceptions.Exception)
	HardDeleteManyByIds(ctx context.Context, reqDto *dtos.HardDeleteManyRoutineTasksByIdsReqDto) (*dtos.HardDeleteManyRoutineTasksByIdsResDto, *exceptions.Exception)
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

func (s *RoutineTaskService) GetOneById(
	ctx context.Context,
	reqDto *dtos.GetOneRoutineTaskByIdReqDto,
) (*dtos.GetOneRoutineTaskByIdResDto, *exceptions.Exception) {
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

	return &dtos.GetOneRoutineTaskByIdResDto{
		Id:              routineTask.Id,
		StationId:       routineTask.StationId,
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

func (s *RoutineTaskService) CreateOneByStationId(
	ctx context.Context,
	reqDto *dtos.CreateOneRoutineTaskByStationIdReqDto,
) (*dtos.CreateOneRoutineTaskByStationIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	newRoutineTaskId, exception := s.routineTaskRepository.CreateOneByStationId(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		inputs.CreateRoutineTaskInput{
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

	return &dtos.CreateOneRoutineTaskByStationIdResDto{
		Id:        *newRoutineTaskId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) UpdateOneById(
	ctx context.Context,
	reqDto *dtos.UpdateOneRoutineTaskByIdReqDto,
) (*dtos.UpdateOneRoutineTaskByIdResDto, *exceptions.Exception) {
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

	return &dtos.UpdateOneRoutineTaskByIdResDto{
		UpdatedAt: updatedRoutineTask.UpdatedAt,
	}, nil
}

func (s *RoutineTaskService) HardDeleteOneById(
	ctx context.Context,
	reqDto *dtos.HardDeleteOneRoutineTaskByIdReqDto,
) (*dtos.HardDeleteOneRoutineTaskByIdResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteOneRoutineTaskByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTaskService) HardDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteManyRoutineTasksByIdsReqDto,
) (*dtos.HardDeleteManyRoutineTasksByIdsResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteManyRoutineTasksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
