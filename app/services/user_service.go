package services

import (
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
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
