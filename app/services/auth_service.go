package services

import (
	"sync"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"golang.org/x/crypto/bcrypt"

	"notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	util "notezy-backend/app/util"
	constants "notezy-backend/global/constants"
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
	db := models.NotezyDB
	tx := db.Begin() // start the transaction

	hashedPassword, exception := hashPassword(reqDto.CreateUserInputData.Password)
	if exception != nil {
		return nil, exception
	}

	reqDto.CreateUserInputData.Password = hashedPassword
	newUser, exception := models.CreateUser(tx, reqDto.CreateUserInputData)
	if exception != nil {
		return nil, exception
	}

	accessToken, exception := util.GenerateAccessToken(newUser.Id.UUID.String(), newUser.Name, newUser.Email)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := util.GenerateRefreshToken(newUser.Id.UUID.String(), newUser.Name, newUser.Email)
	if exception != nil {
		return nil, exception
	}

	var wg sync.WaitGroup
	// once the errCh receive an error, we're going to stop the entire process and do the transaction roll back
	errCh := make(chan *exceptions.Exception, 1)

	run := func(f func() *exceptions.Exception) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := f(); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	// update refreshToken in user table
	run(func() *exceptions.Exception {
		_, exception := models.UpdateUserById(tx, newUser.Id, models.UpdateUserInput{RefreshToken: refreshToken})
		return exception
	})

	// create user info
	run(func() *exceptions.Exception {
		_, exception := models.CreateUserInfoByUserId(tx, newUser.Id, reqDto.CreateUserInfoInputData)
		return exception
	})

	// create user account
	run(func() *exceptions.Exception {
		_, exception := models.CreateUserAccountByUserId(tx, newUser.Id, reqDto.CreateUserAccountInputData)
		return exception
	})

	// create user setting
	run(func() *exceptions.Exception {
		_, exception := models.CreateUserSettingByUserId(tx, newUser.Id, reqDto.CreateUserSettingInputData)
		return exception
	})

	// store user data into the cache
	run(func() *exceptions.Exception {
		exception = caches.SetUserDataCache(
			newUser.Id,
			caches.UserDataCache{
				Name:               newUser.Name,
				DisplayName:        newUser.DisplayName,
				Email:              newUser.Email,
				AccessToken:        *accessToken,
				Role:               newUser.Role,
				Plan:               newUser.Plan,
				Status:             newUser.Status,
				AvatarURL:          nil,
				Theme:              models.Theme_System,
				Language:           models.Language_English,
				GeneralSettingCode: 0,
				PrivacySettingCode: 0,
			},
		)
		return exception
	})

	done := make(chan struct{}) // to indicate a end sign
	go func() {
		wg.Wait()
		errCh <- exceptions.User.FailedToCommitTransaction().WithError(tx.Commit().Error)
		close(done) // to pass the end sign
	}()

	select {
	case err := <-errCh:
		tx.Rollback()
		close(errCh)
		return nil, err
	case <-time.After(constants.RegisterTimeoutDuration):
		// deal with timeout exception
		return nil, exceptions.Util.Timeout(constants.RegisterTimeoutDuration)
	case <-done: // case when it receive the sign from done channel
		return &dtos.RegisterResDto{AccessToken: *accessToken, CreatedAt: newUser.CreatedAt}, nil
	}
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
		user, exception := models.GetUserById(nil, userId)
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
	var user *models.User = nil
	var exception *exceptions.Exception = nil
	if util.IsAlphaNumberString(*reqDto.Account) {
		if user, exception = models.GetUserByName(nil, *reqDto.Account); exception != nil {
			return nil, exception
		}
	} else if util.IsEmailString(*reqDto.Account) {
		if user, exception = models.GetUserByEmail(nil, *reqDto.Account); exception != nil {
			return nil, exception
		}
	} else {
		return nil, exceptions.Auth.InvalidDto()
	}

	if !checkPasswordHash(*reqDto.Password, user.Password) {
		return nil, exceptions.Auth.WrongPassword()
	}

	accessToken, exception := util.GenerateAccessToken(user.Id.UUID.String(), user.Name, user.Email)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := util.GenerateRefreshToken(user.Id.UUID.String(), user.Name, user.Email)
	if exception != nil {
		return nil, exception
	}

	// update the access token of the user
	if exception = caches.UpdateUserDataCache(user.Id, caches.UpdateUserDataCacheDto{AccessToken: accessToken}); exception != nil {
		return nil, exception
	}
	// update the refresh token of the user
	if _, exception = models.UpdateUserById(nil, user.Id, models.UpdateUserInput{RefreshToken: refreshToken}); exception != nil {
		return nil, exception
	}

	return &dtos.LoginResDto{AccessToken: *accessToken, CreatedAt: time.Now()}, nil
}
