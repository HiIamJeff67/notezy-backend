package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineServiceInterface interface {
	GetMyRoutineById(ctx context.Context, reqDto *dtos.GetMyRoutineByIdReqDto) (*dtos.GetMyRoutineByIdResDto, *exceptions.Exception)
	CreateRoutineByStationId(ctx context.Context, reqDto *dtos.CreateRoutineByStationIdReqDto) (*dtos.CreateRoutineByStationIdResDto, *exceptions.Exception)
	CreateRoutinesByStationIds(ctx context.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto) (*dtos.CreateRoutinesByStationIdsResDto, *exceptions.Exception)
	UpdateMyRoutineById(ctx context.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto) (*dtos.UpdateMyRoutineByIdResDto, *exceptions.Exception)
	UpdateMyRoutinesByIds(ctx context.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto) (*dtos.UpdateMyRoutinesByIdsResDto, *exceptions.Exception)
	RestoreMyRoutineById(ctx context.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto) (*dtos.RestoreMyRoutineByIdResDto, *exceptions.Exception)
	RestoreMyRoutinesByIds(ctx context.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto) (*dtos.RestoreMyRoutinesByIdsResDto, *exceptions.Exception)
	DeleteMyRoutineById(ctx context.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto) (*dtos.DeleteMyRoutineByIdResDto, *exceptions.Exception)
	DeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto) (*dtos.DeleteMyRoutinesByIdsResDto, *exceptions.Exception)
	HardDeleteMyRoutineById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto) (*dtos.HardDeleteMyRoutineByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto) (*dtos.HardDeleteMyRoutinesByIdsResDto, *exceptions.Exception)
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

func (s *RoutineService) GetMyRoutineById(
	ctx context.Context,
	reqDto *dtos.GetMyRoutineByIdReqDto,
) (*dtos.GetMyRoutineByIdResDto, *exceptions.Exception) {
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
	}, nil
}

func (s *RoutineService) CreateRoutineByStationId(
	ctx context.Context,
	reqDto *dtos.CreateRoutineByStationIdReqDto,
) (*dtos.CreateRoutineByStationIdResDto, *exceptions.Exception) {
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

	return &dtos.CreateRoutineByStationIdResDto{
		Id:        *newRoutineId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) CreateRoutinesByStationIds(
	ctx context.Context,
	reqDto *dtos.CreateRoutinesByStationIdsReqDto,
) (*dtos.CreateRoutinesByStationIdsResDto, *exceptions.Exception) {
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

	return &dtos.CreateRoutinesByStationIdsResDto{
		Ids:       newRoutineIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) UpdateMyRoutineById(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutineByIdReqDto,
) (*dtos.UpdateMyRoutineByIdResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyRoutineByIdResDto{
		UpdatedAt: updatedRoutine.UpdatedAt,
	}, nil
}

func (s *RoutineService) UpdateMyRoutinesByIds(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutinesByIdsReqDto,
) (*dtos.UpdateMyRoutinesByIdsResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyRoutinesByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) RestoreMyRoutineById(
	ctx context.Context,
	reqDto *dtos.RestoreMyRoutineByIdReqDto,
) (*dtos.RestoreMyRoutineByIdResDto, *exceptions.Exception) {
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
	ctx context.Context,
	reqDto *dtos.RestoreMyRoutinesByIdsReqDto,
) (*dtos.RestoreMyRoutinesByIdsResDto, *exceptions.Exception) {
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
	ctx context.Context,
	reqDto *dtos.DeleteMyRoutineByIdReqDto,
) (*dtos.DeleteMyRoutineByIdResDto, *exceptions.Exception) {
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

	return &dtos.DeleteMyRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) DeleteMyRoutinesByIds(
	ctx context.Context,
	reqDto *dtos.DeleteMyRoutinesByIdsReqDto,
) (*dtos.DeleteMyRoutinesByIdsResDto, *exceptions.Exception) {
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

	return &dtos.DeleteMyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteMyRoutineById(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineByIdReqDto,
) (*dtos.HardDeleteMyRoutineByIdResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyRoutineByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineService) HardDeleteMyRoutinesByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto,
) (*dtos.HardDeleteMyRoutinesByIdsResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyRoutinesByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
