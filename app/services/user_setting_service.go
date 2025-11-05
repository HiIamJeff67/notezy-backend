package services

import (
	"context"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	validation "notezy-backend/app/validation"
)

/* ============================== Interface & Instance ============================== */

type UserSettingServiceInterface interface {
	GetMySetting(ctx context.Context, reqDto *dtos.GetMySettingReqDto) (*dtos.GetMySettingResDto, *exceptions.Exception)
	UpdateMySetting(ctx context.Context, reqDto *dtos.UpdateMySettingReqDto) (*dtos.UpdateMySettingResDto, *exceptions.Exception)
}

type UserSettingService struct {
	db                    *gorm.DB
	userSettingRepository repositories.UserSettingRepositoryInterface
}

func NewUserSettingService(
	db *gorm.DB,
	userSettingRepository repositories.UserSettingRepositoryInterface,
) UserSettingServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserSettingService{
		db:                    db,
		userSettingRepository: userSettingRepository,
	}
}

/* ============================== Service Methods for UserSetting ============================== */

func (s *UserSettingService) GetMySetting(
	ctx context.Context, reqDto *dtos.GetMySettingReqDto,
) (*dtos.GetMySettingResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	userSetting, exception := s.userSettingRepository.GetOneByUserId(
		db,
		reqDto.ContextFields.UserId,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMySettingResDto{
		Language:           userSetting.Language,
		GeneralSettingCode: userSetting.GeneralSettingCode,
		PrivacySettingCode: userSetting.PrivacySettingCode,
	}, nil
}

func (s *UserSettingService) UpdateMySetting(
	ctx context.Context, reqDto *dtos.UpdateMySettingReqDto,
) (*dtos.UpdateMySettingResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	updatedUserSetting, exception := s.userSettingRepository.UpdateOneByUserId(
		db,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateUserSettingInput{
			Values: inputs.UpdateUserSettingInput{
				Language:           &reqDto.Body.Values.Language,
				GeneralSettingCode: &reqDto.Body.Values.GeneralSettingCode,
				PrivacySettingCode: &reqDto.Body.Values.PrivacySettingCode,
			},
			SetNull: reqDto.Body.SetNull,
		})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMySettingResDto{
		UpdatedAt: updatedUserSetting.UpdatedAt,
	}, nil
}
