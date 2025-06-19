package services

import (
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

func FindMe(reqDto *dtos.FindMeReqDto) (*dtos.FindMeResDto, *exceptions.Exception) {
	userDataCache, exception := caches.GetUserDataCache(reqDto.Id)
	if exception != nil {
		return nil, exception
	}

	return userDataCache, nil
}

// for temporary use
func FindAllUsers() (*[]schemas.User, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(nil)
	users, exception := userRepository.GetAll()
	if exception != nil {
		return nil, exception
	}

	return users, nil
}

func UpdateMe(reqDto *dtos.UpdateMeReqDto) (*dtos.UpdateMeResDto, *exceptions.Exception) {
	claims, exception := util.ParseAccessToken(reqDto.AccessToken)
	if exception != nil {
		return nil, exception
	}

	userRepository := repositories.NewUserRepository(nil)
	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	user, exception := userRepository.UpdateOneById(userId, inputs.PartialUpdateUserInput{
		Values: inputs.UpdateUserInput{
			DisplayName: reqDto.Values.DisplayName,
			Status:      reqDto.Values.Status,
		},
		SetNull: reqDto.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMeResDto{UpdatedAt: user.UpdatedAt}, nil
}

// func UpdatePlan(reqDto *dtos.UpdatePlanReqDto) (*dtos.UpdatePlanResDto, *exceptions.Exception) {

// }
