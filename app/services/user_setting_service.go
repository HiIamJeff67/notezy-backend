package services

import (
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"

	"gorm.io/gorm"
)

/* ============================== Interface & Instance ============================== */

type UserSettingServiceInterface interface {
	GetMySetting(reqDto *dtos.GetMySettingReqDto) (*dtos.GetMySettingResDto, *exceptions.Exception)
	UpdateMySetting(reqDto *dtos.UpdateMySettingReqDto) (*dtos.UpdateMySettingResDto, *exceptions.Exception)
}

type UserSettingService struct {
	db *gorm.DB
}

func NewUserSettingService(db *gorm.DB) UserSettingServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserSettingService{db: db}
}

/* ============================== Services for UserSetting ============================== */

func (s *UserSettingService) GetMySetting(reqDto *dtos.GetMySettingReqDto) (*dtos.GetMySettingResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userSettingRepository := repositories.NewUserSettingRepository(nil)

	userSetting, exception := userSettingRepository.GetOneByUserId(reqDto.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMySettingResDto{
		Language:           userSetting.Language,
		GeneralSettingCode: userSetting.GeneralSettingCode,
		PrivacySettingCode: userSetting.PrivacySettingCode,
	}, nil
}

func (s *UserSettingService) UpdateMySetting(reqDto *dtos.UpdateMySettingReqDto) (*dtos.UpdateMySettingResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userSettingRepository := repositories.NewUserSettingRepository(nil)

	updatedUserSetting, exception := userSettingRepository.UpdateOneByUserId(reqDto.UserId, inputs.PartialUpdateUserSettingInput{
		Values: inputs.UpdateUserSettingInput{
			Language:           &reqDto.Values.Language,
			GeneralSettingCode: &reqDto.Values.GeneralSettingCode,
			PrivacySettingCode: &reqDto.Values.PrivacySettingCode,
		},
		SetNull: reqDto.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMySettingResDto{
		UpdatedAt: updatedUserSetting.UpdatedAt,
	}, nil
}
