package services

import (
	"sync"
	"time"

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
		updatedUserInput := models.UpdateUserInput{RefreshToken: &refreshToken}
		_, exception := models.UpdateUserById(tx, newUser.Id, updatedUserInput)
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
		userDataCache := &caches.UserDataCache{
			Name:               newUser.Name,
			DisplayName:        newUser.DisplayName,
			Email:              newUser.Email,
			AccessToken:        accessToken,
			Role:               newUser.Role,
			Plan:               newUser.Plan,
			Status:             newUser.Status,
			AvatarURL:          nil,
			Theme:              models.Theme_System,
			Language:           models.Language_English,
			GeneralSettingCode: 0,
			PrivacySettingCode: 0,
		}
		exception = caches.SetUserDataCache(newUser.Id, userDataCache)
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
		return nil, err
	case <-time.After(constants.RegisterTimeoutDuration):
		// deal with timeout exception
		return nil, exceptions.Util.Timeout(constants.RegisterTimeoutDuration)
	case <-done: // case when it receive the sign from done channel
		return &dtos.RegisterResDto{AccessToken: accessToken, CreatedAt: newUser.CreatedAt}, nil
	}
}

// func Login() {}
