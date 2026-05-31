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
	GetMyRoutineTagById(ctx context.Context, reqDto *dtos.GetMyRoutineTagByIdReqDto) (*dtos.GetMyRoutineTagByIdResDto, *exceptions.Exception)
	CreateRoutineTag(ctx context.Context, reqDto *dtos.CreateRoutineTagReqDto) (*dtos.CreateRoutineTagResDto, *exceptions.Exception)
	CreateRoutineTags(ctx context.Context, reqDto *dtos.CreateRoutineTagsReqDto) (*dtos.CreateRoutineTagsResDto, *exceptions.Exception)
	UpdateMyRoutineTagById(ctx context.Context, reqDto *dtos.UpdateMyRoutineTagByIdReqDto) (*dtos.UpdateMyRoutineTagByIdResDto, *exceptions.Exception)
	UpdateMyRoutineTagsByIds(ctx context.Context, reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto) (*dtos.UpdateMyRoutineTagsByIdsResDto, *exceptions.Exception)
	HardDeleteMyRoutineTagById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto) (*dtos.HardDeleteMyRoutineTagByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutineTagsByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto) (*dtos.HardDeleteMyRoutineTagsByIdsResDto, *exceptions.Exception)
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

	return &dtos.GetMyRoutineTagByIdResDto{
		Id:        routineTag.Id,
		Name:      routineTag.Name,
		Color:     routineTag.Color,
		Icon:      routineTag.Icon,
		UpdatedAt: routineTag.UpdatedAt,
		CreatedAt: routineTag.CreatedAt,
	}, nil
}

func (s *RoutineTagService) CreateRoutineTag(
	ctx context.Context,
	reqDto *dtos.CreateRoutineTagReqDto,
) (*dtos.CreateRoutineTagResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyRoutineTagByIdResDto{
		UpdatedAt: updatedRoutineTag.UpdatedAt,
	}, nil
}

func (s *RoutineTagService) UpdateMyRoutineTagsByIds(
	ctx context.Context,
	reqDto *dtos.UpdateMyRoutineTagsByIdsReqDto,
) (*dtos.UpdateMyRoutineTagsByIdsResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyRoutineTagsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteMyRoutineTagById(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTagByIdReqDto,
) (*dtos.HardDeleteMyRoutineTagByIdResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyRoutineTagByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *RoutineTagService) HardDeleteMyRoutineTagsByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyRoutineTagsByIdsReqDto,
) (*dtos.HardDeleteMyRoutineTagsByIdsResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyRoutineTagsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
