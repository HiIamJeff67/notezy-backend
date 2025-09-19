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
	storages "notezy-backend/app/storages"
	validation "notezy-backend/app/validation"
)

/* ============================== Interface & Instance ============================== */

type MaterialServiceInterface interface {
	GetMyMaterialById(ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception)
	SearchMyMaterialsByShelfId(ctx context.Context, reqDto *dtos.SearchMyMaterialsByShelfIdReqDto) (*dtos.SearchMyMaterialsByShelfIdResDto, *exceptions.Exception)
	CreateTextbookMaterial(ctx context.Context, reqDto *dtos.CreateMaterialReqDto) (*dtos.CreateMaterialResDto, *exceptions.Exception)
	SaveMyTextbookMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception)
	MoveMyMaterialById(ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception)
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

func (s *MaterialService) GetMyMaterialById(ctx context.Context, reqDto *dtos.GetMyMaterialByIdReqDto) (*dtos.GetMyMaterialByIdResDto, *exceptions.Exception) {
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

func (s *MaterialService) SearchMyMaterialsByShelfId(ctx context.Context, reqDto *dtos.SearchMyMaterialsByShelfIdReqDto) (*dtos.SearchMyMaterialsByShelfIdResDto, *exceptions.Exception) {
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
			reqDto.Body.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		)
	if len(strings.ReplaceAll(reqDto.Param.Query, " ", "")) > 0 {
		query = query.Where("\"MaterialTable\".name ILIKE ?", "%"+reqDto.Param.Query+"%")
	}

	result := query.Order("updated_at DESC").
		Limit(int(reqDto.Param.Limit)).
		Offset(int(reqDto.Param.Offset)).
		Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	resDto := dtos.SearchMyMaterialsByShelfIdResDto{}
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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Material.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.CreateMaterialResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) SaveMyTextbookMaterialById(ctx context.Context, reqDto *dtos.SaveMyMaterialByIdReqDto) (*dtos.SaveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	tx := s.db.Begin()
	materialRepository := repositories.NewMaterialRepository(tx)

	partialUpdate := inputs.PartialUpdateMaterialInput{
		Values: inputs.UpdateMaterialInput{
			Name: reqDto.Body.Name,
			// content key remain the same here
		},
		SetNull: nil,
	}
	var object *storages.Object = nil // initialize the object first to be nil
	var contentKey = s.storage.GetKey(reqDto.ContextFields.UserPublicId.String(), reqDto.Body.MaterialId.String())
	// check if the material content type is allowed in the material type of textbook
	materialType := enums.MaterialType_Textbook
	// check if there does exist a file in the reqDto
	hasContentFile := reqDto.Body.ContentFile != nil

	// if there does exist a file, extract the data in it and get its content type, parse media type, and actual size, etc.
	if hasContentFile {
		var exception *exceptions.Exception = nil // initialize the exception while checking the file

		var fileHeaderSize int64 = 0
		if reqDto.Body.Size != nil {
			fileHeaderSize = *reqDto.Body.Size
		}

		object, exception = s.storage.NewObject(contentKey, reqDto.Body.ContentFile, fileHeaderSize)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

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
			exception.Log()
		} else {
			partialUpdate.Values.ContentType = contentType
			partialUpdate.Values.ParseMediaType = object.ParseMediaType
		}
	}

	if object != nil {
		partialUpdate.Values.Size = &object.Size
	}
	material, exception := materialRepository.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		&materialType,
		partialUpdate,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// if there does exist a file, then put the file at the end to ensure the entire operation is consistent
	if hasContentFile {
		exception = s.storage.PutObjectByKey(ctx, material.ContentKey, object)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Material.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.SaveMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) MoveMyMaterialById(ctx context.Context, reqDto *dtos.MoveMyMaterialByIdReqDto) (*dtos.MoveMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	material, exception := materialRepository.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.Body.SourceParentSubShelfId,
		reqDto.ContextFields.UserId,
		nil,
		inputs.PartialUpdateMaterialInput{
			Values: inputs.UpdateMaterialInput{
				ParentSubShelfId: &reqDto.Body.DestinationParentSubShelfId,
			},
			SetNull: nil,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyMaterialByIdResDto{
		UpdatedAt: material.UpdatedAt,
	}, nil
}

func (s *MaterialService) RestoreMyMaterialById(reqDto *dtos.RestoreMyMaterialByIdReqDto) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception) {
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

func (s *MaterialService) RestoreMyMaterialsByIds(reqDto *dtos.RestoreMyMaterialsByIdsReqDto) (*dtos.RestoreMyMaterialsByIdsResDto, *exceptions.Exception) {
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

func (s *MaterialService) DeleteMyMaterialById(reqDto *dtos.DeleteMyMaterialByIdReqDto) (*dtos.DeleteMyMaterialByIdResDto, *exceptions.Exception) {
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

func (s *MaterialService) DeleteMyMaterialsByIds(reqDto *dtos.DeleteMyMaterialsByIdsReqDto) (*dtos.DeleteMyMaterialsByIdsResDto, *exceptions.Exception) {
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
