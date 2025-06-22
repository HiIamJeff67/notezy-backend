package services

import (
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
)

/* ============================== Interface & Instance ============================== */

type UserInfoServiceInterface interface {
	GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception)
	UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception)
}

type userInfoService struct{}

var UserInfoService UserInfoServiceInterface = &userInfoService{}

/* ============================== Services ============================== */

func (s *userInfoService) GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception) {
	userInfoRepository := repositories.NewUserInfoRepository(nil)

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

func (s *userInfoService) UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception) {
	userInfoRepository := repositories.NewUserInfoRepository(nil)

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
