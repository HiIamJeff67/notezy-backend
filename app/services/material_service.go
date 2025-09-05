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

	material, exception := materialRepository.GetOneById(reqDto.Body.MaterialId, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	downloadURL, exception := s.storage.PresignGetObjectByKey(ctx, material.ContentKey, nil)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyMaterialByIdResDto{
		Id:            material.Id,
		RootShelfId:   material.RootShelfId,
		ParentShelfId: material.ParentShelfId,
		Name:          material.Name,
		Type:          material.Type,
		DownloadURL:   downloadURL,
		ContentType:   material.ContentType,
		DeletedAt:     material.DeletedAt,
		UpdatedAt:     material.UpdatedAt,
		CreatedAt:     material.CreatedAt,
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
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON s.id = uts.shelf_id").
		Where("uts.shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Body.ShelfId,
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
			Id:            material.Id,
			RootShelfId:   material.RootShelfId,
			ParentShelfId: material.ParentShelfId,
			Name:          material.Name,
			Type:          material.Type,
			DownloadURL:   downloadURL,
			ContentType:   material.ContentType,
			DeletedAt:     material.DeletedAt,
			UpdatedAt:     material.UpdatedAt,
			CreatedAt:     material.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *MaterialService) CreateTextbookMaterial(ctx context.Context, reqDto *dtos.CreateMaterialReqDto) (*dtos.CreateMaterialResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	newMaterialId := uuid.New()
	newContentKey := s.storage.GenerateKey(
		reqDto.ContextFields.UserPublicId.String(),
		newMaterialId.String(),
	)
	zeroSize := int64(0)
	_, exception := materialRepository.CreateOne(
		reqDto.ContextFields.UserId,
		inputs.CreateMaterialInput{
			Id:            newMaterialId,
			RootShelfId:   reqDto.Body.RootShelfId,
			ParentShelfId: reqDto.Body.ParentShelfId,
			Name:          reqDto.Body.Name,
			Size:          zeroSize,
			Type:          enums.MaterialType_Textbook,
			ContentKey:    newContentKey,
			ContentType:   enums.MaterialContentType_Markdown,
		},
	)
	if exception != nil {
		return nil, exception
	}

	newContent := bytes.NewReader([]byte{})

	_, exception = s.storage.PutObjectByKey(ctx, newContentKey, newContent, zeroSize)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateMaterialResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *MaterialService) RestoreMyMaterialById(reqDto *dtos.RestoreMyMaterialByIdReqDto) (*dtos.RestoreMyMaterialByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	materialRepository := repositories.NewMaterialRepository(s.db)

	exception := materialRepository.RestoreSoftDeletedOneById(reqDto.Body.MaterialId, reqDto.ContextFields.UserId)
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

	exception := materialRepository.RestoreSoftDeletedManyByIds(reqDto.Body.MaterialIds, reqDto.ContextFields.UserId)
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

	exception := materialRepository.SoftDeleteOneById(reqDto.Body.MaterialId, reqDto.ContextFields.UserId)
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

	exception := materialRepository.SoftDeleteManyByIds(reqDto.Body.MaterialIds, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyMaterialsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
