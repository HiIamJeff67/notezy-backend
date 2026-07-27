package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
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
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type StationServiceInterface interface {
	GetMyStationById(ctx context.Context, reqDto *dtos.GetMyStationByIdReqDto) (*dtos.GetMyStationByIdResDto, *exceptions.Exception)
	GetAllMyStations(ctx context.Context, reqDto *dtos.GetAllMyStationsReqDto) (*dtos.GetAllMyStationsResDto, *exceptions.Exception)
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

	VisualizeMyTotalCount(ctx context.Context, reqDto *dtos.VisualizeMyTotalCountReqDto) (*dtos.VisualizeMyTotalCountResDto, *exceptions.Exception)

	GetMyStationPermission(ctx context.Context, reqDto *dtos.GetMyStationPermissionReqDto) (*dtos.GetMyStationPermissionResDto, *exceptions.Exception)
	CreateMyStationPermission(ctx context.Context, reqDto *dtos.CreateMyStationPermissionReqDto) (*dtos.CreateMyStationPermissionResDto, *exceptions.Exception)
	UpsertMyStationPermission(ctx context.Context, reqDto *dtos.UpsertMyStationPermissionReqDto) (*dtos.UpsertMyStationPermissionResDto, *exceptions.Exception)
	UpsertMyStationPermissions(ctx context.Context, reqDto *dtos.UpsertMyStationPermissionsReqDto) (*dtos.UpsertMyStationPermissionsResDto, *exceptions.Exception)
	UpdateMyStationPermission(ctx context.Context, reqDto *dtos.UpdateMyStationPermissionReqDto) (*dtos.UpdateMyStationPermissionResDto, *exceptions.Exception)
	TransferMyStationOwnership(ctx context.Context, reqDto *dtos.TransferMyStationOwnershipReqDto) (*dtos.TransferMyStationOwnershipResDto, *exceptions.Exception)
	DeleteMyStationPermission(ctx context.Context, reqDto *dtos.DeleteMyStationPermissionReqDto) *exceptions.Exception
	DeleteMyStationPermissions(ctx context.Context, reqDto *dtos.DeleteMyStationPermissionsReqDto) *exceptions.Exception
	LeaveMyStation(ctx context.Context, reqDto *dtos.LeaveMyStationReqDto) *exceptions.Exception
	LeaveMyStations(ctx context.Context, reqDto *dtos.LeaveMyStationsReqDto) *exceptions.Exception

	SearchPrivateStations(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchStationInput) (*gqlmodels.SearchStationConnection, *exceptions.Exception)
}

type StationService struct {
	db                        *gorm.DB
	stationScope              scopes.StationScopeInterface
	stationRepository         repositories.StationRepositoryInterface
	usersToStationsRepository repositories.UsersToStationsRepositoryInterface
}

func NewStationService(
	db *gorm.DB,
	stationScope scopes.StationScopeInterface,
	stationRepository repositories.StationRepositoryInterface,
	usersToStationsRepository repositories.UsersToStationsRepositoryInterface,
) StationServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &StationService{
		db:                        db,
		stationScope:              stationScope,
		stationRepository:         stationRepository,
		usersToStationsRepository: usersToStationsRepository,
	}
}

/* ============================== Auxiliary Functions ============================== */

func (s *StationService) saveMyStationPermission(
	ctx context.Context,
	actorUserId uuid.UUID,
	stationId uuid.UUID,
	targetUserPublicId uuid.UUID,
	permission enums.AccessControlPermission,
	requireExisting *bool,
) (*dtos.StationPermissionResDto, *exceptions.Exception) {
	if permission == enums.AccessControlPermission_Owner {
		return nil, exceptions.Station.NoPermission("transfer Station ownership through an access control")
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Station.FailedToBeginTransaction(
			"Failed to begin Station permission transaction",
		).WithOrigin(tx.Error)
	}
	station, actorPermission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		stationId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	var targetUser schemas.User
	if result := tx.Where("public_id = ?", targetUserPublicId).First(&targetUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	targetPermission, targetException := s.usersToStationsRepository.GetOne(
		station.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if targetException != nil && !errors.Is(targetException.Origin, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, targetException
	}
	if requireExisting != nil && *requireExisting != (targetPermission != nil) {
		tx.Rollback()
		if *requireExisting {
			return nil, targetException
		}
		return nil, exceptions.Station.NoChanges()
	}
	if targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Station.NoPermission("modify the Station owner")
	}
	if actorPermission != enums.AccessControlPermission_Owner && (permission == enums.AccessControlPermission_Admin || targetPermission != nil && targetPermission.Permission == enums.AccessControlPermission_Admin) {
		tx.Rollback()
		return nil, exceptions.Station.NoPermission("manage Admin permissions")
	}
	var relation *schemas.UsersToStations
	if targetPermission == nil {
		relation, exception = s.usersToStationsRepository.CreateOne(
			station.Id,
			targetUser.Id,
			permission,
			options.WithTransactionDB(tx),
		)
	} else {
		relation, exception = s.usersToStationsRepository.UpdateOne(
			station.Id,
			targetUser.Id,
			permission,
			options.WithTransactionDB(tx),
		)
	}
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}
	return &dtos.StationPermissionResDto{
		UserPublicId: targetUser.PublicId,
		Permission:   relation.Permission,
		UpdatedAt:    relation.UpdatedAt,
		CreatedAt:    relation.CreatedAt,
	}, nil
}

/* ============================== Service Methods for Station ============================== */

func (s *StationService) GetMyStationById(
	ctx context.Context,
	reqDto *dtos.GetMyStationByIdReqDto,
) (*dtos.GetMyStationByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	station, permission, exception := s.stationRepository.GetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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

func (s *StationService) GetAllMyStations(
	ctx context.Context,
	reqDto *dtos.GetAllMyStationsReqDto,
) (*dtos.GetAllMyStationsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	stations, permissions, exception := s.stationRepository.GetAllByUserId(
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyStationsResDto, len(stations))
	for index, station := range stations {
		resDto[index] = struct {
			Id                  uuid.UUID                     "json:\"id\""
			Name                string                        "json:\"name\""
			Icon                *enums.SupportedIcon          "json:\"icon\""
			HeaderBackgroundURL *string                       "json:\"headerBackgroundURL\""
			Permission          enums.AccessControlPermission "json:\"permission\""
			RoutineCount        int64                         "json:\"routineCount\""
			DeletedAt           *time.Time                    "json:\"deletedAt\""
			UpdatedAt           time.Time                     "json:\"updatedAt\""
			CreatedAt           time.Time                     "json:\"createdAt\""
		}{
			Id:                  station.Id,
			Name:                station.Name,
			Icon:                station.Icon,
			HeaderBackgroundURL: station.HeaderBackgroundURL,
			Permission:          permissions[index],
			RoutineCount:        station.RoutineCount,
			DeletedAt:           station.DeletedAt,
			UpdatedAt:           station.UpdatedAt,
			CreatedAt:           station.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *StationService) CreateStation(
	ctx context.Context,
	reqDto *dtos.CreateStationReqDto,
) (*dtos.CreateStationResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	newStationId, exception := s.stationRepository.CreateOne(
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
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
	newStationIds, exception := s.stationRepository.CreateMany(
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
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
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	input := make([]inputs.UpdateStationByIdInput, len(reqDto.Body.UpdatedStations))
	for index, updatedStation := range reqDto.Body.UpdatedStations {
		input[index] = inputs.UpdateStationByIdInput{
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
	exception = s.stationRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredStation, exception := s.stationRepository.RestoreSoftDeletedOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	restoredStations, exception := s.stationRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	station, permission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if permission == enums.AccessControlPermission_Owner {
		result := tx.
			Model(&schemas.Station{}).
			Where("id = ?", station.Id).
			Update("deleted_at", time.Now())
		if result.Error != nil {
			tx.Rollback()
			return nil, exceptions.Station.FailedToUpdate().WithOrigin(result.Error)
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, exceptions.Station.NoChanges()
		}
	} else {
		exception = s.usersToStationsRepository.DeleteOne(
			station.Id,
			reqDto.ContextFields.UserId,
			options.WithTransactionDB(tx),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	exception = s.stationRepository.SoftDeleteManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	exception = s.stationRepository.HardDeleteOneById(
		reqDto.Body.StationId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	exception = s.stationRepository.HardDeleteManyByIds(
		reqDto.Body.StationIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.HardDeleteMyStationsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for Visualization ============================== */

func (s *StationService) VisualizeMyTotalCount(
	ctx context.Context, reqDto *dtos.VisualizeMyTotalCountReqDto,
) (*dtos.VisualizeMyTotalCountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	var totals struct {
		StationCount     int64 `gorm:"column:station_count;"`
		RoutineCount     int64 `gorm:"column:routine_count;"`
		RoutineTaskCount int64 `gorm:"column:routine_task_count;"`
		RoutineTagCount  int64 `gorm:"column:routine_tag_count;"`
	}

	if reqDto.Param.Permission == enums.AccessControlPermission_Owner {
		result := db.Model(&schemas.UserAccount{}).
			Select("station_count, routine_count, routine_tag_count").
			Where(`user_id = ?`, reqDto.ContextFields.UserId).
			Scan(&totals)
		if result.Error != nil {
			return nil, exceptions.Station.NotFound().WithOrigin(result.Error)
		}

		result = db.Model(&schemas.RoutineTask{}).
			Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id`).
			Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
			Joins(`INNER JOIN "StationTable" station ON station.id = routine.station_id AND station.deleted_at IS NULL`).
			Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
			Count(&totals.RoutineTaskCount)
		if result.Error != nil {
			return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
		}

		return &dtos.VisualizeMyTotalCountResDto{
			Data: []dtos.TwoDimensionalDatum[int64]{
				dtos.TwoDimensionalDatum[int64]{
					Id:    "station-total-count",
					X:     "Station Total Count",
					Value: totals.StationCount,
				},
				dtos.TwoDimensionalDatum[int64]{
					Id:    "routine-total-count",
					X:     "Routine Total Count",
					Value: totals.RoutineCount,
				},
				dtos.TwoDimensionalDatum[int64]{
					Id:    "routine-task-total-count",
					X:     "Routine Task Total Count",
					Value: totals.RoutineTaskCount,
				},
				dtos.TwoDimensionalDatum[int64]{
					Id:    "routine-tag-total-count",
					X:     "Routine Tag Total Count",
					Value: totals.RoutineTagCount,
				},
			},
		}, nil
	}

	result := db.Model(&schemas.Station{}).
		Select(`
			COUNT(DISTINCT "StationTable".id) AS station_count,
			COALESCE(SUM("StationTable".routine_count), 0) AS routine_count
		`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = "StationTable".id`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Where(`"StationTable".deleted_at IS NULL`).
		Scan(&totals)
	if result.Error != nil {
		return nil, exceptions.Station.NotFound().WithOrigin(result.Error)
	}

	result = db.Model(&schemas.RoutineTask{}).
		Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id`).
		Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
		Joins(`INNER JOIN "StationTable" station ON station.id = routine.station_id AND station.deleted_at IS NULL`).
		Where("uts.user_id = ? AND uts.permission = ?", reqDto.ContextFields.UserId, reqDto.Param.Permission).
		Count(&totals.RoutineTaskCount)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.NotFound().WithOrigin(result.Error)
	}

	return &dtos.VisualizeMyTotalCountResDto{
		Data: []dtos.TwoDimensionalDatum[int64]{
			dtos.TwoDimensionalDatum[int64]{
				Id:    "station-total-count",
				X:     "Station Total Count",
				Value: totals.StationCount,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "routine-total-count",
				X:     "Routine Total Count",
				Value: totals.RoutineCount,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "routine-task-total-count",
				X:     "Routine Task Total Count",
				Value: totals.RoutineTaskCount,
			},
			dtos.TwoDimensionalDatum[int64]{
				Id:    "routine-tag-total-count",
				X:     "Routine Tag Total Count",
				Value: totals.RoutineTagCount,
			},
		},
	}, nil
}

/* ============================== Service Methods for Station Permissions ============================== */

func (s *StationService) GetMyStationPermission(
	ctx context.Context, reqDto *dtos.GetMyStationPermissionReqDto,
) (*dtos.GetMyStationPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)
	if _, _, exception = s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	); exception != nil {
		return nil, exception
	}

	var targetUser schemas.User
	if result := db.Where("public_id = ?", reqDto.Param.UserPublicId).First(&targetUser); result.Error != nil {
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	relation, exception := s.usersToStationsRepository.GetOne(
		reqDto.Param.StationId,
		targetUser.Id,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyStationPermissionResDto{UserPublicId: targetUser.PublicId, Permission: relation.Permission, UpdatedAt: relation.UpdatedAt, CreatedAt: relation.CreatedAt}, nil
}

func (s *StationService) CreateMyStationPermission(
	ctx context.Context, reqDto *dtos.CreateMyStationPermissionReqDto,
) (*dtos.CreateMyStationPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}
	requireExisting := false
	return s.saveMyStationPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.StationId, reqDto.Param.UserPublicId, reqDto.Body.Permission, &requireExisting)
}

func (s *StationService) UpsertMyStationPermission(
	ctx context.Context, reqDto *dtos.UpsertMyStationPermissionReqDto,
) (*dtos.UpsertMyStationPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}
	return s.saveMyStationPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.StationId, reqDto.Param.UserPublicId, reqDto.Body.Permission, nil)
}

func (s *StationService) UpsertMyStationPermissions(
	ctx context.Context, reqDto *dtos.UpsertMyStationPermissionsReqDto,
) (*dtos.UpsertMyStationPermissionsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	userPublicIds := make([]uuid.UUID, len(reqDto.Body.Permissions))
	permissionByPublicId := make(map[uuid.UUID]enums.AccessControlPermission, len(reqDto.Body.Permissions))
	for index, input := range reqDto.Body.Permissions {
		if input.Permission == enums.AccessControlPermission_Owner {
			return nil, exceptions.Station.NoPermission("transfer Station ownership through permissions")
		}
		if _, exists := permissionByPublicId[input.UserPublicId]; exists {
			return nil, exceptions.Station.InvalidDto("permissions cannot contain duplicate userPublicIds")
		}

		userPublicIds[index] = input.UserPublicId
		permissionByPublicId[input.UserPublicId] = input.Permission
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Station.FailedToBeginTransaction(
			"Failed to begin Station permission transaction",
		).WithOrigin(tx.Error)
	}

	station, actorPermission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", userPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if len(targetUsers) != len(userPublicIds) {
		tx.Rollback()
		return nil, exceptions.User.NotFound()
	}

	userByPublicId := make(map[uuid.UUID]schemas.User, len(targetUsers))
	userById := make(map[uuid.UUID]schemas.User, len(targetUsers))
	for _, user := range targetUsers {
		userByPublicId[user.PublicId] = user
		userById[user.Id] = user
	}

	userIds := make([]uuid.UUID, len(userPublicIds))
	for index, userPublicId := range userPublicIds {
		userIds[index] = userByPublicId[userPublicId].Id
	}

	existingPermissions, exception := s.usersToStationsRepository.GetMany(
		station.Id,
		userIds,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	existingPermissionByUserId := make(map[uuid.UUID]enums.AccessControlPermission, len(existingPermissions))
	for _, existingPermission := range existingPermissions {
		existingPermissionByUserId[existingPermission.UserId] = existingPermission.Permission
	}

	permissions := make([]enums.AccessControlPermission, len(userIds))
	for index, userId := range userIds {
		user := userById[userId]
		permission := permissionByPublicId[user.PublicId]
		if existingPermissionByUserId[userId] == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return nil, exceptions.Station.NoPermission("modify the Station owner")
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			(permission == enums.AccessControlPermission_Admin ||
				existingPermissionByUserId[userId] == enums.AccessControlPermission_Admin) {
			tx.Rollback()
			return nil, exceptions.Station.NoPermission("manage Admin permissions")
		}

		permissions[index] = permission
	}

	updatedPermissions, exception := s.usersToStationsRepository.UpsertMany(
		station.Id,
		userIds,
		permissions,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}

	updatedPermissionByUserId := make(map[uuid.UUID]schemas.UsersToStations, len(updatedPermissions))
	for _, updatedPermission := range updatedPermissions {
		updatedPermissionByUserId[updatedPermission.UserId] = updatedPermission
	}

	resDto := make([]dtos.UpsertMyStationPermissionResDto, len(userIds))
	for index, userId := range userIds {
		user := userById[userId]
		updatedPermission := updatedPermissionByUserId[userId]
		resDto[index] = dtos.UpsertMyStationPermissionResDto{
			UserPublicId: user.PublicId,
			Permission:   updatedPermission.Permission,
			UpdatedAt:    updatedPermission.UpdatedAt,
			CreatedAt:    updatedPermission.CreatedAt,
		}
	}

	return &dtos.UpsertMyStationPermissionsResDto{Permissions: resDto}, nil
}

func (s *StationService) UpdateMyStationPermission(
	ctx context.Context, reqDto *dtos.UpdateMyStationPermissionReqDto,
) (*dtos.UpdateMyStationPermissionResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}
	requireExisting := true
	return s.saveMyStationPermission(ctx, reqDto.ContextFields.UserId, reqDto.Param.StationId, reqDto.Param.UserPublicId, reqDto.Body.Permission, &requireExisting)
}

func (s *StationService) TransferMyStationOwnership(
	ctx context.Context,
	reqDto *dtos.TransferMyStationOwnershipReqDto,
) (*dtos.TransferMyStationOwnershipResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Station.FailedToBeginTransaction(
			"Failed to begin Station ownership transfer transaction",
		).WithOrigin(tx.Error)
	}
	station, permission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if permission != enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Station.NoPermission("transfer Station ownership")
	}

	var actorUser schemas.User
	if result := tx.Select("id, public_id").Where("id = ?", reqDto.ContextFields.UserId).First(&actorUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	var targetUser schemas.User
	if result := tx.Select("id, public_id").Where("public_id = ?", reqDto.Body.TargetUserPublicId).First(&targetUser); result.Error != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if targetUser.Id == reqDto.ContextFields.UserId {
		tx.Rollback()
		return nil, exceptions.Station.NoChanges()
	}

	targetMembership, exception := s.usersToStationsRepository.GetOne(
		station.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if targetMembership.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return nil, exceptions.Station.NoChanges()
	}

	var accounts []schemas.UserAccount
	result := tx.
		Clauses(clause.Locking{Strength: options.LockingStrengthUpdate}).
		Where("user_id IN ?", []uuid.UUID{reqDto.ContextFields.UserId, targetUser.Id}).
		Order("user_id").
		Find(&accounts)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.Station.FailedToUpdate().WithOrigin(result.Error)
	}
	if len(accounts) != 2 {
		tx.Rollback()
		return nil, exceptions.User.NotFound()
	}

	if _, exception = s.usersToStationsRepository.UpdateOne(
		station.Id,
		reqDto.ContextFields.UserId,
		enums.AccessControlPermission_Admin,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newOwnerMembership, exception := s.usersToStationsRepository.UpdateOne(
		station.Id,
		targetUser.Id,
		enums.AccessControlPermission_Owner,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	result = tx.Model(&schemas.Station{}).
		Where("id = ?", station.Id).
		Update("owner_id", targetUser.Id)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.Station.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, exceptions.Station.NotFound()
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.TransferMyStationOwnershipResDto{
		StationId:                 station.Id,
		PreviousOwnerUserPublicId: actorUser.PublicId,
		NewOwnerUserPublicId:      targetUser.PublicId,
		UpdatedAt:                 newOwnerMembership.UpdatedAt,
	}, nil
}

func (s *StationService) DeleteMyStationPermission(
	ctx context.Context, reqDto *dtos.DeleteMyStationPermissionReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Station.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Station.FailedToBeginTransaction(
			"Failed to begin Station permission transaction",
		).WithOrigin(tx.Error)
	}

	station, actorPermission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	var targetUser schemas.User
	result := tx.
		Model(&schemas.User{}).
		Where("public_id = ?", reqDto.Param.UserPublicId).
		First(&targetUser)
	if result.Error != nil {
		tx.Rollback()
		return exceptions.User.NotFound().WithOrigin(result.Error)
	}

	targetPermission, exception := s.usersToStationsRepository.GetOne(
		station.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if targetPermission.Permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return exceptions.Station.NoPermission("remove the Station owner")
	}
	if actorPermission != enums.AccessControlPermission_Owner &&
		targetPermission.Permission == enums.AccessControlPermission_Admin {
		tx.Rollback()
		return exceptions.Station.NoPermission("revoke Admin access")
	}

	exception = s.usersToStationsRepository.DeleteOne(
		station.Id,
		targetUser.Id,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}

	return nil
}

func (s *StationService) DeleteMyStationPermissions(
	ctx context.Context, reqDto *dtos.DeleteMyStationPermissionsReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Station.InvalidDto().WithOrigin(err)
	}

	userPublicIdSet := make(map[uuid.UUID]struct{}, len(reqDto.Body.UserPublicIds))
	for _, userPublicId := range reqDto.Body.UserPublicIds {
		if _, exists := userPublicIdSet[userPublicId]; exists {
			return exceptions.Station.InvalidDto("userPublicIds cannot contain duplicates")
		}

		userPublicIdSet[userPublicId] = struct{}{}
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return exception
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Station.FailedToBeginTransaction(
			"Failed to begin Station permission transaction",
		).WithOrigin(tx.Error)
	}

	station, actorPermission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	var targetUsers []schemas.User
	result := tx.
		Model(&schemas.User{}).
		Select("id, public_id").
		Where("public_id IN ?", reqDto.Body.UserPublicIds).
		Find(&targetUsers)
	if result.Error != nil {
		tx.Rollback()
		return exceptions.User.NotFound().WithOrigin(result.Error)
	}
	if len(targetUsers) != len(reqDto.Body.UserPublicIds) {
		tx.Rollback()
		return exceptions.User.NotFound()
	}

	userIdByPublicId := make(map[uuid.UUID]uuid.UUID, len(targetUsers))
	for _, targetUser := range targetUsers {
		userIdByPublicId[targetUser.PublicId] = targetUser.Id
	}

	userIds := make([]uuid.UUID, len(reqDto.Body.UserPublicIds))
	for index, userPublicId := range reqDto.Body.UserPublicIds {
		userIds[index] = userIdByPublicId[userPublicId]
	}

	targetPermissions, exception := s.usersToStationsRepository.GetMany(
		station.Id,
		userIds,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if len(targetPermissions) != len(userIds) {
		tx.Rollback()
		return exceptions.Station.NotFound()
	}

	for _, targetPermission := range targetPermissions {
		if targetPermission.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return exceptions.Station.NoPermission("remove the Station owner")
		}
		if actorPermission != enums.AccessControlPermission_Owner &&
			targetPermission.Permission == enums.AccessControlPermission_Admin {
			tx.Rollback()
			return exceptions.Station.NoPermission("revoke Admin access")
		}
	}

	exception = s.usersToStationsRepository.DeleteMany(
		station.Id,
		userIds,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}

	return nil
}

func (s *StationService) LeaveMyStation(
	ctx context.Context, reqDto *dtos.LeaveMyStationReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Station.InvalidDto().WithOrigin(err)
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Station.FailedToBeginTransaction("Failed to begin Station leave transaction").WithOrigin(tx.Error)
	}
	station, permission, exception := s.stationRepository.CheckPermissionAndGetOneById(
		reqDto.Param.StationId,
		reqDto.ContextFields.UserId,
		nil,
		nil,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if permission == enums.AccessControlPermission_Owner {
		tx.Rollback()
		return exceptions.Station.NoPermission("transfer Station ownership before leaving")
	}
	if exception = s.usersToStationsRepository.DeleteOne(
		station.Id,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}
	return nil
}

func (s *StationService) LeaveMyStations(
	ctx context.Context, reqDto *dtos.LeaveMyStationsReqDto,
) *exceptions.Exception {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return exceptions.Station.InvalidDto().WithOrigin(err)
	}
	stationIdSet := make(map[uuid.UUID]struct{}, len(reqDto.Body.Stations))
	stationIds := make([]uuid.UUID, len(reqDto.Body.Stations))
	for index, stationReqDto := range reqDto.Body.Stations {
		if _, exists := stationIdSet[stationReqDto.StationId]; exists {
			return exceptions.Station.InvalidDto("stations cannot contain duplicate stationIds")
		}
		stationIdSet[stationReqDto.StationId] = struct{}{}
		stationIds[index] = stationReqDto.StationId
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exceptions.Station.FailedToBeginTransaction("Failed to begin Station leave transaction").WithOrigin(tx.Error)
	}
	relations, exception := s.usersToStationsRepository.GetManyByStationIdsAndUserId(
		stationIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return exception
	}
	if len(relations) != len(stationIds) {
		tx.Rollback()
		return exceptions.Station.NotFound()
	}
	for _, relation := range relations {
		if relation.Permission == enums.AccessControlPermission_Owner {
			tx.Rollback()
			return exceptions.Station.NoPermission("transfer Station ownership before leaving")
		}
	}

	if exception = s.usersToStationsRepository.DeleteManyByStationIdsAndUserId(
		stationIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
	); exception != nil {
		tx.Rollback()
		return exception
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return exceptions.Station.FailedToCommitTransaction().WithOrigin(err)
	}
	return nil
}

/* ============================== Service Methods for GraphQL Station ============================== */

func (s *StationService) SearchPrivateStations(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchStationInput,
) (*gqlmodels.SearchStationConnection, *exceptions.Exception) {
	type PrivateStation struct {
		schemas.Station
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Negative
	if gqlInput.IsDeletedAt != nil && *gqlInput.IsDeletedAt {
		onlyDeleted = types.Ternary_Positive
	}

	query := db.Model(&schemas.Station{}).
		Select(`"StationTable".*, uts.permission AS permission`).
		Joins(`LEFT JOIN "UsersToStationsTable" uts ON "StationTable".id = uts.station_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.stationScope.FilterOnlyDeleted(onlyDeleted))

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

		query = query.Where("id > ?", searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchStationSortByName:
			query = query.Order("name " + cending).
				Order("routine_count " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByRoutineCount:
			query = query.Order("routine_count " + cending).
				Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByLastUpdate:
			query = query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("routine_count " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchStationSortByCreatedAt:
			query = query.Order("created_at " + cending).
				Order("name " + cending).
				Order("routine_count " + cending).
				Order("updated_at " + cending)
		default:
			query = query.Order("name " + cending).
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
			Node:                station.Station.ToPrivateSearchableStation(station.Permission),
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
