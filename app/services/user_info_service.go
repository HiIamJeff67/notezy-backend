package services

import (
	"context"

	"gorm.io/gorm"

	"notezy-backend/app/caches"
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
	GetPublicUserInfoByUserPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUserInfo, *exceptions.Exception)
	GetPublicUserInfosByUserPublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception)
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

	exception = caches.UpdateUserDataCache(reqDto.UserId, caches.UpdateUserDataCacheDto{
		AvatarURL: reqDto.Values.AvatarURL,
	})
	if exception != nil {
		exception.Log()
	}

	return &dtos.UpdateMyInfoResDto{
		UpdatedAt: updatedUserInfo.UpdatedAt,
	}, nil
}

/* ============================== Service Methods for Public UserInfo (Only available in GraphQL) ============================== */

// use the searchable user cursor (we only give the search functionality on users)
func (s *UserInfoService) GetPublicUserInfoByUserPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUserInfo, *exceptions.Exception) {
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

func (s *UserInfoService) GetPublicUserInfosByUserPublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicUserInfo{}, nil
	}

	uniquePublicIds := make([]string, 0)
	seen := make(map[string]bool)
	for _, publicId := range publicIds {
		if !seen[publicId] {
			uniquePublicIds = append(uniquePublicIds, publicId)
			seen[publicId] = true
		}
	}
	if len(uniquePublicIds) == 0 {
		return make([]*gqlmodels.PublicUserInfo, len(publicIds)), nil
	}

	var userInfosWithPublicUserIds []*struct {
		schemas.UserInfo
		UserPublicId string `gorm:"column:user_public_id"`
	}
	result := s.db.Table(schemas.UserInfo{}.TableName()+" ui").
		Select("ui.*, u.public_id as user_public_id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = ui.user_id").
		Where("u.public_id IN ?", uniquePublicIds).
		Find(&userInfosWithPublicUserIds)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	publicIdToIndexesMap := make(map[string][]int)
	for index, publidId := range publicIds {
		publicIdToIndexesMap[publidId] = append(publicIdToIndexesMap[publidId], index)
	}

	publicUserInfos := make([]*gqlmodels.PublicUserInfo, len(publicIds))
	for _, userInfoWithPublicUserId := range userInfosWithPublicUserIds {
		for _, index := range publicIdToIndexesMap[userInfoWithPublicUserId.UserPublicId] {
			publicUserInfos[index] = userInfoWithPublicUserId.UserInfo.ToPublicUserInfo()
		}
	}

	return publicUserInfos, nil
}
