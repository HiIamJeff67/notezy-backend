package services

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	storages "notezy-backend/app/storages"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type MaterialServiceInterface interface {
	GetMyMaterialById(ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception)
	GetAllMyMaterialsByParentSubShelfId(ctx context.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto) (*dtos.GetAllMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception)
	GetAllMyMaterialsByRootShelfId(ctx context.Context, reqDto *dtos.GetAllMyMaterialsByRootShelfIdReqDto) (*dtos.GetAllMyMaterialsByRootShelfIdResDto, *exceptions.Exception)
	CreateTextbookMaterial(ctx context.Context, reqDto *dtos.CreateMaterialReqDto) (*dtos.CreateMaterialResDto, *exceptions.Exception)
	UpdateMyTextbookMaterialById(reqDto *dtos.UpdateMyMaterialByIdReqDto) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception)
	SaveMyTextbookMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialById(reqDto *dtos.MoveMyMaterialByIdReqDto) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialsByIds(reqDto *dtos.MoveMyMaterialsByIdsReqDto) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception)
	RestoreMyMaterialById(reqDto *dtos.RestoreMyMaterialByIdReqDto) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception)
	RestoreMyMaterialsByIds(reqDto *dtos.RestoreMyMaterialsByIdsReqDto) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception)
	DeleteMyMaterialById(reqDto *dtos.DeleteMyMaterialByIdReqDto) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception)
	DeleteMyMaterialsByIds(reqDto *dtos.DeleteMyMaterialsByIdsReqDto) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception)
}

type MaterialService struct {
	db      *gorm.DB
	storage storages.StorageInterface
}

func NewMaterialService(db *gorm.DB, storage storages.StorageInterface) MaterialServiceInterface {
	return &MaterialService{
		db:      db,
		storage: storage,
	}
}

/* ============================== Service Methods for Material ============================== */

func (s *MaterialService) GetMyMaterialById(
	ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto,
) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	material, exception := materialRepository.GetOneById(
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
		ContentType:      material.ContentType,
		ParseMediaType:   material.ParseMediaType,
		DeletedAt:        material.DeletedAt,
		UpdatedAt:        material.UpdatedAt,
		CreatedAt:        material.CreatedAt,
	}, nil
}

func (s *MaterialService) GetAllMyMaterialsByParentSubShelfId(
	ctx context.Context, reqDto *dtos.GetAllMyMaterialsByParentSubShelfIdReqDto,
) (*dtos.GetAllMyMaterialsByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials := []schemas.Material{}

	query := s.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		)

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
			ContentType:      material.ContentType,
			ParseMediaType:   material.ParseMediaType,
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

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials := []schemas.Material{}

	result := s.db.Model(&schemas.Material{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"MaterialTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Limit(int(constants.MaxMaterialsOfRootShelf)).
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
			ContentType:      material.ContentType,
			ParseMediaType:   material.ParseMediaType,
			DeletedAt:        material.DeletedAt,
			UpdatedAt:        material.UpdatedAt,
			CreatedAt:        material.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *MaterialService) CreateTextbookMaterial(ctx context.Context, reqDto *dtos.CreateMaterialReqDto) (*dtos.CreateMaterialResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	tx := s.db.Begin()
	materialRepository := repositories.NewMaterialRepository(tx)

	newMaterialId := uuid.New()
	newContentKey := s.storage.GetKey(
		reqDto.ContextFields.UserPublicId.String(),
		newMaterialId.String(),
	)
	zeroSize := int64(0)
	_, exception := materialRepository.CreateOne(
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateMaterialInput{
			Id:               newMaterialId,
			ParentSubShelfId: reqDto.Body.ParentSubShelfId,
			Name:             reqDto.Body.Name,
			Size:             zeroSize,
			Type:             enums.MaterialType_Textbook,
			ContentKey:       newContentKey,
			ContentType:      enums.MaterialContentType_PlainText,
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

	return &dtos.CreateMaterialResDto{
		Id:          newMaterialId,
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *MaterialService) UpdateMyTextbookMaterialById(reqDto *dtos.UpdateMyMaterialByIdReqDto,
) (*dtos.UpdateMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialResitory := repositories.NewMaterialRepository(s.db)

	materialType := enums.MaterialType_Textbook

	material, exception := materialResitory.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		&materialType,
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

func (s *MaterialService) SaveMyTextbookMaterialById(
	ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto,
) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}
	// check if there does exist a file in the reqDto
	if reqDto.Body.ContentFile == nil {
		return nil, exceptions.Material.InvalidDto()
	}

	tx := s.db.Begin()
	materialRepository := repositories.NewMaterialRepository(tx)

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

	// check if the material content type is allowed in the material type of textbook
	materialType := enums.MaterialType_Textbook
	if !materialType.IsContentTypeStringAllowed(object.ContentType) {
		tx.Rollback()
		return nil, exceptions.Material.MaterialContentTypeNotAllowedInMaterialType(
			reqDto.Body.MaterialId.String(),
			materialType.String(),
			object.ContentType,
			materialType.AllowedContentTypeStrings(),
		)
	}

	contentType, err := enums.ConvertStringToMaterialContentType(object.ContentType)
	if contentType == nil {
		exception := exceptions.Material.InvalidType(contentType)
		if err != nil {
			exception.WithError(err)
		}
		return nil, exception
	}

	partialUpdate.Values.ContentType = contentType
	partialUpdate.Values.ParseMediaType = object.ParseMediaType
	partialUpdate.Values.Size = &object.Size

	material, exception := materialRepository.UpdateOneById(
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

func (s *MaterialService) MoveMyMaterialById(
	reqDto *dtos.MoveMyMaterialByIdReqDto,
) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	subShelfRepository := repositories.NewSubShelfRepository(s.db)
	materialRepository := repositories.NewMaterialRepository(s.db)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	if hasPermission := subShelfRepository.HasPermission(
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		types.Ternary_Negative,
	); !hasPermission {
		return nil, exceptions.Shelf.NoPermission()
	}
	material, exception := materialRepository.CheckPermissionAndGetOneById(
		reqDto.Body.MaterialId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		types.Ternary_Negative, // exclude the deleted materials
	)
	if exception != nil {
		return nil, exception
	}

	result := s.db.Exec(`
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
	reqDto *dtos.MoveMyMaterialsByIdsReqDto,
) (*dtos.MoveMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	materials, exception := materialRepository.CheckPermissionsAndGetManyByIds(
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

	result := s.db.Exec(`
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
	reqDto *dtos.RestoreMyMaterialByIdReqDto,
) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	exception := materialRepository.RestoreSoftDeletedOneById(
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
	reqDto *dtos.RestoreMyMaterialsByIdsReqDto,
) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	exception := materialRepository.RestoreSoftDeletedManyByIds(
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
	reqDto *dtos.DeleteMyMaterialByIdReqDto,
) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	exception := materialRepository.SoftDeleteOneById(
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
	reqDto *dtos.DeleteMyMaterialsByIdsReqDto,
) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	exception := materialRepository.SoftDeleteManyByIds(
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
