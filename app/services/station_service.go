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

type StationServiceInterface interface {
	GetOneById(ctx context.Context, reqDto *dtos.GetOneStationByIdReqDto) (*dtos.GetOneStationByIdResDto, *exceptions.Exception)
	CreateOneByOwnerId(ctx context.Context, reqDto *dtos.CreateOneStationByOwnerIdReqDto) (*dtos.CreateOneStationByOwnerIdResDto, *exceptions.Exception)
	CreateManyByOwnerId(ctx context.Context, reqDto *dtos.CreateManyStationsByOwnerIdReqDto) (*dtos.CreateManyStationsByOwnerIdResDto, *exceptions.Exception)
	UpdateOneById(ctx context.Context, reqDto *dtos.UpdateOneStationByIdReqDto) (*dtos.UpdateOneStationByIdResDto, *exceptions.Exception)
	BulkUpdateManyByIds(ctx context.Context, reqDto *dtos.BulkUpdateManyStationsByIdsReqDto) (*dtos.BulkUpdateManyStationsByIdsResDto, *exceptions.Exception)
	RestoreSoftDeletedOneById(ctx context.Context, reqDto *dtos.RestoreSoftDeletedOneStationByIdReqDto) (*dtos.RestoreSoftDeletedOneStationByIdResDto, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ctx context.Context, reqDto *dtos.RestoreSoftDeletedManyStationsByIdsReqDto) (*dtos.RestoreSoftDeletedManyStationsByIdsResDto, *exceptions.Exception)
	SoftDeleteOneById(ctx context.Context, reqDto *dtos.SoftDeleteOneStationByIdReqDto) (*dtos.SoftDeleteOneStationByIdResDto, *exceptions.Exception)
	SoftDeleteManyByIds(ctx context.Context, reqDto *dtos.SoftDeleteManyStationsByIdsReqDto) (*dtos.SoftDeleteManyStationsByIdsResDto, *exceptions.Exception)
	HardDeleteOneById(ctx context.Context, reqDto *dtos.HardDeleteOneStationByIdReqDto) (*dtos.HardDeleteOneStationByIdResDto, *exceptions.Exception)
	HardDeleteManyByIds(ctx context.Context, reqDto *dtos.HardDeleteManyStationsByIdsReqDto) (*dtos.HardDeleteManyStationsByIdsResDto, *exceptions.Exception)
}

type StationService struct {
	db                *gorm.DB
	stationRepository repositories.StationRepositoryInterface
}

func NewStationService(
	db *gorm.DB,
	stationRepository repositories.StationRepositoryInterface,
) StationServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &StationService{
		db:                db,
		stationRepository: stationRepository,
	}
}

/* ============================== Service Methods for Station ============================== */

func (s *StationService) GetOneById(
	ctx context.Context,
	reqDto *dtos.GetOneStationByIdReqDto,
) (*dtos.GetOneStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	onlyDeleted := types.Ternary_Negative
	if reqDto.Param.OnlyDeleted != nil {
		onlyDeleted = *reqDto.Param.OnlyDeleted
	}

	db := s.db.WithContext(ctx)
	station, permission, exception := s.stationRepository.GetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetOneStationByIdResDto{
		Id:                  station.Id,
		OwnerId:             station.OwnerId,
		Name:                station.Name,
		Description:         station.Description,
		Icon:                station.Icon,
		HeaderBackgroundURL: station.HeaderBackgroundURL,
		Permission:          permission,
		RoutineCount:        station.RoutineCount,
		DeletedAt:           station.DeletedAt,
		UpdatedAt:           station.UpdatedAt,
		CreatedAt:           station.CreatedAt,
	}, nil
}

func (s *StationService) CreateOneByOwnerId(
	ctx context.Context,
	reqDto *dtos.CreateOneStationByOwnerIdReqDto,
) (*dtos.CreateOneStationByOwnerIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	newStationId, exception := s.stationRepository.CreateOneByOwnerId(
		reqDto.ContextFields.UserId,
		inputs.CreateStationInput{
			Id:                  reqDto.Body.Id,
			Name:                reqDto.Body.Name,
			Description:         reqDto.Body.Description,
			Icon:                reqDto.Body.Icon,
			HeaderBackgroundURL: reqDto.Body.HeaderBackgroundURL,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateOneStationByOwnerIdResDto{
		Id:        *newStationId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *StationService) CreateManyByOwnerId(
	ctx context.Context,
	reqDto *dtos.CreateManyStationsByOwnerIdReqDto,
) (*dtos.CreateManyStationsByOwnerIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.CreateStationInput, len(reqDto.Body.CreatedStations))
	for index, createdStation := range reqDto.Body.CreatedStations {
		input[index] = inputs.CreateStationInput{
			Id:                  createdStation.Id,
			Name:                createdStation.Name,
			Description:         createdStation.Description,
			Icon:                createdStation.Icon,
			HeaderBackgroundURL: createdStation.HeaderBackgroundURL,
		}
	}
	newStationIds, exception := s.stationRepository.CreateManyByOwnerId(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateManyStationsByOwnerIdResDto{
		Ids:       newStationIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *StationService) UpdateOneById(
	ctx context.Context,
	reqDto *dtos.UpdateOneStationByIdReqDto,
) (*dtos.UpdateOneStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	updatedStation, exception := s.stationRepository.UpdateOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateStationInput{
			Values: inputs.UpdateStationInput{
				Name:                reqDto.Body.Values.Name,
				Description:         reqDto.Body.Values.Description,
				Icon:                reqDto.Body.Values.Icon,
				HeaderBackgroundURL: reqDto.Body.Values.HeaderBackgroundURL,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateOneStationByIdResDto{
		UpdatedAt: updatedStation.UpdatedAt,
	}, nil
}

func (s *StationService) BulkUpdateManyByIds(
	ctx context.Context,
	reqDto *dtos.BulkUpdateManyStationsByIdsReqDto,
) (*dtos.BulkUpdateManyStationsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkUpdateStationInput, len(reqDto.Body.UpdatedStations))
	for index, updatedStation := range reqDto.Body.UpdatedStations {
		input[index] = inputs.BulkUpdateStationInput{
			Id: updatedStation.StationId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateStationInput]{
				Values: inputs.UpdateStationInput{
					Name:                updatedStation.Values.Name,
					Description:         updatedStation.Values.Description,
					Icon:                updatedStation.Values.Icon,
					HeaderBackgroundURL: updatedStation.Values.HeaderBackgroundURL,
				},
				SetNull: updatedStation.SetNull,
			},
		}
	}
	exception := s.stationRepository.BulkUpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.BulkUpdateManyStationsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *StationService) RestoreSoftDeletedOneById(
	ctx context.Context,
	reqDto *dtos.RestoreSoftDeletedOneStationByIdReqDto,
) (*dtos.RestoreSoftDeletedOneStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	restoredStation, exception := s.stationRepository.RestoreSoftDeletedOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreSoftDeletedOneStationByIdResDto{
		Id:                  restoredStation.Id,
		OwnerId:             restoredStation.OwnerId,
		Name:                restoredStation.Name,
		Description:         restoredStation.Description,
		Icon:                restoredStation.Icon,
		HeaderBackgroundURL: restoredStation.HeaderBackgroundURL,
		RoutineCount:        restoredStation.RoutineCount,
		DeletedAt:           restoredStation.DeletedAt,
		UpdatedAt:           restoredStation.UpdatedAt,
		CreatedAt:           restoredStation.CreatedAt,
	}, nil
}

func (s *StationService) RestoreSoftDeletedManyByIds(
	ctx context.Context,
	reqDto *dtos.RestoreSoftDeletedManyStationsByIdsReqDto,
) (*dtos.RestoreSoftDeletedManyStationsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	restoredStations, exception := s.stationRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := dtos.RestoreSoftDeletedManyStationsByIdsResDto{}
	for _, restoredStation := range restoredStations {
		resDto = append(resDto, dtos.RestoreSoftDeletedOneStationByIdResDto{
			Id:                  restoredStation.Id,
			OwnerId:             restoredStation.OwnerId,
			Name:                restoredStation.Name,
			Description:         restoredStation.Description,
			Icon:                restoredStation.Icon,
			HeaderBackgroundURL: restoredStation.HeaderBackgroundURL,
			RoutineCount:        restoredStation.RoutineCount,
			DeletedAt:           restoredStation.DeletedAt,
			UpdatedAt:           restoredStation.UpdatedAt,
			CreatedAt:           restoredStation.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *StationService) SoftDeleteOneById(
	ctx context.Context,
	reqDto *dtos.SoftDeleteOneStationByIdReqDto,
) (*dtos.SoftDeleteOneStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.stationRepository.SoftDeleteOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SoftDeleteOneStationByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) SoftDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.SoftDeleteManyStationsByIdsReqDto,
) (*dtos.SoftDeleteManyStationsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.stationRepository.SoftDeleteManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SoftDeleteManyStationsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) HardDeleteOneById(
	ctx context.Context,
	reqDto *dtos.HardDeleteOneStationByIdReqDto,
) (*dtos.HardDeleteOneStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.stationRepository.HardDeleteOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteOneStationByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) HardDeleteManyByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteManyStationsByIdsReqDto,
) (*dtos.HardDeleteManyStationsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	exception := s.stationRepository.HardDeleteManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteManyStationsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL Station ============================== */
