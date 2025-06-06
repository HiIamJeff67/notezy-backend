package services

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	"notezy-backend/app/emails"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	enums "notezy-backend/app/models/enums"
	inputs "notezy-backend/app/models/inputs"
	operations "notezy-backend/app/models/operations"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
	"notezy-backend/global/constants"
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
		UserAgent:   reqDto.UserAgent,
	}
	newUserId, exception := operations.CreateUser(tx, createUserInputData)
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
	_, exception = operations.CreateUserAccountByUserId(tx, *newUserId, inputs.CreateUserAccountInput{AuthCode: authCode, AuthCodeExpiredAt: authCodeExpiredAt})
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

	// ssend the welcome email to the registered user
	emails.SendWelcomeEmail(newUser.Email)

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

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	if util.IsAlphaNumberString(reqDto.Account) { // if the account field contain user name
		if user, exception = operations.GetUserByName(nil, reqDto.Account); exception != nil {
			return nil, exception
		}
	} else if util.IsEmailString(reqDto.Account) { // if the account field contain email
		if user, exception = operations.GetUserByEmail(nil, reqDto.Account); exception != nil {
			return nil, exception
		}
	} else {
		return nil, exceptions.Auth.InvalidDto()
	}

	if !checkPasswordHash(user.Password, reqDto.Password) {
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
	updatedUser, exception := operations.UpdateUserById(
		nil,
		user.Id,
		inputs.UpdateUserInput{
			Status:       &user.PrevStatus,
			RefreshToken: refreshToken,
			UserAgent:    &reqDto.UserAgent,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.LoginResDto{
		AccessToken: *accessToken,
		UpdatedAt:   updatedUser.UpdatedAt,
	}, nil
}

func Logout(reqDto *dtos.LogoutReqDto) (*dtos.LogoutResDto, *exceptions.Exception) {
	claims, exception := util.ParseAccessToken(reqDto.AccessToken)
	if exception != nil {
		return nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	statusOffline := enums.UserStatus_Offline
	emptyRefreshToken := ""
	updatedUser, exception := operations.UpdateUserById(
		nil,
		userId,
		inputs.UpdateUserInput{
			Status:       &statusOffline,
			RefreshToken: &emptyRefreshToken,
		},
	)
	if exception != nil {
		return nil, exception
	}

	exception = caches.DeleteUserDataCache(userId)
	if exception != nil {
		exception.Log()
	}

	return &dtos.LogoutResDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}
