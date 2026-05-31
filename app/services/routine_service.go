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
	types "notezy-backend/shared/types"
)

type RoutineServiceInterface interface {
	GetOneById(ctx context.Context, reqDto *dtos.GetOneRoutineByIdReqDto) (*dtos.GetOneRoutineByIdResDto, *exceptions.Exception)
	CreateOneByStationId(ctx context.Context, reqDto *dtos.CreateOneRoutineByStationIdReqDto) (*dtos.CreateOneRoutineByStationIdResDto, *exceptions.Exception)
	BulkCreateManyByStationIds(ctx context.Context, reqDto *dtos.BulkCreateManyRoutinesByStationIdsReqDto) (*dtos.BulkCreateManyRoutinesByStationIdsResDto, *exceptions.Exception)
	UpdateOneById(ctx context.Context, reqDto *dtos.UpdateOneRoutineByIdReqDto) (*dtos.UpdateOneRoutineByIdResDto, *exceptions.Exception)
	BulkUpdateManyByIds(ctx context.Context, reqDto *dtos.BulkUpdateManyRoutinesByIdsReqDto) (*dtos.BulkUpdateManyRoutinesByIdsResDto, *exceptions.Exception)
	RestoreSoftDeletedOneById(ctx context.Context, reqDto *dtos.RestoreSoftDeletedOneRoutineByIdReqDto) (*dtos.RestoreSoftDeletedOneRoutineByIdResDto, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ctx context.Context, reqDto *dtos.RestoreSoftDeletedManyRoutinesByIdsReqDto) (*dtos.RestoreSoftDeletedManyRoutinesByIdsResDto, *exceptions.Exception)
	SoftDeleteOneById(ctx context.Context, reqDto *dtos.SoftDeleteOneRoutineByIdReqDto) (*dtos.SoftDeleteOneRoutineByIdResDto, *exceptions.Exception)
	SoftDeleteManyByIds(ctx context.Context, reqDto *dtos.SoftDeleteManyRoutinesByIdsReqDto) (*dtos.SoftDeleteManyRoutinesByIdsResDto, *exceptions.Exception)
	HardDeleteOneById(ctx context.Context, reqDto *dtos.HardDeleteOneRoutineByIdReqDto) (*dtos.HardDeleteOneRoutineByIdResDto, *exceptions.Exception)
	HardDeleteManyByIds(ctx context.Context, reqDto *dtos.HardDeleteManyRoutinesByIdsReqDto) (*dtos.HardDeleteManyRoutinesByIdsResDto, *exceptions.Exception)
}

type RoutineService struct {
	db                *gorm.DB
	routineRepository repositories.RoutineRepositoryInterface
}

func NewRoutineService(
	db *gorm.DB,
	routineRepository repositories.RoutineRepositoryInterface,
) RoutineServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RoutineService{
		db:                db,
		routineRepository: routineRepository,
	}
}

/* ============================== Service Methods for Routine ============================== */

func (s *RoutineService) GetOneById(
	ctx context.Context,
	reqDto *dtos.GetOneRoutineByIdReqDto,
) (*dtos.GetOneRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	onlyDeleted := types.Ternary_Negative
	if reqDto.Param.OnlyDeleted != nil {
		onlyDeleted = *reqDto.Param.OnlyDeleted
	}

	db := s.db.WithContext(ctx)
	routine, exception := s.routineRepository.GetOneById(
		reqDto.Param.RoutineId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetOneRoutineByIdResDto{
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
	}, nil
}

func (s *RoutineService) CreateOneByStationId(
	ctx context.Context,
	reqDto *dtos.CreateOneRoutineByStationIdReqDto,
) (*dtos.CreateOneRoutineByStationIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.CreateOneRoutineByStationIdResDto{
		Id:        *newRoutineId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkCreateManyByStationIds(
	ctx context.Context,
	reqDto *dtos.BulkCreateManyRoutinesByStationIdsReqDto,
) (*dtos.BulkCreateManyRoutinesByStationIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.BulkCreateManyRoutinesByStationIdsResDto{
		Ids:       newRoutineIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) UpdateOneById(
	ctx context.Context,
	reqDto *dtos.UpdateOneRoutineByIdReqDto,
) (*dtos.UpdateOneRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.UpdateOneRoutineByIdResDto{
		UpdatedAt: updatedRoutine.UpdatedAt,
	}, nil
}

func (s *RoutineService) BulkUpdateManyByIds(
	ctx context.Context,
	reqDto *dtos.BulkUpdateManyRoutinesByIdsReqDto,
) (*dtos.BulkUpdateManyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.BulkUpdateManyRoutinesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) RestoreSoftDeletedOneById(
	ctx context.Context,
	reqDto *dtos.RestoreSoftDeletedOneRoutineByIdReqDto,
) (*dtos.RestoreSoftDeletedOneRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.RestoreSoftDeletedOneRoutineByIdResDto{
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

func (s *RoutineService) RestoreSoftDeletedManyByIds(
	ctx context.Context,
	reqDto *dtos.RestoreSoftDeletedManyRoutinesByIdsReqDto,
) (*dtos.RestoreSoftDeletedManyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	resDto := dtos.RestoreSoftDeletedManyRoutinesByIdsResDto{}
	for _, restoredRoutine := range restoredRoutines {
		resDto = append(resDto, dtos.RestoreSoftDeletedOneRoutineByIdResDto{
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

func (s *RoutineService) SoftDeleteOneById(
	ctx context.Context,
	reqDto *dtos.SoftDeleteOneRoutineByIdReqDto,
) (*dtos.SoftDeleteOneRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.SoftDeleteOneRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) SoftDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.SoftDeleteManyRoutinesByIdsReqDto,
) (*dtos.SoftDeleteManyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.SoftDeleteManyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteOneById(
	ctx context.Context,
	reqDto *dtos.HardDeleteOneRoutineByIdReqDto,
) (*dtos.HardDeleteOneRoutineByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.HardDeleteOneRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteManyRoutinesByIdsReqDto,
) (*dtos.HardDeleteManyRoutinesByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.HardDeleteManyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
