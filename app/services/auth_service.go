package services

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	dtos "notezy-backend/app/dtos"
	emails "notezy-backend/app/emails"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	tokens "notezy-backend/app/tokens"
	util "notezy-backend/app/util"
	constants "notezy-backend/shared/constants"

	authsql "notezy-backend/app/models/sql/auth"
)

/* ============================== Interface & Instance ============================== */

type AuthServiceInterface interface {
	Register(reqDto *dtos.RegisterReqDto) (*dtos.RegisterResDto, *exceptions.Exception)
	Login(reqDto *dtos.LoginReqDto) (*dtos.LoginResDto, *exceptions.Exception)
	Logout(reqDto *dtos.LogoutReqDto) (*dtos.LogoutResDto, *exceptions.Exception)
	SendAuthCode(reqDto *dtos.SendAuthCodeReqDto) (*dtos.SendAuthCodeResDto, *exceptions.Exception)
	ValidateEmail(reqDto *dtos.ValidateEmailReqDto) (*dtos.ValidateEmailResDto, *exceptions.Exception)
	ResetEmail(reqDto *dtos.ResetEmailReqDto) (*dtos.ResetEmailResDto, *exceptions.Exception)
	ForgetPassword(reqDto *dtos.ForgetPasswordReqDto) (*dtos.ForgetPasswordResDto, *exceptions.Exception)
	DeleteMe(reqDto *dtos.DeleteMeReqDto) (*dtos.DeleteMeResDto, *exceptions.Exception)
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthServiceInterface {
	return &AuthService{
		db: db,
	}
}

/* ============================== Auxiliary Functions ============================== */

func (s *AuthService) hashPassword(password string) (string, *exceptions.Exception) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", exceptions.Util.FailedToGenerateHashValue().WithError(err)
	}

	return string(bytes), nil
}

func (s *AuthService) checkPasswordHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

/* ============================== Services ============================== */

func (s *AuthService) Register(reqDto *dtos.RegisterReqDto) (*dtos.RegisterResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Auth.InvalidDto().WithError(err)
	}

	// Start transaction
	tx := s.db.Begin()
	userRepository := repositories.NewUserRepository(tx)
	userInfoRepository := repositories.NewUserInfoRepository(tx)
	userAccountRepository := repositories.NewUserAccountRepository(tx)
	userSettingRepository := repositories.NewUserSettingRepository(tx)

	hashedPassword, exception := s.hashPassword(reqDto.Password)
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
	accessToken, exception := tokens.GenerateAccessToken(
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
	refreshToken, exception := tokens.GenerateRefreshToken(
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
		Values: inputs.UpdateUserInput{
			RefreshToken: refreshToken,
		},
		SetNull: nil,
	})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user info
	_, exception = userInfoRepository.CreateOneByUserId(*newUserId, inputs.CreateUserInfoInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user account
	_, exception = userAccountRepository.CreateOneByUserId(*newUserId, inputs.CreateUserAccountInput{AuthCode: authCode, AuthCodeExpiredAt: authCodeExpiredAt})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user setting
	_, exception = userSettingRepository.CreateOneByUserId(*newUserId, inputs.CreateUserSettingInput{})
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

func (s *AuthService) Login(reqDto *dtos.LoginReqDto) (*dtos.LoginResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userRepository := repositories.NewUserRepository(s.db)

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

	if !s.checkPasswordHash(user.Password, reqDto.Password) {
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

	accessToken, exception := tokens.GenerateAccessToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := tokens.GenerateRefreshToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
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

func (s *AuthService) Logout(reqDto *dtos.LogoutReqDto) (*dtos.LogoutResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Auth.InvalidDto().WithError(err)
	}

	userRepository := repositories.NewUserRepository(s.db)

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
		return nil, exception
	}

	return &dtos.LogoutResDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) SendAuthCode(reqDto *dtos.SendAuthCodeReqDto) (*dtos.SendAuthCodeResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err).Log()
	}

	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	output := struct {
		Name      string `json:"name"`
		UserAgent string `json:"userAgent"`
	}{}
	result := s.db.Raw(authsql.UpdateAuthCodeQuery, authCode, authCodeExpiredAt, reqDto.Email).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToUpdate().WithError(err)
	}

	exception := emails.SendValidationEmail(reqDto.Email, output.Name, authCode, output.UserAgent, authCodeExpiredAt)
	if exception != nil {
		return nil, exception
	}

	return &dtos.SendAuthCodeResDto{
		AuthCodeExpiredAt: authCodeExpiredAt,
		UpdatedAt:         time.Now(),
	}, nil
}

func (s *AuthService) ValidateEmail(reqDto *dtos.ValidateEmailReqDto) (*dtos.ValidateEmailResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	output := struct {
		UpdatedAt time.Time `json:"updatedAt"`
	}{}
	result := s.db.Raw(authsql.ValidateEmailQuery, reqDto.UserId, reqDto.AuthCode).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}

	return &dtos.ValidateEmailResDto{
		UpdatedAt: output.UpdatedAt,
	}, nil
}

func (s *AuthService) ResetEmail(reqDto *dtos.ResetEmailReqDto) (*dtos.ResetEmailResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userAccountRepository := repositories.NewUserAccountRepository(nil)

	output := struct {
		UpdatedAt time.Time `json:"updatedAt"`
	}{}
	result := s.db.Raw(authsql.ResetEmailQuery, reqDto.NewEmail, reqDto.AuthCode, reqDto.UserId).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}

	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	_, exception := userAccountRepository.UpdateOneByUserId(
		reqDto.UserId,
		inputs.PartialUpdateUserAccountInput{
			Values: inputs.UpdateUserAccountInput{
				AuthCode:          &authCode,
				AuthCodeExpiredAt: &authCodeExpiredAt,
			},
			SetNull: nil,
		},
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.ResetEmailResDto{
		UpdatedAt: output.UpdatedAt,
	}, nil
}

func (s *AuthService) ForgetPassword(reqDto *dtos.ForgetPasswordReqDto) (*dtos.ForgetPasswordResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userRepository := repositories.NewUserRepository(s.db)

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

	accessToken, exception := tokens.GenerateAccessToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := tokens.GenerateRefreshToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}

	// update the access token of the user
	exception = caches.UpdateUserDataCache(user.Id, caches.UpdateUserDataCacheDto{AccessToken: accessToken})
	if exception != nil {
		return nil, exception
	}

	hashedPassword, exception := s.hashPassword(reqDto.NewPassword)
	if exception != nil {
		return nil, exception
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := userRepository.UpdateOneById(
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Password:     &hashedPassword,
				RefreshToken: refreshToken,
				UserAgent:    &reqDto.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		})
	if exception != nil {
		return nil, exception
	}

	return &dtos.ForgetPasswordResDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) DeleteMe(reqDto *dtos.DeleteMeReqDto) (*dtos.DeleteMeResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	output := struct {
		DeletedAt time.Time `json:"deletedAt" gorm:"deleted_at"`
	}{}
	result := s.db.Raw(authsql.DeleteMeQuery, reqDto.UserId, reqDto.AuthCode).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}

	return &dtos.DeleteMeResDto{
		DeletedAt: output.DeletedAt,
	}, nil
}
