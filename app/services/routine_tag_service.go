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

type RoutineTagServiceInterface interface {
	GetOneById(ctx context.Context, reqDto *dtos.GetOneRoutineTagByIdReqDto) (*dtos.GetOneRoutineTagByIdResDto, *exceptions.Exception)
	CreateOneByUserId(ctx context.Context, reqDto *dtos.CreateOneRoutineTagByUserIdReqDto) (*dtos.CreateOneRoutineTagByUserIdResDto, *exceptions.Exception)
	BulkCreateManyByUserId(ctx context.Context, reqDto *dtos.BulkCreateManyRoutineTagsByUserIdReqDto) (*dtos.BulkCreateManyRoutineTagsByUserIdResDto, *exceptions.Exception)
	UpdateOneById(ctx context.Context, reqDto *dtos.UpdateOneRoutineTagByIdReqDto) (*dtos.UpdateOneRoutineTagByIdResDto, *exceptions.Exception)
	BulkUpdateManyByIds(ctx context.Context, reqDto *dtos.BulkUpdateManyRoutineTagsByIdsReqDto) (*dtos.BulkUpdateManyRoutineTagsByIdsResDto, *exceptions.Exception)
	HardDeleteOneById(ctx context.Context, reqDto *dtos.HardDeleteOneRoutineTagByIdReqDto) (*dtos.HardDeleteOneRoutineTagByIdResDto, *exceptions.Exception)
	HardDeleteManyByIds(ctx context.Context, reqDto *dtos.HardDeleteManyRoutineTagsByIdsReqDto) (*dtos.HardDeleteManyRoutineTagsByIdsResDto, *exceptions.Exception)
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

func (s *RoutineTagService) GetOneById(
	ctx context.Context,
	reqDto *dtos.GetOneRoutineTagByIdReqDto,
) (*dtos.GetOneRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.GetOneRoutineTagByIdResDto{
		Id:        routineTag.Id,
		Name:      routineTag.Name,
		Color:     routineTag.Color,
		Icon:      routineTag.Icon,
		UpdatedAt: routineTag.UpdatedAt,
		CreatedAt: routineTag.CreatedAt,
	}, nil
}

func (s *RoutineTagService) CreateOneByUserId(
	ctx context.Context,
	reqDto *dtos.CreateOneRoutineTagByUserIdReqDto,
) (*dtos.CreateOneRoutineTagByUserIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.CreateOneRoutineTagByUserIdResDto{
		Id:        *newRoutineTagId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) BulkCreateManyByUserId(
	ctx context.Context,
	reqDto *dtos.BulkCreateManyRoutineTagsByUserIdReqDto,
) (*dtos.BulkCreateManyRoutineTagsByUserIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.BulkCreateManyRoutineTagsByUserIdResDto{
		Ids:       newRoutineTagIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) UpdateOneById(
	ctx context.Context,
	reqDto *dtos.UpdateOneRoutineTagByIdReqDto,
) (*dtos.UpdateOneRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.UpdateOneRoutineTagByIdResDto{
		UpdatedAt: updatedRoutineTag.UpdatedAt,
	}, nil
}

func (s *RoutineTagService) BulkUpdateManyByIds(
	ctx context.Context,
	reqDto *dtos.BulkUpdateManyRoutineTagsByIdsReqDto,
) (*dtos.BulkUpdateManyRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.BulkUpdateManyRoutineTagsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteOneById(
	ctx context.Context,
	reqDto *dtos.HardDeleteOneRoutineTagByIdReqDto,
) (*dtos.HardDeleteOneRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.HardDeleteOneRoutineTagByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteManyRoutineTagsByIdsReqDto,
) (*dtos.HardDeleteManyRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
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

	return &dtos.HardDeleteManyRoutineTagsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
