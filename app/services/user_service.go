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
	util "notezy-backend/app/util"
)

/* ============================== Interface & Instance ============================== */

type UserServiceInterface interface {
	GetMe(reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception)
	GetAllUsers() (*[]schemas.User, *exceptions.Exception)
	SearchUsers(ctx context.Context, gqlInput gqlmodels.SearchableUserInput) (*gqlmodels.SearchableUserConnection, *exceptions.Exception)
	UpdateMe(reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserServiceInterface {
	return &UserService{
		db: db,
	}
}

/* ============================== Auxilary Functions for GraphQL ============================== */

func (s *UserService) convertUserToPublicUser(user *schemas.User) *gqlmodels.PublicUser {
	return &gqlmodels.PublicUser{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
		Plan:        user.Plan,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		UserInfo:    &gqlmodels.PublicUserInfo{},
		Badges:      []*gqlmodels.PublicBadge{},
		Themes:      []*gqlmodels.PublicTheme{},
	}
}

func (s *UserService) convertPublicUserToUser() {}

/* ============================== Services ============================== */

func (s *UserService) GetMe(reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(reqDto.UserId)
	if exception != nil {
		return nil, exception
	}

	return userDataCache, nil
}

// for temporary use
func (s *UserService) GetAllUsers() (*[]schemas.User, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(nil)

	users, exception := userRepository.GetAll()
	if exception != nil {
		return nil, exception
	}

	return users, nil
}

func (s *UserService) SearchUsers(ctx context.Context, gqlInput gqlmodels.SearchableUserInput) (*gqlmodels.SearchableUserConnection, *exceptions.Exception) {
	startTime := time.Now()

	query := s.db.WithContext(ctx).Model(&schemas.User{})

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ? OR display_name ILIKE ? OR email ILIKE ?",
			"%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, exception := util.DecodeSearchCursor(*gqlInput.After)
		if exception != nil {
			return nil, exception
		}

		for _, field := range searchCursor.Fields {
			query.Where(field.Name+" = ?", field.Value)
		}
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := "ASC"
		if *gqlInput.SortOrder == gqlmodels.SearchableSortOrderDesc {
			cending = "DESC"
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchableUserSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchableUserSortByLastActive:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchableUserSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		}
	}

	limit := 10
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	query = query.Limit(limit + 1)

	var users []schemas.User
	if err := query.Find(&users).Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	hasNextPage := len(users) > limit // since we fetch an additional one

	searchEdges := make([]*gqlmodels.SearchableUserEdge, len(users))
	for index, user := range users {
		searchCursorFields := map[string]interface{}{
			"name":         user.Name,
			"display_name": user.DisplayName,
			"email":        user.Email,
		}
		searchCursor, err := util.EncodeSearchCursor(searchCursorFields)
		if err != nil {
			return nil, err
		}
		if searchCursor == nil {
			return nil, exceptions.Searchable.FailedToUnMarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchableUserEdge{
			SearchCursor: *searchCursor,
			Node:         s.convertUserToPublicUser(&user),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) != 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartSearchCursor = &searchEdges[0].SearchCursor
		searchPageInfo.EndSearchCursor = &searchEdges[len(searchEdges)-1].SearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchableUserConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}

func (s *UserService) UpdateMe(reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userRepository := repositories.NewUserRepository(s.db)

	updatedUser, exception := userRepository.UpdateOneById(reqDto.UserId, inputs.PartialUpdateUserInput{
		Values: inputs.UpdateUserInput{
			DisplayName: reqDto.Values.DisplayName,
			Status:      reqDto.Values.Status,
		},
		SetNull: reqDto.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMeResDto{UpdatedAt: updatedUser.UpdatedAt}, nil
}

// may add some business logic of payment
// func UpdatePlan(reqDto *dtos.UpdatePlanReqDto) (*dtos.UpdatePlanResDto, *exceptions.Exception) {

// }
