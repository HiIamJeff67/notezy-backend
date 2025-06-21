package services

import (
	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
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

type userService struct{}

var UserService UserServiceInterface = &userService{}

/* ============================== Services ============================== */

func (u *userService) GetMe(reqDto *dtos.GetMeReqDto) (*dtos.GetMeResDto, *exceptions.Exception) {
	userDataCache, exception := caches.GetUserDataCache(reqDto.UserId)
	if exception != nil {
		return nil, exception
	}

	return userDataCache, nil
}

// for temporary use
func (u *userService) GetAllUsers() (*[]schemas.User, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(nil)

	users, exception := userRepository.GetAll()
	if exception != nil {
		return nil, exception
	}

	return users, nil
}

func (u *userService) UpdateMe(reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(nil)

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
