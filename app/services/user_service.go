package services

import (
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/models/inputs"
	operations "notezy-backend/app/models/operations"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

func FindMe(reqDto *dtos.FindMeReqDto) (*dtos.FindMeResDto, *exceptions.Exception) {
	claims, exception := util.ParseAccessToken(reqDto.AccessToken)
	if exception != nil {
		return nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(userId)
	if exception != nil {
		return nil, exception
	}

	return userDataCache, nil
}

// for temporary use
func FindAllUsers() (*[]schemas.User, *exceptions.Exception) {
	users, exception := operations.GetAllUsers(nil)
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

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	user, exception := operations.UpdateUserById(nil, userId, inputs.UpdateUserInput{
		DisplayName: &reqDto.DisplayName,
		Email:       &reqDto.Email,
		Password:    &reqDto.Password,
		Status:      &reqDto.Status,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMeResDto{UpdatedAt: user.UpdatedAt}, nil
}
