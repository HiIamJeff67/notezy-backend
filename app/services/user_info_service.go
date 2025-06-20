package services

import (
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/models/repositories"
)

func GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception) {
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
