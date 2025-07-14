package services

import (
	"context"
	"strings"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Interface & Instance ============================== */

type UserInfoServiceInterface interface {
	GetMyInfo(reqDto *dtos.GetMyInfoReqDto) (*dtos.GetMyInfoResDto, *exceptions.Exception)
	UpdateMyInfo(reqDto *dtos.UpdateMyInfoReqDto) (*dtos.UpdateMyInfoResDto, *exceptions.Exception)

	// services for public userInfos
	GetPublicUserInfoByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicUserInfo, *exceptions.Exception)
	GetPublicUserInfosByEncodedSearchCursor(ctx context.Context, encodedSearchCursors []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception)
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

/* ============================== Services for UserInfo ============================== */

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

/* ============================== Services for Public UserInfo (Only available in GraphQL) ============================== */

// use the searchable user cursor (we only give the search functionality on users)
func (s *UserInfoService) GetPublicUserInfoByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchUserCursorFields](encodedSearchCursor)
	if exception != nil {
		return nil, exception
	}

	userInfo := schemas.UserInfo{}
	result := s.db.Table(schemas.UserInfo{}.TableName()).
		Joins("LEFT JOIN \"UserTable\" u ON u.id = user_id").
		Where("u.search_cursor_id = ?", searchCursor.Fields.SearchCursorID).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return userInfo.ToPublicUserInfo(), nil
}

func (s *UserInfoService) GetPublicUserInfosByEncodedSearchCursor(ctx context.Context, encodedSearchCursors []string) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	if len(encodedSearchCursors) == 0 {
		return []*gqlmodels.PublicUserInfo{}, nil
	}

	var invalidIndiceCount int = 0
	var searchCursorIds []*string
	for _, encodedSearchCursor := range encodedSearchCursors {
		searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchUserCursorFields](encodedSearchCursor)
		if exception != nil || len(strings.ReplaceAll(searchCursor.Fields.SearchCursorID, " ", "")) == 0 {
			invalidIndiceCount++
		} else {
			searchCursorIds = append(searchCursorIds, &searchCursor.Fields.SearchCursorID)
		}
	}

	if invalidIndiceCount == len(encodedSearchCursors) {
		return make([]*gqlmodels.PublicUserInfo, len(encodedSearchCursors)), nil
	}

	var userInfos []struct {
		schemas.UserInfo
		SearchCursorId string `gorm:"column:search_cursor_id"`
	}

	result := s.db.Table(schemas.UserInfo{}.TableName()+" ui").
		Select("ui.*, u.search_cursor_id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = ui.user_id").
		Where("u.search_cursor_id IN ?", searchCursorIds).
		Find(&userInfos)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	publicUserInfos := make([]*gqlmodels.PublicUserInfo, len(userInfos))
	for index, userInfo := range userInfos {
		publicUserInfos[index] = userInfo.UserInfo.ToPublicUserInfo()
	}

	return publicUserInfos, nil
}
