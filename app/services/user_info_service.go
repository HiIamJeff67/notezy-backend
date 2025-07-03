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

type UserInfoServiceInterface interface {
	GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception)
	UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception)
}

type UserInfoService struct {
	db *gorm.DB
}

func NewUserInfoService(db *gorm.DB) UserInfoServiceInterface {
	return &UserInfoService{
		db: db,
	}
}

/* ============================== Services ============================== */

func (s *UserInfoService) GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userInfoRepository := repositories.NewUserInfoRepository(s.db)

	userInfo, exception := userInfoRepository.GetOneByUserId(reqDto.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyInfoResDto{
		CoverBackgroundURL: userInfo.CoverBackgroundURL,
		AvatarURL:          userInfo.AvatarURL,
		Header:             userInfo.Header,
		Introduction:       userInfo.Introduction,
		Gender:             userInfo.Gender,
		Country:            userInfo.Country,
		BirthDate:          userInfo.BirthDate,
		UpdatedAt:          userInfo.UpdatedAt,
	}, nil
}

func (s *UserInfoService) UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userInfoRepository := repositories.NewUserInfoRepository(s.db)

	updatedUserInfo, exception := userInfoRepository.UpdateOneByUserId(reqDto.UserId, inputs.PartialUpdateUserInfoInput{
		Values: inputs.UpdateUserInfoInput{
			CoverBackgroundURL: reqDto.Values.CoverBackgroundURL,
			AvatarURL:          reqDto.Values.AvatarURL,
			Header:             reqDto.Values.Header,
			Introduction:       reqDto.Values.Introduction,
			Gender:             reqDto.Values.Gender,
			Country:            reqDto.Values.Country,
			BirthDate:          reqDto.Values.BirthDate,
		},
		SetNull: reqDto.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyInfoResDto{
		UpdatedAt: updatedUserInfo.UpdatedAt,
	}, nil
}
