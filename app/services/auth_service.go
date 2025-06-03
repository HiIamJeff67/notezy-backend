package services

import (
	"time"

	uuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	"notezy-backend/app/models/enums"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/operations"
	"notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Auxiliary Function ============================== */
func hashPassword(password string) (string, *exceptions.Exception) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", exceptions.Util.FailedToGenerateHashValue().WithError(err)
	}

	return string(bytes), nil
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/* ============================== Service ============================== */
func Register(reqDto *dtos.RegisterReqDto) (*dtos.RegisterResDto, *exceptions.Exception) {
	// Start transaction
	tx := models.NotezyDB.Begin()

	hashedPassword, exception := hashPassword(reqDto.Password)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	createUserInputData := inputs.CreateUserInput{
		Name:        reqDto.Name,
		DisplayName: util.GenerateRandomFakeName(),
		Email:       reqDto.Email,
		Password:    hashedPassword,
	}
	newUserId, exception := operations.CreateUser(tx, createUserInputData)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Generate tokens
	accessToken, exception := util.GenerateAccessToken((*newUserId).String(), createUserInputData.Name, createUserInputData.Email)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	refreshToken, exception := util.GenerateRefreshToken((*newUserId).String(), createUserInputData.Name, createUserInputData.Email)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Update user refresh token
	newUser, exception := operations.UpdateUserById(tx, *newUserId, inputs.UpdateUserInput{RefreshToken: refreshToken})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user info
	_, exception = operations.CreateUserInfoByUserId(tx, *newUserId, inputs.CreateUserInfoInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user account
	_, exception = operations.CreateUserAccountByUserId(tx, *newUserId, inputs.CreateUserAccountInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user setting
	_, exception = operations.CreateUserSettingByUserId(tx, *newUserId, inputs.CreateUserSettingInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.User.FailedToCommitTransaction().WithError(err)
	}

	// Create user data cache
	exception = caches.SetUserDataCache(
		*newUserId,
		caches.UserDataCache{
			Name:               newUser.Name,
			DisplayName:        newUser.DisplayName,
			Email:              newUser.Email,
			AccessToken:        *accessToken,
			Role:               newUser.Role,
			Plan:               newUser.Plan,
			Status:             newUser.Status,
			AvatarURL:          "",
			Theme:              enums.Theme_System,
			Language:           enums.Language_English,
			GeneralSettingCode: 0,
			PrivacySettingCode: 0,
			UpdatedAt:          newUser.UpdatedAt,
		},
	)
	if exception != nil {
		exception.Log()
	}

	return &dtos.RegisterResDto{
		AccessToken: *accessToken,
		CreatedAt:   newUser.CreatedAt,
	}, nil
}

func Login(reqDto *dtos.LoginReqDto) (*dtos.LoginResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	// check if we can just use the access token
	if reqDto.AccessToken != nil {
		claims, exception := util.ParseAccessToken(*reqDto.AccessToken)
		if exception != nil {
			return nil, exception
		}
		var userId uuid.UUID
		if err := userId.Scan(claims.Id); err != nil {
			return nil, exceptions.User.InvalidInput().WithError(err)
		}
		userCacheData, exception := caches.GetUserDataCache(userId)
		if exception != nil {
			return nil, exception
		}
		if userCacheData.AccessToken != *reqDto.AccessToken {
			return nil, exceptions.Auth.WrongAccessToken()
		}

		return &dtos.LoginResDto{AccessToken: *reqDto.AccessToken, CreatedAt: claims.ExpiresAt.Time}, nil
	}

	// else check if we can use the refresh token
	if reqDto.RefreshToken != nil {
		claims, exception := util.ParseRefreshToken(*reqDto.RefreshToken)
		if exception != nil {
			return nil, exception
		}

		var userId uuid.UUID
		if err := userId.Scan(claims.Id); err != nil {
			return nil, exceptions.User.InvalidInput().WithError(err)
		}
		user, exception := operations.GetUserById(nil, userId)
		if exception != nil {
			return nil, exception
		}
		if user.RefreshToken != *reqDto.RefreshToken {
			return nil, exceptions.Auth.WrongRefreshToken()
		}

		accessToken, exception := util.GenerateAccessToken(claims.Id, claims.Name, claims.Email)
		if exception != nil {
			return nil, exception
		}
		return &dtos.LoginResDto{AccessToken: *accessToken, CreatedAt: time.Now()}, nil
	}

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	if util.IsAlphaNumberString(*reqDto.Account) {
		if user, exception = operations.GetUserByName(nil, *reqDto.Account); exception != nil {
			return nil, exception
		}
	} else if util.IsEmailString(*reqDto.Account) {
		if user, exception = operations.GetUserByEmail(nil, *reqDto.Account); exception != nil {
			return nil, exception
		}
	} else {
		return nil, exceptions.Auth.InvalidDto()
	}

	if !checkPasswordHash(*reqDto.Password, user.Password) {
		return nil, exceptions.Auth.WrongPassword()
	}

	accessToken, exception := util.GenerateAccessToken(user.Id.String(), user.Name, user.Email)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := util.GenerateRefreshToken(user.Id.String(), user.Name, user.Email)
	if exception != nil {
		return nil, exception
	}

	// update the access token of the user
	if exception = caches.UpdateUserDataCache(user.Id, caches.UpdateUserDataCacheDto{AccessToken: accessToken}); exception != nil {
		return nil, exception
	}
	// update the refresh token of the user
	if _, exception = operations.UpdateUserById(nil, user.Id, inputs.UpdateUserInput{RefreshToken: refreshToken}); exception != nil {
		return nil, exception
	}

	return &dtos.LoginResDto{AccessToken: *accessToken, CreatedAt: time.Now()}, nil
}
