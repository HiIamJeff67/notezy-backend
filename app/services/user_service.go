package services

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/models/operations"
	"notezy-backend/app/models/schemas"
)

/* ============================== Services ============================== */
// func FindMe(reqDto *dtos.FindMeReqDto) (*dtos.FindMeResDto, *exceptions.Exception) {
// 	if reqDto.AccessToken != nil {
// 		claims, exception := util.ParseAccessToken(*reqDto.AccessToken)
// 		if exception != nil {
// 			return nil, exception
// 		}

// 		var userId uuid.UUID
// 		if err := userId.Scan(claims.Id); err != nil {
// 			return nil, exceptions.User.InvalidInput().WithError(err)
// 		}

// 		userDataCache, exception := caches.GetUserDataCache(userId)
// 		if exception != nil {
// 			return nil, exception
// 		}

// 		return userDataCache, nil
// 	} else if reqDto.RefreshToken != nil {
// 		claims, exception := util.ParseRefreshToken(*reqDto.RefreshToken)
// 		if exception != nil {
// 			return nil, exception
// 		}

// 		var userId uuid.UUID
// 		if err := userId.Scan(claims.Id); err != nil {
// 			return nil, exceptions.User.InvalidInput().WithError(err)
// 		}

// 		userDataCache, exception := caches.GetUserDataCache(userId)
// 		if exception != nil {
// 			return nil, exception
// 		}

// 		return userDataCache, nil
// 	}
// 	return nil, nil
// }

// for temporary use
func FindAllUsers() (*[]schemas.User, *exceptions.Exception) {
	users, exception := operations.GetAllUsers(nil)
	if exception != nil {
		return nil, exception
	}

	return users, nil
}
