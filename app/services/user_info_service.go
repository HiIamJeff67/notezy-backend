package services

import (
	"context"

	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	validation "notezy-backend/app/validation"
)

/* ============================== Interface & Instance ============================== */

type UserInfoServiceInterface interface {
	GetMyInfo(ctx context.Context, reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception)
	UpdateMyInfo(ctx context.Context, reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception)

	// services for public userInfos
	GetPublicUserInfoByUserPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUserInfo, *exceptions.Exception)
	GetPublicUserInfosByUserPublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception)
}

type UserInfoService struct {
	db                 *gorm.DB
	userInfoRepository repositories.UserInfoRepositoryInterface
}

func NewUserInfoService(
	db *gorm.DB,
	userInfoRepository repositories.UserInfoRepositoryInterface,
) UserInfoServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserInfoService{
		db:                 db,
		userInfoRepository: userInfoRepository,
	}
}

/* ============================== Service Methods for UserInfo ============================== */

func (s *UserInfoService) GetMyInfo(
	ctx context.Context, reqDto *dtos.GetMyInfoReqDto,
) (*dtos.GetMyInfoResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	userInfo, exception := s.userInfoRepository.GetOneByUserId(
		db,
		reqDto.ContextFields.UserId,
	)
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

func (s *UserInfoService) UpdateMyInfo(
	ctx context.Context, reqDto *dtos.UpdateMyInfoReqDto,
) (*dtos.UpdateMyInfoResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	updatedUserInfo, exception := s.userInfoRepository.UpdateOneByUserId(
		db,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateUserInfoInput{
			Values: inputs.UpdateUserInfoInput{
				CoverBackgroundURL: reqDto.Body.Values.CoverBackgroundURL,
				AvatarURL:          reqDto.Body.Values.AvatarURL,
				Header:             reqDto.Body.Values.Header,
				Introduction:       reqDto.Body.Values.Introduction,
				Gender:             reqDto.Body.Values.Gender,
				Country:            reqDto.Body.Values.Country,
				BirthDate:          reqDto.Body.Values.BirthDate,
			},
			SetNull: reqDto.Body.SetNull,
		})
	if exception != nil {
		return nil, exception
	}

	exception = caches.UpdateUserDataCache(reqDto.ContextFields.UserId, caches.UpdateUserDataCacheDto{
		AvatarURL: reqDto.Body.Values.AvatarURL,
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
func (s *UserInfoService) GetPublicUserInfoByUserPublicId(
	ctx context.Context,
	publicId string,
) (*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	userInfo := schemas.UserInfo{}
	result := db.Table(schemas.UserInfo{}.TableName()).
		Joins("LEFT JOIN \"UserTable\" u ON u.id = user_id").
		Where("u.public_id = ?", publicId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return userInfo.ToPublicUserInfo(), nil
}

func (s *UserInfoService) GetPublicUserInfosByUserPublicIds(
	ctx context.Context, publicIds []string,
) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicUserInfo{}, nil
	}

	db := s.db.WithContext(ctx)

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
	result := db.Table(schemas.UserInfo{}.TableName()+" ui").
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
