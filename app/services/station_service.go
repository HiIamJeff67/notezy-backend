package services

import (
	"context"
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
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type StationServiceInterface interface {
	GetMyStationById(ctx context.Context, reqDto *dtos.GetMyStationByIdReqDto) (*dtos.GetMyStationByIdResDto, *exceptions.Exception)
	CreateStation(ctx context.Context, reqDto *dtos.CreateStationReqDto) (*dtos.CreateStationResDto, *exceptions.Exception)
	CreateStations(ctx context.Context, reqDto *dtos.CreateStationsReqDto) (*dtos.CreateStationsResDto, *exceptions.Exception)
	UpdateMyStationById(ctx context.Context, reqDto *dtos.UpdateMyStationByIdReqDto) (*dtos.UpdateMyStationByIdResDto, *exceptions.Exception)
	UpdateMyStationsByIds(ctx context.Context, reqDto *dtos.UpdateMyStationsByIdsReqDto) (*dtos.UpdateMyStationsByIdsResDto, *exceptions.Exception)
	RestoreMyStationById(ctx context.Context, reqDto *dtos.RestoreMyStationByIdReqDto) (*dtos.RestoreMyStationByIdResDto, *exceptions.Exception)
	RestoreMyStationsByIds(ctx context.Context, reqDto *dtos.RestoreMyStationsByIdsReqDto) (*dtos.RestoreMyStationsByIdsResDto, *exceptions.Exception)
	DeleteMyStationById(ctx context.Context, reqDto *dtos.DeleteMyStationByIdReqDto) (*dtos.DeleteMyStationByIdResDto, *exceptions.Exception)
	DeleteMyStationsByIds(ctx context.Context, reqDto *dtos.DeleteMyStationsByIdsReqDto) (*dtos.DeleteMyStationsByIdsResDto, *exceptions.Exception)
	HardDeleteMyStationById(ctx context.Context, reqDto *dtos.HardDeleteMyStationByIdReqDto) (*dtos.HardDeleteMyStationByIdResDto, *exceptions.Exception)
	HardDeleteMyStationsByIds(ctx context.Context, reqDto *dtos.HardDeleteMyStationsByIdsReqDto) (*dtos.HardDeleteMyStationsByIdsResDto, *exceptions.Exception)

	// services for graphql stations
	SearchPrivateStations(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchStationInput) (*gqlmodels.SearchStationConnection, *exceptions.Exception)
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

func (s *StationService) GetMyStationById(
	ctx context.Context,
	reqDto *dtos.GetMyStationByIdReqDto,
) (*dtos.GetMyStationByIdResDto, *exceptions.Exception) {
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

	return &dtos.GetMyStationByIdResDto{
		Id:                  station.Id,
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

func (s *StationService) CreateStation(
	ctx context.Context,
	reqDto *dtos.CreateStationReqDto,
) (*dtos.CreateStationResDto, *exceptions.Exception) {
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

	return &dtos.CreateStationResDto{
		Id:        *newStationId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *StationService) CreateStations(
	ctx context.Context,
	reqDto *dtos.CreateStationsReqDto,
) (*dtos.CreateStationsResDto, *exceptions.Exception) {
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

	return &dtos.CreateStationsResDto{
		Ids:       newStationIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *StationService) UpdateMyStationById(
	ctx context.Context,
	reqDto *dtos.UpdateMyStationByIdReqDto,
) (*dtos.UpdateMyStationByIdResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyStationByIdResDto{
		UpdatedAt: updatedStation.UpdatedAt,
	}, nil
}

func (s *StationService) UpdateMyStationsByIds(
	ctx context.Context,
	reqDto *dtos.UpdateMyStationsByIdsReqDto,
) (*dtos.UpdateMyStationsByIdsResDto, *exceptions.Exception) {
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

	return &dtos.UpdateMyStationsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *StationService) RestoreMyStationById(
	ctx context.Context,
	reqDto *dtos.RestoreMyStationByIdReqDto,
) (*dtos.RestoreMyStationByIdResDto, *exceptions.Exception) {
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

	return &dtos.RestoreMyStationByIdResDto{
		Id:                  restoredStation.Id,
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

func (s *StationService) RestoreMyStationsByIds(
	ctx context.Context,
	reqDto *dtos.RestoreMyStationsByIdsReqDto,
) (*dtos.RestoreMyStationsByIdsResDto, *exceptions.Exception) {
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

	resDto := dtos.RestoreMyStationsByIdsResDto{}
	for _, restoredStation := range restoredStations {
		resDto = append(resDto, dtos.RestoreMyStationByIdResDto{
			Id:                  restoredStation.Id,
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

func (s *StationService) DeleteMyStationById(
	ctx context.Context,
	reqDto *dtos.DeleteMyStationByIdReqDto,
) (*dtos.DeleteMyStationByIdResDto, *exceptions.Exception) {
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

	return &dtos.DeleteMyStationByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) DeleteMyStationsByIds(
	ctx context.Context,
	reqDto *dtos.DeleteMyStationsByIdsReqDto,
) (*dtos.DeleteMyStationsByIdsResDto, *exceptions.Exception) {
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

	return &dtos.DeleteMyStationsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) HardDeleteMyStationById(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyStationByIdReqDto,
) (*dtos.HardDeleteMyStationByIdResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyStationByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *StationService) HardDeleteMyStationsByIds(
	ctx context.Context,
	reqDto *dtos.HardDeleteMyStationsByIdsReqDto,
) (*dtos.HardDeleteMyStationsByIdsResDto, *exceptions.Exception) {
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

	return &dtos.HardDeleteMyStationsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL Station ============================== */

func (s *StationService) SearchPrivateStations(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchStationInput,
) (*gqlmodels.SearchStationConnection, *exceptions.Exception) {
	type PrivateStation struct {
		schemas.Station
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	allowedPermisssions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.Station{}).
		Select("\"StationTable\".*, uts.permission AS permission").
		Joins("LEFT JOIN \"UsersToStationsTable\" uts ON \"StationTable\".id = uts.station_id").
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermisssions).
		Where("\"StationTable\".deleted_at IS NULL")

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchStationCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchStationSortByName:
			query.Order("name " + cending).
				Order("routine_count " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByRoutineCount:
			query.Order("routine_count " + cending).
				Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("routine_count " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("routine_count " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("name " + cending).
				Order("routine_count " + cending).
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

	var stations []PrivateStation
	if err := query.Find(&stations).Error; err != nil {
		return nil, exceptions.Station.NotFound().WithOrigin(err)
	}

	hasNextPage := len(stations) > limit
	searchEdges := make([]*gqlmodels.SearchStationEdge, len(stations))

	for index, station := range stations {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchStationCursorFields]{
			Fields: gqlmodels.SearchStationCursorFields{
				ID: station.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchStationEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                station.Station.ToPrivateStation(station.Permission),
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

	return &gqlmodels.SearchStationConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
