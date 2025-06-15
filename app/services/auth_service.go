package services

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	emails "notezy-backend/app/emails"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	enums "notezy-backend/app/models/enums"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	constants "notezy-backend/app/shared/constants"
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

func checkPasswordHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

/* ============================== Service ============================== */
func Register(reqDto *dtos.RegisterReqDto) (*dtos.RegisterResDto, *exceptions.Exception) {
	// Start transaction
	tx := models.NotezyDB.Begin()
	userRepository := repositories.NewUserRepository(tx)
	userInfoRepository := repositories.NewUserInfoRepository(tx)
	userAccountRepository := repositories.NewUserAccountRepository(tx)
	userSettingRepository := repositories.NewUserSettingRepository(tx)

	hashedPassword, exception := hashPassword(reqDto.Password)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	createUserInputData := inputs.CreateUserInput{
		Name:        reqDto.Name,
		DisplayName: util.GenerateRandomFakeName(), // we generate a default display name for the new user
		Email:       reqDto.Email,
		Password:    hashedPassword,
		UserAgent:   reqDto.UserAgent,
	}
	newUserId, exception := userRepository.CreateOne(createUserInputData)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Generate accessToken
	accessToken, exception := util.GenerateAccessToken(
		(*newUserId).String(),
		createUserInputData.Name,
		createUserInputData.Email,
		createUserInputData.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	// Generate refreshToken
	refreshToken, exception := util.GenerateRefreshToken(
		(*newUserId).String(),
		createUserInputData.Name,
		createUserInputData.Email,
		createUserInputData.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	// Generate authCode and its expiration time
	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)

	// Update user refresh token
	newUser, exception := userRepository.UpdateOneById(*newUserId, inputs.PartialUpdateUserInput{
		Values:  inputs.UpdateUserInput{RefreshToken: refreshToken},
		SetNull: nil,
	})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user info
	exception = userInfoRepository.CreateOneByUserId(*newUserId, inputs.CreateUserInfoInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user account
	exception = userAccountRepository.CreateOneByUserId(*newUserId, inputs.CreateUserAccountInput{AuthCode: authCode, AuthCodeExpiredAt: authCodeExpiredAt})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user setting
	exception = userSettingRepository.CreateOneByUserId(*newUserId, inputs.CreateUserSettingInput{})
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
			Language:           enums.Language_English,
			GeneralSettingCode: 0,
			PrivacySettingCode: 0,
			UpdatedAt:          newUser.UpdatedAt,
		},
	)
	if exception != nil {
		exception.Log()
	}

	// ssend the welcome email to the registered user
	exception = emails.SendWelcomeEmail(newUser.Email, newUser.Name, newUser.Status.String())
	if exception != nil {
		exception.Log()
	}

	return &dtos.RegisterResDto{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		CreatedAt:    newUser.CreatedAt,
	}, nil
}

func Login(reqDto *dtos.LoginReqDto) (*dtos.LoginResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userRepository := repositories.NewUserRepository(nil)

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	if util.IsAlphaNumberString(reqDto.Account) { // if the account field contains user name
		if user, exception = userRepository.GetOneByName(reqDto.Account); exception != nil {
			return nil, exception
		}
	} else if util.IsEmailString(reqDto.Account) { // if the account field contains email
		if user, exception = userRepository.GetOneByEmail(reqDto.Account); exception != nil {
			return nil, exception
		}
	} else {
		return nil, exceptions.Auth.InvalidDto()
	}

	if user.BlockLoginUtil.After(time.Now()) {
		return nil, exceptions.Auth.LoginBlockedDueToTryingTooManyTimes(user.BlockLoginUtil)
	}

	if !checkPasswordHash(user.Password, reqDto.Password) {
		newLoginCount := user.LoginCount + 1
		blockLoginUntil := util.GetLoginBlockedUntilByLoginCount(newLoginCount)
		updateInvalidUserInput := inputs.UpdateUserInput{
			LoginCount: &newLoginCount,
		}
		if blockLoginUntil != nil {
			updateInvalidUserInput.BlockLoginUtil = blockLoginUntil
		}

		_, exception := userRepository.UpdateOneById(user.Id, inputs.PartialUpdateUserInput{
			Values:  updateInvalidUserInput,
			SetNull: nil,
		})
		if exception != nil {
			return nil, exception
		}

		if blockLoginUntil != nil {
			return nil, exceptions.Auth.LoginBlockedDueToTryingTooManyTimes(*blockLoginUntil)
		}

		return nil, exceptions.Auth.WrongPassword()
	}

	if user.UserAgent != reqDto.UserAgent {
		// send a security email to warn the user
	}

	accessToken, exception := util.GenerateAccessToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := util.GenerateRefreshToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}

	// update the access token of the user
	exception = caches.UpdateUserDataCache(user.Id, caches.UpdateUserDataCacheDto{AccessToken: accessToken})
	if exception != nil {
		return nil, exception
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := userRepository.UpdateOneById(
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &user.PrevStatus,
				RefreshToken: refreshToken,
				UserAgent:    &reqDto.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		})
	if exception != nil {
		return nil, exception
	}

	return &dtos.LoginResDto{
		AccessToken: *accessToken,
		UpdatedAt:   updatedUser.UpdatedAt,
	}, nil
}

func Logout(reqDto *dtos.LogoutReqDto) (*dtos.LogoutResDto, *exceptions.Exception) {
	userRepository := repositories.NewUserRepository(nil)

	offlineStatus := enums.UserStatus_Offline
	emptyString := ""
	updatedUser, exception := userRepository.UpdateOneById(
		reqDto.UserId,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &offlineStatus,
				RefreshToken: &emptyString,
			},
			SetNull: nil,
		})
	if exception != nil {
		return nil, exception
	}

	exception = caches.DeleteUserDataCache(reqDto.UserId)
	if exception != nil {
		exception.Log()
	}

	return &dtos.LogoutResDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func SendAuthCode(reqDto *dtos.SendAuthCodeReqDto) (*dtos.SendAuthCodeResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err).Log()
	}

	db := models.NotezyDB

	// Use CTE(Common Table Expression) tp speed up
	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	result := struct {
		Name      string `json:"name"`
		UserAgent string `json:"userAgent"`
	}{}
	err := db.Raw(`
        UPDATE`+schemas.UserAccount{}.TableName()+` ua
		SET auth_code = ?, auth_code_expired_at = ?
		FROM `+schemas.User{}.TableName()+` u
		WHERE ua.user_id = u.id AND u.email = ?
		RETURNING u.name, u.user_agent
    `, authCode, authCodeExpiredAt, reqDto.Email).Scan(&result).Error
	if err != nil {
		return nil, exceptions.UserAccount.FailedToUpdate().WithError(err)
	}

	exception := emails.SendValidationEmail(reqDto.Email, result.Name, authCode, result.UserAgent, authCodeExpiredAt)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SendAuthCodeResDto{
		AuthCodeExpiredAt: authCodeExpiredAt,
		UpdatedAt:         time.Now(),
	}, nil
}

func ResetEmail() {}

func ResetPassword() {}
