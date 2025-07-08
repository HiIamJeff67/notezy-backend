package services

import (
	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
)

/* ============================== Interface & Instance ============================== */

type UserServiceInterface interface {
	GetMe(reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception)
	GetAllUsers() (*[]schemas.User, *exceptions.Exception)
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

/* ============================== Services ============================== */

func (u *UserService) GetMe(reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception) {
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

// func (s *UserService) SearchUsers(ctx context.Context, gqlInput gqlmodels.SearchableUserInput) (*gqlmodels.SearchableUserConnection, error) {
// 	startTime := time.Now()

// 	query := s.db.WithContext(ctx).Model(&schemas.User{})

// 	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
// 		searchCursor, err := util.DecodeSearchCursor(*gqlInput.After)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid cursor: %w", err)
// 		}

// 	}
// }

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
