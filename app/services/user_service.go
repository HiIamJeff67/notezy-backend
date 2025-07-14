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
	UpdateMe(reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception)

	// services for public users
	GetPublicUserByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicUser, *exceptions.Exception)
	SearchPublicUsers(ctx context.Context, gqlInput gqlmodels.SearchUserInput) (*gqlmodels.SearchUserConnection, *exceptions.Exception)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserService{db: db}
}

/* ============================== Services for Users ============================== */

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

/* ============================== Services for Public User (Only available in GraphQL) ============================== */

func (s *UserService) GetPublicUserByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicUser, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(s.db)

	searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchUserCursorFields](encodedSearchCursor)
	if exception != nil {
		return nil, exception
	}

	publicUser, exception := userRepository.GetPublicOneBySearchCursorId(searchCursor.Fields.SearchCursorID)
	if exception != nil {
		// try to get the user from email
		// user, exception := userRepository.GetOneByEmail(searchCursor.Fields.Email)
		// if exception != nil {
		// 	return nil, exception
		// }
		// return user.ToPublicUser(), nil
		return nil, exception
	}

	return publicUser, nil
}

func (s *UserService) SearchPublicUsers(ctx context.Context, gqlInput gqlmodels.SearchUserInput) (*gqlmodels.SearchUserConnection, *exceptions.Exception) {
	startTime := time.Now()

	query := s.db.WithContext(ctx).Model(&schemas.User{})

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ? OR display_name ILIKE ? OR email ILIKE ?",
			"%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%", "%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchUserCursorFields](*gqlInput.After)
		if exception != nil {
			return nil, exception
		}

		query.Where("search_cursor_id > ?", searchCursor.Fields.SearchCursorID)
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
	searchEdges := make([]*gqlmodels.SearchUserEdge, len(users))

	for index, user := range users {
		searchCursor := util.SearchCursor[gqlmodels.SearchUserCursorFields]{
			Fields: gqlmodels.SearchUserCursorFields{
				SearchCursorID: user.SearchCursorId,
			},
		}
		encodedSearchCursor, exception := searchCursor.EncodeSearchCursor()
		if exception != nil {
			return nil, exception
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
