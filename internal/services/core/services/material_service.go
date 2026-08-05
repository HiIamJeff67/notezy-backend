package services

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	materialsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/materials"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	materialsql "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/sqls/material"
	storage "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/storage"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
	validator "github.com/go-playground/validator/v10"
)

type MaterialServiceInterface interface {
	GetMyMaterialById(ctx context.Context, requestDto *materialsdto.GetMyMaterialByIdRequestDto) (*materialsdto.GetMyMaterialByIdResponseDto, *exceptions.Exception)
	GetMyMaterialAndItsParentById(ctx context.Context, requestDto *materialsdto.GetMyMaterialAndItsParentByIdRequestDto) (*materialsdto.GetMyMaterialAndItsParentByIdResponseDto, *exceptions.Exception)
	GetMyMaterialsByParentSubShelfId(ctx context.Context, requestDto *materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto) (*materialsdto.GetMyMaterialsByParentSubShelfIdResponseDto, *exceptions.Exception)
	GetAllMyMaterialsByRootShelfId(ctx context.Context, requestDto *materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto) (*materialsdto.GetAllMyMaterialsByRootShelfIdResponseDto, *exceptions.Exception)
	CreateMyMaterial(ctx context.Context, requestDto *materialsdto.CreateMyMaterialRequestDto) (*materialsdto.CreateMyMaterialResponseDto, *exceptions.Exception)
	UpdateMyMaterialById(ctx context.Context, requestDto *materialsdto.UpdateMyMaterialByIdRequestDto) (*materialsdto.UpdateMyMaterialByIdResponseDto, *exceptions.Exception)
	SaveMyMaterialById(ctx context.Context, requestDto *materialsdto.SaveMyMaterialByIdRequestDto) (*materialsdto.SaveMyMaterialByIdResponseDto, *exceptions.Exception)
	MoveMyMaterialById(ctx context.Context, requestDto *materialsdto.MoveMyMaterialByIdRequestDto) (*materialsdto.MoveMyMaterialByIdResponseDto, *exceptions.Exception)
	MoveMyMaterialsByIds(ctx context.Context, requestDto *materialsdto.MoveMyMaterialsByIdsRequestDto) (*materialsdto.MoveMyMaterialsByIdsResponseDto, *exceptions.Exception)
	RestoreMyMaterialById(ctx context.Context, requestDto *materialsdto.RestoreMyMaterialByIdRequestDto) (*materialsdto.RestoreMyMaterialByIdResponseDto, *exceptions.Exception)
	RestoreMyMaterialsByIds(ctx context.Context, requestDto *materialsdto.RestoreMyMaterialsByIdsRequestDto) (*materialsdto.RestoreMyMaterialsByIdsResponseDto, *exceptions.Exception)
	DeleteMyMaterialById(ctx context.Context, requestDto *materialsdto.DeleteMyMaterialByIdRequestDto) (*materialsdto.DeleteMyMaterialByIdResponseDto, *exceptions.Exception)
	DeleteMyMaterialsByIds(ctx context.Context, requestDto *materialsdto.DeleteMyMaterialsByIdsRequestDto) (*materialsdto.DeleteMyMaterialsByIdsResponseDto, *exceptions.Exception)

	SearchPrivateMaterials(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchMaterialInput) (*gqlmodels.SearchMaterialConnection, *exceptions.Exception)
}

type MaterialService struct {
	validator          *validator.Validate
	db                 *gorm.DB
	storage            storage.StorageInterface
	materialScope      scopes.MaterialScopeInterface
	subShelfRepository repositories.SubShelfRepositoryInterface
	materialRepository repositories.MaterialRepositoryInterface
	storageKeySalt     string
}

func NewMaterialService(
	validator *validator.Validate,
	db *gorm.DB,
	storage storage.StorageInterface,
	materialScope scopes.MaterialScopeInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	materialRepository repositories.MaterialRepositoryInterface,
	storageKeySalt string,
) MaterialServiceInterface {
	return &MaterialService{
		validator:          validator,
		db:                 db,
		storage:            storage,
		materialScope:      materialScope,
		subShelfRepository: subShelfRepository,
		materialRepository: materialRepository,
		storageKeySalt:     storageKeySalt,
	}
}

func (s *MaterialService) GetMyMaterialById(
	ctx context.Context, requestDto *materialsdto.GetMyMaterialByIdRequestDto,
) (*materialsdto.GetMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.IsDeleted != nil {
		if *requestDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	material, exception := s.materialRepository.GetOneById(
		requestDto.Param.MaterialId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	downloadURL, err := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
	if err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
	}

	return &materialsdto.GetMyMaterialByIdResponseDto{
		Id:               material.Id,
		ParentSubShelfId: material.ParentSubShelfId,
		Name:             material.Name,
		Size:             material.Size,
		ContentType:      *material.ContentType.ToContractable(),
		ParseMediaType:   material.ParseMediaType,
		DownloadURL:      downloadURL,
		DeletedAt:        material.DeletedAt,
		UpdatedAt:        material.UpdatedAt,
		CreatedAt:        material.CreatedAt,
	}, nil
}

func (s *MaterialService) GetMyMaterialAndItsParentById(
	ctx context.Context, requestDto *materialsdto.GetMyMaterialAndItsParentByIdRequestDto,
) (*materialsdto.GetMyMaterialAndItsParentByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.IsDeleted != nil {
		if *requestDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := materialsdto.GetMyMaterialAndItsParentByIdResponseDto{}
	var contentKey string
	err := db.Raw(materialsql.GetMyMaterialAndItsParentByIdSQL,
		requestDto.Param.MaterialId, actorUserId, pg.Array(allowedPermissions), onlyDeleted,
	).Row().
		Scan(&resDto.Id,
			&resDto.Name,
			&resDto.Size,
			&resDto.ContentType,
			&resDto.ParseMediaType,
			&contentKey,
			&resDto.DeletedAt,
			&resDto.UpdatedAt,
			&resDto.CreatedAt,
			&resDto.RootShelfId,
			&resDto.ParentSubShelfId,
			&resDto.ParentSubShelfName,
			&resDto.ParentSubShelfPrevSubShelfId,
			&resDto.ParentSubShelfPath,
			&resDto.ParentSubShelfDeletedAt,
			&resDto.ParentSubShelfUpdatedAt,
			&resDto.ParentSubShelfCreatedAt,
		)
	if err != nil {
		return nil, apiexceptions.Material.NotFound().WithOrigin(err)
	}
	if len(strings.TrimSpace(contentKey)) == 0 {
		return nil, apiexceptions.Material.NotFound()
	}

	downloadURL, err := s.storage.PresignGetObjectByKey(ctx, contentKey, nil)
	if err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
	}
	resDto.DownloadURL = downloadURL // could be empty string

	return &resDto, nil
}

func (s *MaterialService) GetMyMaterialsByParentSubShelfId(
	ctx context.Context, requestDto *materialsdto.GetMyMaterialsByParentSubShelfIdRequestDto,
) (*materialsdto.GetMyMaterialsByParentSubShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.AreDeleted != nil {
		if *requestDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	materials := []schemas.Material{}
	result := db.Model(&schemas.Material{}).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "MaterialTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			requestDto.Param.ParentSubShelfId,
			actorUserId,
			allowedPermissions,
		).Scopes(scopes.NewMaterialScope().FilterOnlyDeleted(onlyDeleted)).
		Order("name ASC").
		Limit(int(data.MaxMaterialsOfSubShelf)).
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, apiexceptions.Material.NotFound().WithOrigin(err)
	}

	resDto := materialsdto.GetMyMaterialsByParentSubShelfIdResponseDto{}
	for _, material := range materials {
		downloadURL, err := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
		if err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
		}
		resDto = append(resDto, materialsdto.GetMyMaterialByIdResponseDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Size:             material.Size,
			ContentType:      *material.ContentType.ToContractable(),
			ParseMediaType:   material.ParseMediaType,
			DownloadURL:      downloadURL,
			DeletedAt:        material.DeletedAt,
			UpdatedAt:        material.UpdatedAt,
			CreatedAt:        material.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *MaterialService) GetAllMyMaterialsByRootShelfId(
	ctx context.Context, requestDto *materialsdto.GetAllMyMaterialsByRootShelfIdRequestDto,
) (*materialsdto.GetAllMyMaterialsByRootShelfIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if requestDto.Param.AreDeleted != nil {
		if *requestDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	materials := []schemas.Material{}
	result := db.Model(&schemas.Material{}).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "MaterialTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			requestDto.Param.RootShelfId, actorUserId, allowedPermissions,
		).Scopes(scopes.NewMaterialScope().FilterOnlyDeleted(onlyDeleted)).
		Limit(int(data.MaxMaterialsOfRootShelf)).
		Order("name ASC").
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, apiexceptions.Material.NotFound()
	}

	resDto := materialsdto.GetAllMyMaterialsByRootShelfIdResponseDto{}
	for _, material := range materials {
		downloadURL, err := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
		if err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
		}
		resDto = append(resDto, materialsdto.GetMyMaterialByIdResponseDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Size:             material.Size,
			ContentType:      *material.ContentType.ToContractable(),
			ParseMediaType:   material.ParseMediaType,
			DownloadURL:      downloadURL,
			DeletedAt:        material.DeletedAt,
			UpdatedAt:        material.UpdatedAt,
			CreatedAt:        material.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *MaterialService) CreateMyMaterial(
	ctx context.Context, requestDto *materialsdto.CreateMyMaterialRequestDto,
) (*materialsdto.CreateMyMaterialResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserPublicId, exception := contexts.GetActorUserPublicId(ctx)
	if exception != nil {
		return nil, exception
	}

	newMaterialId := uuid.New()
	newContentKey := s.storage.GetKey(
		actorUserPublicId.String(),
		newMaterialId.String(),
		s.storageKeySalt,
	)
	zeroSize := int64(0)
	_, exception = s.materialRepository.CreateOneBySubShelfId(
		requestDto.Body.ParentSubShelfId,
		actorUserId,
		inputs.CreateMaterialInput{
			Id:             newMaterialId,
			Name:           requestDto.Body.Name,
			Size:           zeroSize,
			ContentKey:     newContentKey,
			ParseMediaType: "",
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	newContentFile := bytes.NewReader([]byte{})

	object, err := s.storage.NewObject(newContentKey, newContentFile, zeroSize)
	if err != nil {
		return nil, apiexceptions.Storage.FailedToReadObjectBytes().WithOrigin(err)
	}

	if err := s.storage.PutObjectByKey(ctx, newContentKey, object); err != nil {
		return nil, apiexceptions.Storage.FailedToPutObject(object).WithOrigin(err)
	}

	return &materialsdto.CreateMyMaterialResponseDto{
		Id:        newMaterialId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) UpdateMyMaterialById(
	ctx context.Context, requestDto *materialsdto.UpdateMyMaterialByIdRequestDto,
) (*materialsdto.UpdateMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	material, exception := s.materialRepository.UpdateOneById(
		requestDto.Param.MaterialId,
		actorUserId,
		inputs.PartialUpdateMaterialInput{
			Values: inputs.UpdateMaterialInput{
				Name: requestDto.Body.Values.Name,
			},
			SetNull: requestDto.Body.SetNull,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &materialsdto.UpdateMyMaterialByIdResponseDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) SaveMyMaterialById(
	ctx context.Context, requestDto *materialsdto.SaveMyMaterialByIdRequestDto,
) (*materialsdto.SaveMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}
	// check if there does exist a file in the requestDto
	if requestDto.Body.ContentFile == nil {
		return nil, apiexceptions.Material.InvalidDto()
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserPublicId, exception := contexts.GetActorUserPublicId(ctx)
	if exception != nil {
		return nil, exception
	}

	partialUpdate := inputs.PartialUpdateMaterialInput{
		Values: inputs.UpdateMaterialInput{
			// content key remain the same here
		},
		SetNull: nil,
	}
	var contentKey = s.storage.GetKey(
		actorUserPublicId.String(),
		requestDto.Param.MaterialId.String(),
		s.storageKeySalt,
	)

	fileHeaderSize := int64(len(requestDto.Body.ContentFile))

	// extract the data in it and get its content type, parse media type, and actual size, etc.
	object, err := s.storage.NewObject(contentKey, bytes.NewReader(requestDto.Body.ContentFile), fileHeaderSize)
	if err != nil {
		return nil, apiexceptions.Storage.FailedToReadObjectBytes().WithOrigin(err)
	}
	if object == nil {
		return nil, apiexceptions.Material.CannotGetFileObjects()
	}

	size := object.Size
	contentType, err := enums.ConvertStringToMaterialContentType(object.ContentType)
	if err != nil {
		return nil, apiexceptions.Material.InvalidType(object.ContentType).WithOrigin(err)
	}
	partialUpdate.Values.ParseMediaType = object.ParseMediaType
	partialUpdate.Values.Size = &size
	partialUpdate.Values.ContentType = contentType

	material, exception := s.materialRepository.UpdateOneById(
		requestDto.Param.MaterialId,
		actorUserId,
		partialUpdate,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	// if there does exist a file, then put the file at the end to ensure the entire operation is consistent
	if err := s.storage.PutObjectByKey(ctx, material.ContentKey, object); err != nil {
		return nil, apiexceptions.Storage.FailedToPutObject(object).WithOrigin(err)
	}

	return &materialsdto.SaveMyMaterialByIdResponseDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) MoveMyMaterialById(
	ctx context.Context, requestDto *materialsdto.MoveMyMaterialByIdRequestDto,
) (*materialsdto.MoveMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	result := db.Exec(materialsql.MoveMyMaterialByIdSQL,
		requestDto.Body.DestinationParentSubShelfId,
		requestDto.Body.MaterialId,
		actorUserId,
		pg.Array(allowedPermissions),
		requestDto.Body.DestinationParentSubShelfId,
		actorUserId,
		pg.Array(allowedPermissions),
	)
	if err := result.Error; err != nil {
		return nil, apiexceptions.Material.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, apiexceptions.Material.NoChanges()
	}

	return &materialsdto.MoveMyMaterialByIdResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) MoveMyMaterialsByIds(
	ctx context.Context, requestDto *materialsdto.MoveMyMaterialsByIdsRequestDto,
) (*materialsdto.MoveMyMaterialsByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	result := db.Exec(materialsql.MoveMyMaterialsByIdsSQL,
		requestDto.Body.DestinationParentSubShelfId,
		requestDto.Body.MaterialIds,
		actorUserId,
		pg.Array(allowedPermissions),
		requestDto.Body.DestinationParentSubShelfId,
		actorUserId,
		pg.Array(allowedPermissions),
	)
	if err := result.Error; err != nil {
		return nil, apiexceptions.Material.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, apiexceptions.Material.NoChanges()
	}

	return &materialsdto.MoveMyMaterialsByIdsResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) RestoreMyMaterialById(
	ctx context.Context, requestDto *materialsdto.RestoreMyMaterialByIdRequestDto,
) (*materialsdto.RestoreMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredMaterial, exception := s.materialRepository.RestoreSoftDeletedOneById(
		requestDto.Param.MaterialId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	downloadURL, err := s.storage.PresignGetObjectByKey(ctx, restoredMaterial.ContentKey, nil)
	if err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
	}

	return &materialsdto.RestoreMyMaterialByIdResponseDto{
		Id:               restoredMaterial.Id,
		ParentSubShelfId: restoredMaterial.ParentSubShelfId,
		Name:             restoredMaterial.Name,
		Size:             restoredMaterial.Size,
		ContentType:      *restoredMaterial.ContentType.ToContractable(),
		ParseMediaType:   restoredMaterial.ParseMediaType,
		DownloadURL:      downloadURL,
		DeletedAt:        restoredMaterial.DeletedAt,
		UpdatedAt:        restoredMaterial.UpdatedAt,
		CreatedAt:        restoredMaterial.CreatedAt,
	}, nil
}

func (s *MaterialService) RestoreMyMaterialsByIds(
	ctx context.Context, requestDto *materialsdto.RestoreMyMaterialsByIdsRequestDto,
) (*materialsdto.RestoreMyMaterialsByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredMaterials, exception := s.materialRepository.RestoreSoftDeletedManyByIds(
		requestDto.Body.MaterialIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := materialsdto.RestoreMyMaterialsByIdsResponseDto{}
	for _, restoredMaterial := range restoredMaterials {
		downloadURL, err := s.storage.PresignGetObjectByKey(ctx, restoredMaterial.ContentKey, nil)
		if err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to presign Material object")
		}
		resDto = append(resDto, materialsdto.RestoreMyMaterialByIdResponseDto{
			Id:               restoredMaterial.Id,
			ParentSubShelfId: restoredMaterial.ParentSubShelfId,
			Name:             restoredMaterial.Name,
			Size:             restoredMaterial.Size,
			ContentType:      *restoredMaterial.ContentType.ToContractable(),
			ParseMediaType:   restoredMaterial.ParseMediaType,
			DownloadURL:      downloadURL,
			DeletedAt:        restoredMaterial.DeletedAt,
			UpdatedAt:        restoredMaterial.UpdatedAt,
			CreatedAt:        restoredMaterial.CreatedAt,
		})
	}
	return &resDto, nil
}

func (s *MaterialService) DeleteMyMaterialById(
	ctx context.Context, requestDto *materialsdto.DeleteMyMaterialByIdRequestDto,
) (*materialsdto.DeleteMyMaterialByIdResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.materialRepository.SoftDeleteOneById(
		requestDto.Param.MaterialId,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &materialsdto.DeleteMyMaterialByIdResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *MaterialService) DeleteMyMaterialsByIds(
	ctx context.Context, requestDto *materialsdto.DeleteMyMaterialsByIdsRequestDto,
) (*materialsdto.DeleteMyMaterialsByIdsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, apiexceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.materialRepository.SoftDeleteManyByIds(
		requestDto.Body.MaterialIds,
		actorUserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &materialsdto.DeleteMyMaterialsByIdsResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL Material ============================== */

func (s *MaterialService) SearchPrivateMaterials(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchMaterialInput,
) (*gqlmodels.SearchMaterialConnection, *exceptions.Exception) {
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

	query := db.Model(&schemas.Material{}).
		Select(`"MaterialTable".*`).
		Joins(`INNER JOIN "SubShelfTable" ss ON ss.id = "MaterialTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.materialScope.FilterOnlyDeleted(onlyDeleted))

	if gqlInput.ParentSubShelfID != nil {
		query = query.Where(`"MaterialTable".parent_sub_shelf_id = ?`, *gqlInput.ParentSubShelfID)
	}

	if gqlInput.RootShelfID != nil {
		query = query.Where("ss.root_shelf_id = ?", *gqlInput.RootShelfID)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			`"MaterialTable".name ILIKE ? OR "MaterialTable".content_type::text ILIKE ? OR "MaterialTable".parse_media_type ILIKE ?`,
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchMaterialCursorFields](*gqlInput.After)
		if err != nil {
			return nil, apiexceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"MaterialTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchMaterialSortByName:
			query = query.Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".updated_at ` + cending).
				Order(`"MaterialTable".created_at ` + cending)
		case gqlmodels.SearchMaterialSortBySize:
			query = query.Order(`"MaterialTable".size ` + cending).
				Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".updated_at ` + cending).
				Order(`"MaterialTable".created_at ` + cending)
		case gqlmodels.SearchMaterialSortByContentType:
			query = query.Order(`"MaterialTable".content_type ` + cending).
				Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".updated_at ` + cending).
				Order(`"MaterialTable".created_at ` + cending)
		case gqlmodels.SearchMaterialSortByLastUpdate:
			query = query.Order(`"MaterialTable".updated_at ` + cending).
				Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".created_at ` + cending)
		case gqlmodels.SearchMaterialSortByCreatedAt:
			query = query.Order(`"MaterialTable".created_at ` + cending).
				Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".updated_at ` + cending)
		default:
			query = query.Order(`"MaterialTable".name ` + cending).
				Order(`"MaterialTable".updated_at ` + cending).
				Order(`"MaterialTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var materials []schemas.Material
	if err := query.Find(&materials).Error; err != nil {
		return nil, apiexceptions.Material.NotFound().WithOrigin(err)
	}

	hasNextPage := len(materials) > limit
	searchEdges := make([]*gqlmodels.SearchMaterialEdge, len(materials))

	for index, material := range materials {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchMaterialCursorFields]{
			Fields: gqlmodels.SearchMaterialCursorFields{
				ID: material.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, apiexceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, apiexceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchMaterialEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                material.ToPrivateMaterial(),
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

	return &gqlmodels.SearchMaterialConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
