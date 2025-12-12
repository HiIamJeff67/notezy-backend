package services

import (
	"context"
	"strings"
	"time"

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
	constants "notezy-backend/shared/constants"
	searchcursor "notezy-backend/shared/lib/searchcursor"
)

/* ============================== Interface & Instance ============================== */

type UserServiceInterface interface {
	GetUserData(ctx context.Context, reqDto *dtos.GetUserDataReqDto) (*dtos.GetUserDataResDto, *exceptions.Exception)
	GetMe(ctx context.Context, reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception)
	UpdateMe(ctx context.Context, reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception)

	// services for public users
	GetPublicUserByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicUser, *exceptions.Exception)
	GetPublicAuthorByThemePublicIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicUser, *exceptions.Exception)
	SearchPublicUsers(ctx context.Context, gqlInput gqlmodels.SearchUserInput) (*gqlmodels.SearchUserConnection, *exceptions.Exception)
}

type UserService struct {
	db             *gorm.DB
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(
	db *gorm.DB,
	userRepository repositories.UserRepositoryInterface,
) UserServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserService{
		db:             db,
		userRepository: userRepository,
	}
}

/* ============================== Service Methods for Users ============================== */

func (s *UserService) GetUserData(
	ctx context.Context, reqDto *dtos.GetUserDataReqDto,
) (*dtos.GetUserDataResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return userDataCache, nil
}

func (s *UserService) GetMe(
	ctx context.Context, reqDto *dtos.GetMeReqDto,
) (*dtos.GetMeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	user, exception := s.userRepository.GetOneById(
		db,
		reqDto.ContextFields.UserId,
		nil,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMeResDto{
		PublicId:    user.PublicId,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
		Plan:        user.Plan,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserService) UpdateMe(
	ctx context.Context, reqDto *dtos.UpdateMeReqDto,
) (*dtos.UpdateMeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	updatedUser, exception := s.userRepository.UpdateOneById(
		db,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				DisplayName: reqDto.Body.Values.DisplayName,
				Status:      reqDto.Body.Values.Status,
			},
			SetNull: reqDto.Body.SetNull,
		})
	if exception != nil {
		return nil, exception
	}

	if reqDto.Body.Values.DisplayName != nil {
		exception = caches.UpdateUserDataCache(reqDto.ContextFields.UserId, caches.UpdateUserDataCacheDto{
			DisplayName: reqDto.Body.Values.DisplayName,
		})
		if exception != nil {
			exception.Log()
		}
	}
	if reqDto.Body.Values.Status != nil {
		exception = caches.UpdateUserDataCache(reqDto.ContextFields.UserId, caches.UpdateUserDataCacheDto{
			Status: reqDto.Body.Values.Status,
		})
		if exception != nil {
			exception.Log()
		}
	}

	return &dtos.UpdateMeResDto{UpdatedAt: updatedUser.UpdatedAt}, nil
}

// may add some business logic of payment
// func UpdatePlan(reqDto *dtos.UpdatePlanReqDto) (*dtos.UpdatePlanResDto, *exceptions.Exception) {

// }

/* ============================== Service Methods for Public User (Only available in GraphQL) ============================== */

func (s *UserService) GetPublicUserByPublicId(
	ctx context.Context, publicId string,
) (*gqlmodels.PublicUser, *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	user := schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Where("public_id = ?", publicId).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return user.ToPublicUser(), nil
}

func (s *UserService) GetPublicAuthorByThemePublicIds(
	ctx context.Context, publicIds []string,
) ([]*gqlmodels.PublicUser, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicUser{}, nil
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
		return make([]*gqlmodels.PublicUser, len(publicIds)), nil
	}

	var authorsWithPublicThemeIds []*struct {
		schemas.User
		ThemePublicId string `gorm:"theme_public_id"`
	}
	result := db.Table(schemas.User{}.TableName()+" u").
		Select("u.*, t.public_id as theme_public_id").
		Joins("LEFT JOIN \"ThemeTable\" t ON t.author_id = u.id").
		Where("t.public_id IN ?", uniquePublicIds).
		Find(&authorsWithPublicThemeIds)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	publicIdToIndexesMap := make(map[string][]int)
	for index, publicId := range publicIds {
		publicIdToIndexesMap[publicId] = append(publicIdToIndexesMap[publicId], index)
	}

	publicUsers := make([]*gqlmodels.PublicUser, len(publicIds))
	for _, authorWithPublicThemeId := range authorsWithPublicThemeIds {
		for _, index := range publicIdToIndexesMap[authorWithPublicThemeId.ThemePublicId] {
			publicUsers[index] = authorWithPublicThemeId.ToPublicUser()
		}
	}

	return publicUsers, nil
}

func (s *UserService) SearchPublicUsers(
	ctx context.Context, gqlInput gqlmodels.SearchUserInput,
) (*gqlmodels.SearchUserConnection, *exceptions.Exception) {
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.User{})

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ? OR display_name ILIKE ? OR email ILIKE ?",
			"%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchUserCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithError(err)
		}

		query.Where("public_id > ?", searchCursor.Fields.PublicID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchUserSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchUserSortByLastActive:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchUserSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = max(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var users []schemas.User
	if err := query.Find(&users).Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	hasNextPage := len(users) > limit // since we fetch an additional one
	searchEdges := make([]*gqlmodels.SearchUserEdge, len(users))

	for index, user := range users {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchUserCursorFields]{
			Fields: gqlmodels.SearchUserCursorFields{
				PublicID: user.PublicId,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithError(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchUserEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                user.ToPublicUser(),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchUserConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
