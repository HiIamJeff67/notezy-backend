package services

import (
	"context"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
)

/* ============================== Interface & Instance ============================== */

type UserInfoServiceInterface interface {
	GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception)
	UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception)

	// services for public userInfos
	GetPublicUserInfoByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUserInfo, *exceptions.Exception)
	GetPublicUserInfosByPublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception)
}

type UserInfoService struct {
	db *gorm.DB
}

func NewUserInfoService(db *gorm.DB) UserInfoServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserInfoService{db: db}
}

/* ============================== Service Methods for UserInfo ============================== */

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
		Country:            *userInfo.Country,
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

/* ============================== Service Methods for Public UserInfo (Only available in GraphQL) ============================== */

// use the searchable user cursor (we only give the search functionality on users)
func (s *UserInfoService) GetPublicUserInfoByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	userInfo := schemas.UserInfo{}
	result := s.db.Table(schemas.UserInfo{}.TableName()).
		Joins("LEFT JOIN \"UserTable\" u ON u.id = user_id").
		Where("u.public_id = ?", publicId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return userInfo.ToPublicUserInfo(), nil
}

func (s *UserInfoService) GetPublicUserInfosByPublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicUserInfo{}, nil
	}

	var userInfos []schemas.UserInfo
	result := s.db.Table(schemas.UserInfo{}.TableName()+" ui").
		Select("ui.*, u.public_id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = ui.user_id").
		Where("u.public_id IN ?", publicIds).
		Find(&userInfos)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	publicUserInfos := make([]*gqlmodels.PublicUserInfo, len(userInfos))
	for index, userInfo := range userInfos {
		publicUserInfos[index] = userInfo.ToPublicUserInfo()
	}

	return publicUserInfos, nil
}
