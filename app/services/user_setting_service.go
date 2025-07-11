package services

import (
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/repositories"
)

/* ============================== Interface & Instance ============================== */

type UserSettingServiceInterface interface {
	GetMySetting(reqDto *dtos.GetMySettingReqDto) (*dtos.GetMySettingResDto, *exceptions.Exception)
	UpdateMySetting(reqDto *dtos.UpdateMySettingReqDto) (*dtos.UpdateMySettingResDto, *exceptions.Exception)
}

type userSettingService struct{}

var UserSettingService UserSettingServiceInterface = &userSettingService{}

/* ============================== Services for UserSetting ============================== */

func (s *userSettingService) GetMySetting(reqDto *dtos.GetMySettingReqDto) (*dtos.GetMySettingResDto, *exceptions.Exception) {
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

func (s *userSettingService) UpdateMySetting(reqDto *dtos.UpdateMySettingReqDto) (*dtos.UpdateMySettingResDto, *exceptions.Exception) {
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
