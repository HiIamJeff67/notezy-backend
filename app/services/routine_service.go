package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineServiceInterface interface {
	GetMyRoutineById(ctx context.Context, reqDto *dtos.GetMyRoutineByIdReqDto) (*dtos.GetMyRoutineByIdResDto, *exceptions.Exception)
	GetAllMyRoutinesByTimeRange(ctx context.Context, reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto) (*dtos.GetAllMyRoutinesByTimeRangeResDto, *exceptions.Exception)
	CreateRoutineByStationId(ctx context.Context, reqDto *dtos.CreateRoutineByStationIdReqDto) (*dtos.CreateRoutineByStationIdResDto, *exceptions.Exception)
	CreateRoutinesByStationIds(ctx context.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto) (*dtos.CreateRoutinesByStationIdsResDto, *exceptions.Exception)
	UpdateMyRoutineById(ctx context.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto) (*dtos.UpdateMyRoutineByIdResDto, *exceptions.Exception)
	UpdateMyRoutinesByIds(ctx context.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto) (*dtos.UpdateMyRoutinesByIdsResDto, *exceptions.Exception)
	LinkRoutineTagById(ctx context.Context, reqDto *dtos.LinkRoutineTagByIdReqDto) (*dtos.LinkRoutineTagByIdResDto, *exceptions.Exception)
	BulkLinkRoutineTagsByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineTagsByIdsReqDto) (*dtos.BulkLinkRoutineTagsByIdsResDto, *exceptions.Exception)
	LinkRoutineTaskById(ctx context.Context, reqDto *dtos.LinkRoutineTaskByIdReqDto) (*dtos.LinkRoutineTaskByIdResDto, *exceptions.Exception)
	BulkLinkRoutineTasksByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineTasksByIdsReqDto) (*dtos.BulkLinkRoutineTasksByIdsResDto, *exceptions.Exception)
	LinkRoutineItemById(ctx context.Context, reqDto *dtos.LinkRoutineItemByIdReqDto) (*dtos.LinkRoutineItemByIdResDto, *exceptions.Exception)
	BulkLinkRoutineItemsByIds(ctx context.Context, reqDto *dtos.BulkLinkRoutineItemsByIdsReqDto) (*dtos.BulkLinkRoutineItemsByIdsResDto, *exceptions.Exception)
	RestoreMyRoutineById(ctx context.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto) (*dtos.RestoreMyRoutineByIdResDto, *exceptions.Exception)
	RestoreMyRoutinesByIds(ctx context.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto) (*dtos.RestoreMyRoutinesByIdsResDto, *exceptions.Exception)
	DeleteMyRoutineById(ctx context.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto) (*dtos.DeleteMyRoutineByIdResDto, *exceptions.Exception)
	DeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto) (*dtos.DeleteMyRoutinesByIdsResDto, *exceptions.Exception)
	HardDeleteMyRoutineById(ctx context.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto) (*dtos.HardDeleteMyRoutineByIdResDto, *exceptions.Exception)
	HardDeleteMyRoutinesByIds(ctx context.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto) (*dtos.HardDeleteMyRoutinesByIdsResDto, *exceptions.Exception)

	// services for graphql routines
	SearchPrivateRoutines(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineInput) (*gqlmodels.SearchRoutineConnection, *exceptions.Exception)
}

type RoutineService struct {
	db                    *gorm.DB
	stationRepository     repositories.StationRepositoryInterface
	routineRepository     repositories.RoutineRepositoryInterface
	routineTagRepository  repositories.RoutineTagRepositoryInterface
	routineTaskRepository repositories.RoutineTaskRepositoryInterface
	itemRepository        repositories.ItemRepositoryInterface
}

func NewRoutineService(
	db *gorm.DB,
	stationRepository repositories.StationRepositoryInterface,
	routineRepository repositories.RoutineRepositoryInterface,
	routineTagRepository repositories.RoutineTagRepositoryInterface,
	routineTaskRepository repositories.RoutineTaskRepositoryInterface,
	itemRepository repositories.ItemRepositoryInterface,
) RoutineServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &RoutineService{
		db:                    db,
		stationRepository:     stationRepository,
		routineRepository:     routineRepository,
		routineTagRepository:  routineTagRepository,
		routineTaskRepository: routineTaskRepository,
		itemRepository:        itemRepository,
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

func (s *RoutineService) GetAllMyRoutinesByTimeRange(
	ctx context.Context,
	reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto,
) (*dtos.GetAllMyRoutinesByTimeRangeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}
	if !reqDto.Param.From.Before(reqDto.Param.To) {
		return nil, exceptions.Routine.InvalidInput().WithOrigin(fmt.Errorf("from must be before to"))
	}

	db := s.db.WithContext(ctx)
	routines, exception := s.routineRepository.GetAllByTimeRange(
		reqDto.Param.From,
		reqDto.Param.To,
		reqDto.Param.StationIds,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.GetAllMyRoutinesByTimeRangeResDto, len(routines))
	for index, routine := range routines {
		resDto[index] = dtos.GetMyRoutineByIdResDto{
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
		}
	}

	return &resDto, nil
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

func (s *RoutineService) LinkRoutineTagById(
	ctx context.Context,
	reqDto *dtos.LinkRoutineTagByIdReqDto,
) (*dtos.LinkRoutineTagByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.routineTagRepository.HasPermission(
		reqDto.Body.RoutineTagId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.RoutineTag.NoPermission("get the routine tag")
	}

	var newRoutinesToTags schemas.RoutinesToTags
	newRoutinesToTags.RoutineId = reqDto.Body.RoutineId
	newRoutinesToTags.TagId = reqDto.Body.RoutineTagId

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Where("routine_id = ? AND tag_id = ?", newRoutinesToTags.RoutineId, newRoutinesToTags.TagId).
			Delete(&schemas.RoutinesToTags{})
	} else {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Create(&newRoutinesToTags)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTags().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineTagByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineTagsByIds(
	ctx context.Context,
	reqDto *dtos.BulkLinkRoutineTagsByIdsReqDto,
) (*dtos.BulkLinkRoutineTagsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	isRoutineExist := make(map[uuid.UUID]bool)
	isRoutineTagExist := make(map[uuid.UUID]bool)
	var routineIds []uuid.UUID
	var routineTagIds []uuid.UUID
	for _, linkedRoutineAndTag := range reqDto.Body.LinkedRoutinesAndTags {
		if !isRoutineExist[linkedRoutineAndTag.RoutineId] {
			isRoutineExist[linkedRoutineAndTag.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndTag.RoutineId)
		}
		if !isRoutineTagExist[linkedRoutineAndTag.RoutineTagId] {
			isRoutineTagExist[linkedRoutineAndTag.RoutineTagId] = true
			routineTagIds = append(routineTagIds, linkedRoutineAndTag.RoutineTagId)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isRoutineTagValid := make(map[uuid.UUID]bool)
	validRoutineTags, exception := s.routineTagRepository.CheckPermissionsAndGetManyByIds(
		routineTagIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutineTag := range validRoutineTags {
		isRoutineTagValid[validRoutineTag.Id] = true
	}

	var newRoutinesToTags []schemas.RoutinesToTags
	for _, linkedRoutineAndTag := range reqDto.Body.LinkedRoutinesAndTags {
		if !isRoutineValid[linkedRoutineAndTag.RoutineId] ||
			!isRoutineTagValid[linkedRoutineAndTag.RoutineTagId] {
			continue
		}
		newRoutinesToTags = append(newRoutinesToTags, schemas.RoutinesToTags{
			RoutineId: linkedRoutineAndTag.RoutineId,
			TagId:     linkedRoutineAndTag.RoutineTagId,
		})
	}
	if len(newRoutinesToTags) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToTags))
	for index, newRoutineToTag := range newRoutinesToTags {
		values[index] = []any{newRoutineToTag.RoutineId, newRoutineToTag.TagId}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Where("(routine_id, tag_id) IN ?", values).
			Delete(&schemas.RoutinesToTags{})
	} else {
		result = tx.Model(&schemas.RoutinesToTags{}).
			Create(&newRoutinesToTags)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTags().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineTagsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) LinkRoutineTaskById(
	ctx context.Context,
	reqDto *dtos.LinkRoutineTaskByIdReqDto,
) (*dtos.LinkRoutineTaskByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.routineTaskRepository.HasPermission(
		reqDto.Body.RoutineTaskId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.RoutineTask.NoPermission("get the routine task")
	}

	var newRoutinesToTasks schemas.RoutinesToTasks
	newRoutinesToTasks.RoutineId = reqDto.Body.RoutineId
	newRoutinesToTasks.TaskId = reqDto.Body.RoutineTaskId

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Where("routine_id = ? AND task_id = ?", newRoutinesToTasks.RoutineId, newRoutinesToTasks.TaskId).
			Delete(&schemas.RoutinesToTasks{})
	} else {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Create(&newRoutinesToTasks)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTasks().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineTaskByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineTasksByIds(
	ctx context.Context,
	reqDto *dtos.BulkLinkRoutineTasksByIdsReqDto,
) (*dtos.BulkLinkRoutineTasksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	isRoutineExist := make(map[uuid.UUID]bool)
	isRoutineTaskExist := make(map[uuid.UUID]bool)
	var routineIds []uuid.UUID
	var routineTaskIds []uuid.UUID
	for _, linkedRoutineAndTask := range reqDto.Body.LinkedRoutinesAndTasks {
		if !isRoutineExist[linkedRoutineAndTask.RoutineId] {
			isRoutineExist[linkedRoutineAndTask.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndTask.RoutineId)
		}
		if !isRoutineTaskExist[linkedRoutineAndTask.RoutineTaskId] {
			isRoutineTaskExist[linkedRoutineAndTask.RoutineTaskId] = true
			routineTaskIds = append(routineTaskIds, linkedRoutineAndTask.RoutineTaskId)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isRoutineTaskValid := make(map[uuid.UUID]bool)
	validRoutineTasks, exception := s.routineTaskRepository.CheckPermissionsAndGetManyByIds(
		routineTaskIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutineTag := range validRoutineTasks {
		isRoutineTaskValid[validRoutineTag.Id] = true
	}

	var newRoutinesToTasks []schemas.RoutinesToTasks
	for _, linkedRoutineAndTask := range reqDto.Body.LinkedRoutinesAndTasks {
		if !isRoutineValid[linkedRoutineAndTask.RoutineId] ||
			!isRoutineTaskValid[linkedRoutineAndTask.RoutineTaskId] {
			continue
		}
		newRoutinesToTasks = append(newRoutinesToTasks, schemas.RoutinesToTasks{
			RoutineId: linkedRoutineAndTask.RoutineId,
			TaskId:    linkedRoutineAndTask.RoutineTaskId,
		})
	}
	if len(newRoutinesToTasks) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToTasks))
	for index, newRoutineToTask := range newRoutinesToTasks {
		values[index] = []any{newRoutineToTask.RoutineId, newRoutineToTask.TaskId}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Where("(routine_id, task_id) IN ?", values).
			Delete(&schemas.RoutinesToTasks{})
	} else {
		result = tx.Model(&schemas.RoutinesToTasks{}).
			Create(&newRoutinesToTasks)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkRoutineTasks().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineTasksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) LinkRoutineItemById(
	ctx context.Context,
	reqDto *dtos.LinkRoutineItemByIdReqDto,
) (*dtos.LinkRoutineItemByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	if !s.routineRepository.HasPermission(
		reqDto.Body.RoutineId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Routine.NoPermission("get the routine")
	}

	if !s.itemRepository.HasPermission(
		reqDto.Body.ItemId,
		reqDto.Body.ItemType,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Item.NoPermission("get the item")
	}

	var newRoutinesToItems schemas.RoutinesToItems
	newRoutinesToItems.RoutineId = reqDto.Body.RoutineId
	newRoutinesToItems.ItemId = reqDto.Body.ItemId
	newRoutinesToItems.ItemType = reqDto.Body.ItemType

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Where(
				"routine_id = ? AND item_id = ? AND item_type = ?",
				newRoutinesToItems.RoutineId,
				newRoutinesToItems.ItemId,
				newRoutinesToItems.ItemType,
			).
			Delete(&schemas.RoutinesToItems{})
	} else {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Create(&newRoutinesToItems)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkItems().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.LinkRoutineItemByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *RoutineService) BulkLinkRoutineItemsByIds(
	ctx context.Context,
	reqDto *dtos.BulkLinkRoutineItemsByIdsReqDto,
) (*dtos.BulkLinkRoutineItemsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithOrigin(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := s.db.WithContext(ctx).Begin()

	isRoutineExist := make(map[uuid.UUID]bool)
	isItemExist := make(map[types.Pair[uuid.UUID, enums.ItemType]]bool)
	var routineIds []uuid.UUID
	var itemIdentities []types.Pair[uuid.UUID, enums.ItemType]
	for _, linkedRoutineAndItem := range reqDto.Body.LinkedRoutinesAndItems {
		if !isRoutineExist[linkedRoutineAndItem.RoutineId] {
			isRoutineExist[linkedRoutineAndItem.RoutineId] = true
			routineIds = append(routineIds, linkedRoutineAndItem.RoutineId)
		}
		itemIdentity := types.Pair[uuid.UUID, enums.ItemType]{
			First:  linkedRoutineAndItem.ItemId,
			Second: linkedRoutineAndItem.ItemType,
		}
		if !isItemExist[itemIdentity] {
			isItemExist[itemIdentity] = true
			itemIdentities = append(itemIdentities, itemIdentity)
		}
	}

	isRoutineValid := make(map[uuid.UUID]bool)
	validRoutines, exception := s.routineRepository.CheckPermissionsAndGetManyByIds(
		routineIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validRoutine := range validRoutines {
		isRoutineValid[validRoutine.Id] = true
	}

	isItemValid := make(map[types.Pair[uuid.UUID, enums.ItemType]]bool)
	validItems, exception := s.itemRepository.CheckPermissionsAndGetManyByIds(
		itemIdentities,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for _, validItem := range validItems {
		isItemValid[types.Pair[uuid.UUID, enums.ItemType]{
			First:  validItem.Id,
			Second: validItem.Type,
		}] = true
	}

	var newRoutinesToItems []schemas.RoutinesToItems
	for _, linkedRoutineAndItem := range reqDto.Body.LinkedRoutinesAndItems {
		itemIdentity := types.Pair[uuid.UUID, enums.ItemType]{
			First:  linkedRoutineAndItem.ItemId,
			Second: linkedRoutineAndItem.ItemType,
		}
		if !isRoutineValid[linkedRoutineAndItem.RoutineId] ||
			!isItemValid[itemIdentity] {
			continue
		}
		newRoutinesToItems = append(newRoutinesToItems, schemas.RoutinesToItems{
			RoutineId: linkedRoutineAndItem.RoutineId,
			ItemId:    linkedRoutineAndItem.ItemId,
			ItemType:  linkedRoutineAndItem.ItemType,
		})
	}
	if len(newRoutinesToItems) == 0 {
		tx.Rollback()
		return nil, exceptions.Routine.NoChanges()
	}

	values := make([][]any, len(newRoutinesToItems))
	for index, newRoutineToItem := range newRoutinesToItems {
		values[index] = []any{newRoutineToItem.RoutineId, newRoutineToItem.ItemId, newRoutineToItem.ItemType}
	}

	var result *gorm.DB
	if reqDto.Body.IsUnlink {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Where("(routine_id, item_id, item_type) IN ?", values).
			Delete(&schemas.RoutinesToItems{})
	} else {
		result = tx.Model(&schemas.RoutinesToItems{}).
			Create(&newRoutinesToItems)
	}
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Routine.FailedToLinkItems().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Routine.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Routine.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BulkLinkRoutineItemsByIdsResDto{
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

/* ============================== Service Methods for GraphQL Routine ============================== */

func (s *RoutineService) SearchPrivateRoutines(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchRoutineInput,
) (*gqlmodels.SearchRoutineConnection, *exceptions.Exception) {
	type PrivateRoutine struct {
		schemas.Routine
		Permission enums.AccessControlPermission `gorm:"column:permission"`
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.Routine{}).
		Select("\"RoutineTable\".*, uts.permission AS permission").
		Joins("LEFT JOIN \"UsersToStationsTable\" uts ON \"RoutineTable\".station_id = uts.station_id").
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Where("\"RoutineTable\".deleted_at IS NULL")

	if gqlInput.StationID != nil {
		query = query.Where(
			"\"RoutineTable\".station_id = ?",
			*gqlInput.StationID,
		)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"title ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchRoutineSortByTitle:
			query.Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByStatus:
			query.Order("status " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByScheduledStartAt:
			query.Order("scheduled_start_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByScheduledEndAt:
			query.Order("scheduled_end_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByPeriod:
			query.Order("period " + cending).
				Order("title " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("title " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchRoutineSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("title " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("title " + cending).
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

	var routines []PrivateRoutine
	if err := query.Find(&routines).Error; err != nil {
		return nil, exceptions.Routine.NotFound().WithOrigin(err)
	}

	hasNextPage := len(routines) > limit
	searchEdges := make([]*gqlmodels.SearchRoutineEdge, len(routines))

	for index, routine := range routines {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchRoutineCursorFields]{
			Fields: gqlmodels.SearchRoutineCursorFields{
				ID: routine.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchRoutineEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                routine.Routine.ToPrivateRoutine(),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: false,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchRoutineConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
