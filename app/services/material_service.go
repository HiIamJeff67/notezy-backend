package services

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	materialsql "notezy-backend/app/models/sql/material"
	storages "notezy-backend/app/storages"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type MaterialServiceInterface interface {
	GetMyMaterialById(ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception)
	GetMyMaterialAndItsParentById(ctx context.Context, reqDto *dtos.GetMyMaterialAndItsParentByIdReqDto) (*dtos.GetMyMaterialAndItsParentByIdResDto, *exceptions.Exception)
	GetAllMyMaterialsByParentSubShelfId(ctx context.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto) (*dtos.GetAllMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception)
	GetAllMyMaterialsByRootShelfId(ctx context.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto) (*dtos.GetAllMyMaterialsByRootShelfIdResDto, *exceptions.Exception)
	CreateTextbookMaterial(ctx context.Context, reqDto *dtos.CreateTextbookMaterialReqDto) (*dtos.CreateTextbookMaterialResDto, *exceptions.Exception)
	CreateNotebookMaterial(ctx context.Context, reqDto *dtos.CreateNotebookMaterialReqDto) (*dtos.CreateNotebookMaterialResDto, *exceptions.Exception)
	UpdateMyMaterialById(ctx context.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception)
	SaveMyTextbookMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception)
	SaveMyNotebookMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialById(ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialsByIds(ctx context.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception)
	RestoreMyMaterialById(ctx context.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception)
	RestoreMyMaterialsByIds(ctx context.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception)
	DeleteMyMaterialById(ctx context.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception)
	DeleteMyMaterialsByIds(ctx context.Context, reqDto *dtos.DeleteMyMaterialsByIdsReqDto) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception)
}

type MaterialService struct {
	db                 *gorm.DB
	storage            storages.StorageInterface
	subShelfRepository repositories.SubShelfRepositoryInterface
	materialRepository repositories.MaterialRepositoryInterface
}

func NewMaterialService(
	db *gorm.DB,
	storage storages.StorageInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	materialRepository repositories.MaterialRepositoryInterface,
) MaterialServiceInterface {
	return &MaterialService{
		db:                 db,
		storage:            storage,
		subShelfRepository: subShelfRepository,
		materialRepository: materialRepository,
	}
}

/* ============================== Service Methods for Material ============================== */

func (s *MaterialService) GetMyMaterialById(
	ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto,
) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	material, exception := s.materialRepository.GetOneById(
		db,
		reqDto.Param.MaterialId,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyMaterialByIdResDto{
		Id:               material.Id,
		ParentSubShelfId: material.ParentSubShelfId,
		Name:             material.Name,
		Type:             material.Type,
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
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}
	onlyDeleted := types.Ternary_Negative
	output := struct {
		Id                           uuid.UUID          `gorm:"id"`
		Name                         string             `gorm:"name"`
		Type                         enums.MaterialType `gorm:"type"`
		Size                         int64              `gorm:"size"`
		ContentKey                   string             `gorm:"content_key"`
		DeletedAt                    *time.Time         `gorm:"deleted_at"`
		UpdatedAt                    time.Time          `gorm:"updated_at"`
		CreatedAt                    time.Time          `gorm:"created_at"`
		RootShelfId                  uuid.UUID          `gorm:"root_shelf_id"`
		ParentSubShelfId             uuid.UUID          `gorm:"parent_sub_shelf_id"`
		ParentSubShelfName           string             `gorm:"parent_sub_shelf_name"`
		ParentSubShelfPrevSubShelfId *uuid.UUID         `gorm:"parent_sub_shelf_prev_sub_shelf_id"`
		ParentSubShelfPath           types.UUIDArray    `gorm:"parent_sub_shelf_path"`
		ParentSubShelfDeletedAt      time.Time          `gorm:"parent_sub_shelf_deleted_at"`
		ParentSubShelfUpdatedAt      time.Time          `gorm:"parent_sub_shelf_updated_at"`
		ParentSubShelfCreatedAt      time.Time          `gorm:"parent_sub_shelf_created_at"`
	}{}
	result := db.Raw(materialsql.GetMyMaterialAndItsParentByIdSQL,
		reqDto.Param.MaterialId, reqDto.ContextFields.UserId, allowedPermissions, onlyDeleted, onlyDeleted,
	).Scan(&output)
	if err := result.Error; err != nil || result.RowsAffected == 0 || len(strings.TrimSpace(output.ContentKey)) == 0 {
		exception := exceptions.Material.NotFound().WithError(err).Log()
		if err != nil {
			exception.WithError(err)
		}
		// should be exists in database
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, output.ContentKey, nil)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyMaterialAndItsParentByIdResDto{
		Id:                           output.Id,
		Name:                         output.Name,
		Type:                         output.Type,
		Size:                         output.Size,
		DownloadURL:                  downloadURL,
		DeletedAt:                    output.DeletedAt,
		UpdatedAt:                    output.UpdatedAt,
		CreatedAt:                    output.CreatedAt,
		RootShelfId:                  output.RootShelfId,
		ParentSubShelfId:             output.ParentSubShelfId,
		ParentSubShelfName:           output.ParentSubShelfName,
		ParentSubShelfPrevSubShelfId: output.ParentSubShelfPrevSubShelfId,
		ParentSubShelfPath:           output.ParentSubShelfPath,
		ParentSubShelfDeletedAt:      output.ParentSubShelfDeletedAt,
		ParentSubShelfUpdatedAt:      output.ParentSubShelfUpdatedAt,
		ParentSubShelfCreatedAt:      output.ParentSubShelfCreatedAt,
	}, nil
}

func (s *MaterialService) GetAllMyMaterialsByParentSubShelfId(
	ctx context.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto,
) (*dtos.GetAllMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials := []schemas.Material{}

	query := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		).Where("\"MaterialTable\".deleted_at IS NULL")

	result := query.Order("name ASC").
		Limit(int(constants.MaxMaterialsOfSubShelf)).
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	resDto := dtos.GetAllMyMaterialsByParentSubShelfIdResDto{}
	for _, material := range materials {
		downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
		if exception != nil {
			return nil, exception
		}
		resDto = append(resDto, dtos.GetMyMaterialByIdResDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Type:             material.Type,
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
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials := []schemas.Material{}

	result := db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Where("\"MaterialTable\".deleted_at IS NULL").
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
			return nil, exception
		}
		resDto = append(resDto, dtos.GetMyMaterialByIdResDto{
			Id:               material.Id,
			ParentSubShelfId: material.ParentSubShelfId,
			Name:             material.Name,
			Type:             material.Type,
			DownloadURL:      downloadURL,
			DeletedAt:        material.DeletedAt,
			UpdatedAt:        material.UpdatedAt,
			CreatedAt:        material.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *MaterialService) CreateTextbookMaterial(
	ctx context.Context, reqDto *dtos.CreateTextbookMaterialReqDto,
) (*dtos.CreateTextbookMaterialResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	newMaterialId := uuid.New()
	newContentKey := s.storage.GetKey(
		reqDto.ContextFields.UserPublicId.String(),
		newMaterialId.String(),
	)
	zeroSize := int64(0)
	_, exception := s.materialRepository.CreateOne(
		tx,
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateMaterialInput{
			Id:               newMaterialId,
			ParentSubShelfId: reqDto.Body.ParentSubShelfId,
			Name:             reqDto.Body.Name,
			Size:             zeroSize,
			Type:             enums.MaterialType_Textbook,
			ContentKey:       newContentKey,
		},
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newContentFile := bytes.NewReader([]byte{})

	object, exception := s.storage.NewObject(newContentKey, newContentFile, zeroSize)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	exception = s.storage.PutObjectByKey(ctx, newContentKey, object)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, newContentKey, nil)
	if exception != nil {
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Material.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.CreateTextbookMaterialResDto{
		Id:          newMaterialId,
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *MaterialService) CreateNotebookMaterial(
	ctx context.Context, reqDto *dtos.CreateNotebookMaterialReqDto,
) (*dtos.CreateNotebookMaterialResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	newMaterialId := uuid.New()
	newContentKey := s.storage.GetKey(
		reqDto.ContextFields.UserPublicId.String(),
		newMaterialId.String(),
	)
	zeroSize := int64(0)
	_, exception := s.materialRepository.CreateOne(
		tx,
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateMaterialInput{
			Id:               newMaterialId,
			ParentSubShelfId: reqDto.Body.ParentSubShelfId,
			Name:             reqDto.Body.Name,
			Size:             zeroSize,
			Type:             enums.MaterialType_Notebook,
			ContentKey:       newContentKey,
		},
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newContentFile := bytes.NewReader([]byte{})

	object, exception := s.storage.NewObject(newContentKey, newContentFile, zeroSize)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	exception = s.storage.PutObjectByKey(ctx, newContentKey, object)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, newContentKey, nil)
	if exception != nil {
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Material.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.CreateNotebookMaterialResDto{
		Id:          newMaterialId,
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *MaterialService) UpdateMyMaterialById(
	ctx context.Context, reqDto *dtos.UpdateMyMaterialByIdReqDto,
) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	material, exception := s.materialRepository.UpdateOneById(
		db,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		&reqDto.Body.MaterialType,
		inputs.PartialUpdateMaterialInput{
			Values: inputs.UpdateMaterialInput{
				Name: reqDto.Body.Values.Name,
			},
			SetNull: reqDto.Body.SetNull,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

// helper function for Save Material Services
func (s *MaterialService) saveMyMaterialById(
	ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto, materialType enums.MaterialType,
) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}
	// check if there does exist a file in the reqDto
	if reqDto.Body.ContentFile == nil {
		return nil, exceptions.Material.InvalidDto()
	}

	tx := s.db.WithContext(ctx).Begin()

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
		tx.Rollback()
		return nil, exception
	}
	if object == nil {
		tx.Rollback()
		return nil, exceptions.Material.CannotGetFileObjects()
	}

	partialUpdate.Values.ParseMediaType = object.ParseMediaType
	partialUpdate.Values.Size = &object.Size

	material, exception := s.materialRepository.UpdateOneById(
		tx,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		&materialType,
		partialUpdate,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// if there does exist a file, then put the file at the end to ensure the entire operation is consistent
	exception = s.storage.PutObjectByKey(ctx, material.ContentKey, object)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Material.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.SaveMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) SaveMyTextbookMaterialById(
	ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto,
) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	return s.saveMyMaterialById(ctx, reqDto, enums.MaterialType_Textbook)
}

func (s *MaterialService) SaveMyNotebookMaterialById(
	ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto,
) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	return s.saveMyMaterialById(ctx, reqDto, enums.MaterialType_Notebook)
}

func (s *MaterialService) MoveMyMaterialById(
	ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto,
) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := s.subShelfRepository.HasPermission(
		db,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		types.Ternary_Negative,
	); !hasPermission {
		return nil, exceptions.Shelf.NoPermission()
	}
	material, exception := s.materialRepository.CheckPermissionAndGetOneById(
		db,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		types.Ternary_Negative, // exclude the deleted materials
	)
	if exception != nil {
		return nil, exception
	}

	result := db.Exec(`
		UPDATE "MaterialTable"
		SET "parent_sub_shelf_id" = ?, "updated_at" = NOW()
		WHERE id = ? AND deleted_at IS NULL`,
		reqDto.Body.DestinationParentSubShelfId, material.Id,
	)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithError(err)
	}

	return &dtos.MoveMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) MoveMyMaterialsByIds(
	ctx context.Context, reqDto *dtos.MoveMyMaterialsByIdsReqDto,
) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials, exception := s.materialRepository.CheckPermissionsAndGetManyByIds(
		db,
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		types.Ternary_Negative, // exclude the deleted materials
	)
	if exception != nil {
		return nil, exception
	}
	materialIds := []uuid.UUID{}
	for _, material := range materials {
		materialIds = append(materialIds, material.Id)
	}

	result := db.Exec(`
		UPDATE "MaterialTable"
		SET "parent_sub_shelf_id" = ?, "updated_at" = NOW()
		WHERE id IN ? AND deleted_at IS NULL`,
		reqDto.Body.DestinationParentSubShelfId, materialIds, // use the extracted material to update
	)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithError(err)
	}

	return &dtos.MoveMyMaterialsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) RestoreMyMaterialById(
	ctx context.Context, reqDto *dtos.RestoreMyMaterialByIdReqDto,
) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.materialRepository.RestoreSoftDeletedOneById(
		db,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyMaterialByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) RestoreMyMaterialsByIds(
	ctx context.Context, reqDto *dtos.RestoreMyMaterialsByIdsReqDto,
) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.materialRepository.RestoreSoftDeletedManyByIds(
		db,
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyMaterialsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) DeleteMyMaterialById(
	ctx context.Context, reqDto *dtos.DeleteMyMaterialByIdReqDto,
) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.materialRepository.SoftDeleteOneById(
		db,
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
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
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	exception := s.materialRepository.SoftDeleteManyByIds(
		db,
		reqDto.Body.MaterialIds,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyMaterialsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
