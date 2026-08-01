package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-infos"
	caches "github.com/HiIamJeff67/notezy-backend/internal/caches"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/caches/inputs"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
)

type UserInfoServiceInterface interface {
	GetMyInfo(ctx context.Context, requestDto *userinfosdto.GetMyInfoRequestDto) (*userinfosdto.GetMyInfoResponseDto, *exceptions.Exception)
	UpdateMyInfo(ctx context.Context, requestDto *userinfosdto.UpdateMyInfoRequestDto) (*userinfosdto.UpdateMyInfoResponseDto, *exceptions.Exception)

	// services for public userInfos
	GetPublicUserInfoByUserPublicId(ctx context.Context, publicId uuid.UUID) (*gqlmodels.PublicUserInfo, *exceptions.Exception)
	GetPublicUserInfosByUserPublicIds(ctx context.Context, publicIds []uuid.UUID) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception)
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
		db = data.NotezyDB
	}
	return &UserInfoService{
		db:                 db,
		userInfoRepository: userInfoRepository,
	}
}

/* ============================== Service Methods for UserInfo ============================== */

func (s *UserInfoService) GetMyInfo(
	ctx context.Context, requestDto *userinfosdto.GetMyInfoRequestDto,
) (*userinfosdto.GetMyInfoResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"UserInfo",
			"GetMyInfo",
			"User info request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	userInfo, exception := s.userInfoRepository.GetOneByUserId(
		actorUserId,
		options.WithDB(s.db.WithContext(ctx)),
	)
	if exception != nil {
		return nil, exception
	}

	var country *string
	if userInfo.Country != nil {
		countryString := userInfo.Country.String()
		country = &countryString
	}
	return &userinfosdto.GetMyInfoResponseDto{
		CoverBackgroundURL: userInfo.CoverBackgroundURL,
		AvatarURL:          userInfo.AvatarURL,
		Header:             userInfo.Header,
		Introduction:       userInfo.Introduction,
		Gender:             userInfo.Gender.String(),
		Country:            country,
		BirthDate:          userInfo.BirthDate,
		UpdatedAt:          userInfo.UpdatedAt,
	}, nil
}

func (s *UserInfoService) UpdateMyInfo(
	ctx context.Context, requestDto *userinfosdto.UpdateMyInfoRequestDto,
) (*userinfosdto.UpdateMyInfoResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"UserInfo",
			"UpdateMyInfo",
			"User info request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	var gender *enums.UserGender
	if requestDto.Body.Values.Gender != nil {
		parsedGender, err := enums.ConvertStringToUserGender(*requestDto.Body.Values.Gender)
		if err != nil {
			return nil, exceptions.InvalidInput("UserInfo").WithOrigin(err)
		}
		gender = parsedGender
	}
	var country *enums.Country
	if requestDto.Body.Values.Country != nil {
		parsedCountry, err := enums.ConvertStringToCountry(*requestDto.Body.Values.Country)
		if err != nil {
			return nil, exceptions.InvalidInput("UserInfo").WithOrigin(err)
		}
		country = parsedCountry
	}

	updatedUserInfo, exception := s.userInfoRepository.UpdateOneByUserId(
		actorUserId,
		inputs.PartialUpdateUserInfoInput{
			Values: inputs.UpdateUserInfoInput{
				CoverBackgroundURL: requestDto.Body.Values.CoverBackgroundURL,
				AvatarURL:          requestDto.Body.Values.AvatarURL,
				Header:             requestDto.Body.Values.Header,
				Introduction:       requestDto.Body.Values.Introduction,
				Gender:             gender,
				Country:            country,
				BirthDate:          requestDto.Body.Values.BirthDate,
			},
			SetNull: requestDto.Body.SetNull,
		},
		options.WithDB(s.db.WithContext(ctx)),
	)
	if exception != nil {
		return nil, exception
	}

	var user schemas.User
	if result := s.db.WithContext(ctx).
		Select("name").
		Where("id = ?", actorUserId).
		First(&user); result.Error != nil {
		return nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveUser",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}

	exception = caches.UserDataStore.Update(user.Name, cacheinputs.UpdateUserDataCacheInput{
		AvatarURL: requestDto.Body.Values.AvatarURL,
	})
	if exception != nil && logs.NotezyLogger != nil {
		logs.NotezyLogger.Error(
			ctx,
			exception.Origin(),
			exception.String(),
		)
	}

	return &userinfosdto.UpdateMyInfoResponseDto{
		UpdatedAt: updatedUserInfo.UpdatedAt,
	}, nil
}

/* ============================== Service Methods for Public UserInfo (Only available in GraphQL) ============================== */

// use the searchable user cursor (we only give the search functionality on users)
func (s *UserInfoService) GetPublicUserInfoByUserPublicId(
	ctx context.Context,
	publicId uuid.UUID,
) (*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	userInfo := schemas.UserInfo{}
	result := db.Model(&schemas.UserInfo{}).
		Joins(`LEFT JOIN "UserTable" u ON u.id = "UserInfoTable".user_id`).
		Where("u.public_id = ?", publicId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.New(
			"NotFound",
			"UserInfo",
			"GetPublicUserInfoByUserPublicId",
			"User info was not found",
			http.StatusNotFound,
		).WithOrigin(err)
	}

	return userInfo.ToPublicUserInfo(), nil
}

func (s *UserInfoService) GetPublicUserInfosByUserPublicIds(
	ctx context.Context, publicIds []uuid.UUID,
) ([]*gqlmodels.PublicUserInfo, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicUserInfo{}, nil
	}

	db := s.db.WithContext(ctx)

	uniquePublicIds := make([]uuid.UUID, 0)
	seen := make(map[uuid.UUID]bool)
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
		UserPublicId uuid.UUID `gorm:"column:user_public_id"`
	}
	result := db.Model(&schemas.UserInfo{}).
		Select(`"UserInfoTable".*, u.public_id AS user_public_id`).
		Joins(`LEFT JOIN "UserTable" u ON u.id = "UserInfoTable".user_id`).
		Where("u.public_id IN ?", uniquePublicIds).
		Find(&userInfosWithPublicUserIds)
	if err := result.Error; err != nil {
		return nil, exceptions.New(
			"QueryFailed",
			"UserInfo",
			"GetPublicUserInfosByUserPublicIds",
			"Failed to retrieve user infos",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	publicIdToIndexesMap := make(map[uuid.UUID][]int)
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
