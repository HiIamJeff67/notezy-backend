package services

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	materialsql "github.com/HiIamJeff67/notezy-backend/app/models/sqls/material"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	storages "github.com/HiIamJeff67/notezy-backend/app/storages"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type MaterialServiceInterface interface {
	GetMyMaterialById(ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception)
	GetMyMaterialAndItsParentById(ctx context.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto) (*dtos.GetMyMaterialAndItsParentByIdResDto, *exceptions.Exception)
	GetMyMaterialsByParentSubShelfId(ctx context.Context, reqDto *dtos.GetMyMaterialsByParentSubShelfIdReqDto) (*dtos.GetMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception)
	GetAllMyMaterialsByRootShelfId(ctx context.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto) (*dtos.GetAllMyMaterialsByRootShelfIdResDto, *exceptions.Exception)
	CreateMyMaterial(ctx context.Context, reqDto *dtos.CreateMyMaterialReqDto) (*dtos.CreateMyMaterialResDto, *exceptions.Exception)
	UpdateMyMaterialById(ctx context.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception)
	SaveMyMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialById(ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialsByIds(ctx context.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception)
	RestoreMyMaterialById(ctx context.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception)
	RestoreMyMaterialsByIds(ctx context.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception)
	DeleteMyMaterialById(ctx context.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception)
	DeleteMyMaterialsByIds(ctx context.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception)

	SearchPrivateMaterials(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchMaterialInput) (*gqlmodels.SearchMaterialConnection, *exceptions.Exception)
}

type MaterialService struct {
	db                 *gorm.DB
	storage            storages.StorageInterface
	materialScope      scopes.MaterialScopeInterface
	subShelfRepository repositories.SubShelfRepositoryInterface
	materialRepository repositories.MaterialRepositoryInterface
}

func NewMaterialService(
	db *gorm.DB,
	storage storages.StorageInterface,
	materialScope scopes.MaterialScopeInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	materialRepository repositories.MaterialRepositoryInterface,
) MaterialServiceInterface {
	return &MaterialService{
		db:                 db,
		storage:            storage,
		materialScope:      materialScope,
		subShelfRepository: subShelfRepository,
		materialRepository: materialRepository,
	}
}

func (s *MaterialService) GetMyMaterialById(
	ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto,
) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	material, exception := s.materialRepository.GetOneById(
		reqDto.Param.MaterialId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
	if exception != nil {
		exception.Log() // ignore the missing file in storage error
	}

	return &dtos.GetMyMaterialByIdResDto{
		Id:               material.Id,
		ParentSubShelfId: material.ParentSubShelfId,
		Name:             material.Name,
		Size:             material.Size,
		ContentType:      material.ContentType,
		ParseMediaType:   material.ParseMediaType,
		DownloadURL:      downloadURL,
		DeletedAt:        material.DeletedAt,
		UpdatedAt:        material.UpdatedAt,
		CreatedAt:        material.CreatedAt,
	}, nil
}

func (s *MaterialService) GetMyMaterialAndItsParentById(
	ctx context.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto,
) (*dtos.GetMyMaterialAndItsParentByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := dtos.GetMyMaterialAndItsParentByIdResDto{}
	var contentKey string
	err := db.Raw(materialsql.GetMyMaterialAndItsParentByIdSQL,
		reqDto.Param.MaterialId, reqDto.ContextFields.UserId, pg.Array(allowedPermissions), onlyDeleted,
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
		return nil, exceptions.Material.NotFound().WithOrigin(err)
	}
	if len(strings.TrimSpace(contentKey)) == 0 {
		return nil, exceptions.Material.NotFound()
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, contentKey, nil)
	if exception != nil {
		exception.Log() // ignore the missing file in storage error
	}
	resDto.DownloadURL = downloadURL // could be empty string

	return &resDto, nil
}

func (s *MaterialService) GetMyMaterialsByParentSubShelfId(
	ctx context.Context, reqDto *dtos.GetMyMaterialsByParentSubShelfIdReqDto,
) (*dtos.GetMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
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
			reqDto.Param.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		).Scopes(scopes.NewMaterialScope().FilterOnlyDeleted(onlyDeleted)).
		Order("name ASC").
		Limit(int(constants.MaxMaterialsOfSubShelf)).
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithOrigin(err)
	}

	resDto := dtos.GetMyMaterialsByParentSubShelfIdResDto{}
	for _, material := range materials {
		downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
		if exception != nil {
			exception.Log() // ignore the missing file in storage error
		}
		resDto = append(resDto, dtos.GetMyMaterialByIdResDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Size:             material.Size,
			ContentType:      material.ContentType,
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
	ctx context.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto,
) (*dtos.GetAllMyMaterialsByRootShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
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
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Scopes(scopes.NewMaterialScope().FilterOnlyDeleted(onlyDeleted)).
		Limit(int(constants.MaxMaterialsOfRootShelf)).
		Order("name ASC").
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound()
	}

	resDto := dtos.GetAllMyMaterialsByRootShelfIdResDto{}
	for _, material := range materials {
		downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
		if exception != nil {
			exception.Log() // ignore the missing file in storage error
		}
		resDto = append(resDto, dtos.GetMyMaterialByIdResDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Size:             material.Size,
			ContentType:      material.ContentType,
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
	ctx context.Context, reqDto *dtos.CreateMyMaterialReqDto,
) (*dtos.CreateMyMaterialResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	newMaterialId := uuid.New()
	newContentKey := s.storage.GetKey(
		reqDto.ContextFields.UserPublicId.String(),
		newMaterialId.String(),
	)
	zeroSize := int64(0)
	_, exception = s.materialRepository.CreateOneBySubShelfId(
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateMaterialInput{
			Id:             newMaterialId,
			Name:           reqDto.Body.Name,
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

	object, exception := s.storage.NewObject(newContentKey, newContentFile, zeroSize)
	if exception != nil {
		return nil, exception
	}

	exception = s.storage.PutObjectByKey(ctx, newContentKey, object)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateMyMaterialResDto{
		Id:        newMaterialId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) UpdateMyMaterialById(
	ctx context.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto,
) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	material, exception := s.materialRepository.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateMaterialInput{
			Values: inputs.UpdateMaterialInput{
				Name: reqDto.Body.Values.Name,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) SaveMyMaterialById(
	ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto,
) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}
	// check if there does exist a file in the reqDto
	if reqDto.Body.ContentFile == nil {
		return nil, exceptions.Material.InvalidDto()
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	partialUpdate := inputs.PartialUpdateMaterialInput{
		Values: inputs.UpdateMaterialInput{
			// content key remain the same here
		},
		SetNull: nil,
	}
	var contentKey = s.storage.GetKey(reqDto.ContextFields.UserPublicId.String(), reqDto.Body.MaterialId.String())

	var fileHeaderSize int64 = 0
	if reqDto.ContextFields.Size != nil {
		fileHeaderSize = *reqDto.ContextFields.Size
	}

	// extract the data in it and get its content type, parse media type, and actual size, etc.
	object, exception := s.storage.NewObject(contentKey, reqDto.Body.ContentFile, fileHeaderSize)
	if exception != nil {
		return nil, exception
	}
	if object == nil {
		return nil, exceptions.Material.CannotGetFileObjects()
	}

	size := object.Size
	contentType, err := enums.ConvertStringToMaterialContentType(object.ContentType)
	if err != nil {
		return nil, exceptions.Material.InvalidType(object.ContentType).WithOrigin(err)
	}
	partialUpdate.Values.ParseMediaType = object.ParseMediaType
	partialUpdate.Values.Size = &size
	partialUpdate.Values.ContentType = contentType

	material, exception := s.materialRepository.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		partialUpdate,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	// if there does exist a file, then put the file at the end to ensure the entire operation is consistent
	exception = s.storage.PutObjectByKey(ctx, material.ContentKey, object)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SaveMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) MoveMyMaterialById(
	ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto,
) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	result := db.Exec(materialsql.MoveMyMaterialByIdSQL,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		pg.Array(allowedPermissions),
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		pg.Array(allowedPermissions),
	)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Material.NoChanges()
	}

	return &dtos.MoveMyMaterialByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) MoveMyMaterialsByIds(
	ctx context.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto,
) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	result := db.Exec(materialsql.MoveMyMaterialsByIdsSQL,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
		pg.Array(allowedPermissions),
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		pg.Array(allowedPermissions),
	)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Material.NoChanges()
	}

	return &dtos.MoveMyMaterialsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) RestoreMyMaterialById(
	ctx context.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto,
) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredMaterial, exception := s.materialRepository.RestoreSoftDeletedOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, restoredMaterial.ContentKey, nil)
	if exception != nil {
		exception.Log() // ignore the missing file in storage error
	}

	return &dtos.RestoreMyMaterialByIdResDto{
		Id:               restoredMaterial.Id,
		ParentSubShelfId: restoredMaterial.ParentSubShelfId,
		Name:             restoredMaterial.Name,
		Size:             restoredMaterial.Size,
		ContentType:      restoredMaterial.ContentType,
		ParseMediaType:   restoredMaterial.ParseMediaType,
		DownloadURL:      downloadURL,
		DeletedAt:        restoredMaterial.DeletedAt,
		UpdatedAt:        restoredMaterial.UpdatedAt,
		CreatedAt:        restoredMaterial.CreatedAt,
	}, nil
}

func (s *MaterialService) RestoreMyMaterialsByIds(
	ctx context.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto,
) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredMaterials, exception := s.materialRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := dtos.RestoreMyMaterialsByIdsResDto{}
	for _, restoredMaterial := range restoredMaterials {
		downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, restoredMaterial.ContentKey, nil)
		if exception != nil {
			exception.Log() // ignore the missing file in storage error
		}
		resDto = append(resDto, dtos.RestoreMyMaterialByIdResDto{
			Id:               restoredMaterial.Id,
			ParentSubShelfId: restoredMaterial.ParentSubShelfId,
			Name:             restoredMaterial.Name,
			Size:             restoredMaterial.Size,
			ContentType:      restoredMaterial.ContentType,
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
	ctx context.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto,
) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.materialRepository.SoftDeleteOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyMaterialByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *MaterialService) DeleteMyMaterialsByIds(
	ctx context.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto,
) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	exception = s.materialRepository.SoftDeleteManyByIds(
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyMaterialsByIdsResDto{
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
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
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
		return nil, exceptions.Material.NotFound().WithOrigin(err)
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
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
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
