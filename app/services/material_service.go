package services

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
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
		reqDto.Param.RootShelfId,
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
		Joins("LEFT JOIN \"ShelfTable\" s ON \"MaterialTable\".root_shelf_id = s.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.root_shelf_id").
		Where("uts.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Body.RootShelfId,
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
		reqDto.ContextFields.UserId,
		reqDto.Body.RootShelfId,
		inputs.CreateMaterialInput{
			Id:            newMaterialId,
			ParentShelfId: reqDto.Body.ParentShelfId,
			Name:          reqDto.Body.Name,
			Size:          zeroSize,
			Type:          enums.MaterialType_Textbook,
			ContentKey:    newContentKey,
			ContentType:   enums.MaterialContentType_PlainText,
		},
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newContentFile := bytes.NewReader([]byte{})

	object, exception := s.storage.GetObject(newContentKey, newContentFile, zeroSize)
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

	hasContentFile := reqDto.Body.ContentFile != nil

	// check if the material content type is allowed in the material type of textbook
	materialType := enums.MaterialType_Textbook
	var contentType *enums.MaterialContentType = nil // initialize the content type to be nil
	var object *storages.Object = nil                // initialize the object first to be nil
	var contentKey = s.storage.GetKey(reqDto.ContextFields.UserPublicId.String(), reqDto.Body.MaterialId.String())
	if hasContentFile {
		var exception *exceptions.Exception = nil // initialize the exception while checking the file
		bufReader := bufio.NewReader(reqDto.Body.ContentFile)
		peekBytes, err := bufReader.Peek(512) // try to peek the first 512 bytes
		if err != nil && err != io.EOF {      // if err == io.EOF, then the total number of bytes is not greater than 512
			tx.Rollback()
			return nil, exceptions.Material.CannotPeekFiles()
		}
		actualContentType := http.DetectContentType(peekBytes)
		if !materialType.IsContentTypeStringAllowed(actualContentType) {
			tx.Rollback()
			return nil, exceptions.Material.MaterialContentTypeNotAllowedInMaterialType(
				reqDto.Body.MaterialId.String(),
				materialType.String(),
				actualContentType,
				materialType.AllowedContentTypeStrings(),
			)
		}
		contentType, err := enums.ConvertStringToMaterialContentType(actualContentType)
		if contentType == nil {
			exception := exceptions.Material.InvalidType(contentType)
			if err != nil {
				exception.WithError(err)
			}
			exception.Log()
		}

		var fileHeaderSize int64 = 0
		if reqDto.Body.Size != nil {
			fileHeaderSize = *reqDto.Body.Size
		}

		object, exception = s.storage.GetObject(contentKey, reqDto.Body.ContentFile, fileHeaderSize)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	partialUpdate := inputs.PartialUpdateMaterialInput{
		Values: inputs.UpdateMaterialInput{
			Name: reqDto.Body.PartialUpdate.Values.Name,
			// content key remain the same here
			ContentType: contentType, // since content type is just a pointer, we can directly place it here
		},
		SetNull: reqDto.Body.PartialUpdate.SetNull,
	}
	if object != nil {
		partialUpdate.Values.Size = &object.Size
	}
	material, exception := materialRepository.UpdateOneById(
		reqDto.Body.MaterialId,
		reqDto.Body.RootShelfId,
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
		reqDto.Body.SourceRootShelfId,
		reqDto.ContextFields.UserId,
		nil,
		inputs.PartialUpdateMaterialInput{
			Values: inputs.UpdateMaterialInput{
				RootShelfId:   &reqDto.Body.DestinationRootShelfId,
				ParentShelfId: &reqDto.Body.DestinationParentShelfId,
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
		reqDto.Body.RootShelfId,
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
		reqDto.Body.RootShelfId,
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
		reqDto.Body.RootShelfId,
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
		reqDto.Body.RootShelfId,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyMaterialsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
