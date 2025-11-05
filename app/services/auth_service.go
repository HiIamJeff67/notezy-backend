package services

import (
	"context"
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
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"

	authsql "notezy-backend/app/models/sql/auth"
	usersql "notezy-backend/app/models/sql/user"
)

/* ============================== Interface & Instance ============================== */

type AuthServiceInterface interface {
	Register(ctx context.Context, reqDto *dtos.RegisterReqDto) (*dtos.RegisterResDto, *exceptions.Exception)
	Login(ctx context.Context, reqDto *dtos.LoginReqDto) (*dtos.LoginResDto, *exceptions.Exception)
	Logout(ctx context.Context, reqDto *dtos.LogoutReqDto) (*dtos.LogoutResDto, *exceptions.Exception)
	SendAuthCode(ctx context.Context, reqDto *dtos.SendAuthCodeReqDto) (*dtos.SendAuthCodeResDto, *exceptions.Exception)
	ValidateEmail(ctx context.Context, reqDto *dtos.ValidateEmailReqDto) (*dtos.ValidateEmailResDto, *exceptions.Exception)
	ResetEmail(ctx context.Context, reqDto *dtos.ResetEmailReqDto) (*dtos.ResetEmailResDto, *exceptions.Exception)
	ForgetPassword(ctx context.Context, reqDto *dtos.ForgetPasswordReqDto) (*dtos.ForgetPasswordResDto, *exceptions.Exception)
	DeleteMe(ctx context.Context, reqDto *dtos.DeleteMeReqDto) (*dtos.DeleteMeResDto, *exceptions.Exception)
}

type AuthService struct {
	db                    *gorm.DB
	userRepository        repositories.UserRepositoryInterface
	userInfoRepository    repositories.UserInfoRepositoryInterface
	userAccountRepository repositories.UserAccountRepositoryInterface
	userSettingRepository repositories.UserSettingRepositoryInterface
}

func NewAuthService(
	db *gorm.DB,
	userRepository repositories.UserRepositoryInterface,
	userInfoRepository repositories.UserInfoRepositoryInterface,
	userAccountRepository repositories.UserAccountRepositoryInterface,
	userSettingRepository repositories.UserSettingRepositoryInterface,
) AuthServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &AuthService{
		db:                    db,
		userRepository:        userRepository,
		userInfoRepository:    userInfoRepository,
		userAccountRepository: userAccountRepository,
		userSettingRepository: userSettingRepository,
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

/* ============================== Service Methods for Authentication ============================== */

func (s *AuthService) Register(
	ctx context.Context, reqDto *dtos.RegisterReqDto,
) (*dtos.RegisterResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Auth.InvalidDto().WithError(err)
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()

	hashedPassword, exception := s.hashPassword(reqDto.Body.Password)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	createUserInputData := inputs.CreateUserInput{
		Name:        reqDto.Body.Name,
		DisplayName: util.GenerateRandomFakeName(), // we generate a default display name for the new user
		Email:       reqDto.Body.Email,
		Password:    hashedPassword,
		UserAgent:   reqDto.Header.UserAgent,
	}
	newUserId, exception := s.userRepository.CreateOne(tx, createUserInputData)
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
	// Generate csrfToken
	csrfToken, exception := tokens.GenerateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Generate authCode and its expiration time
	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)

	// Update user refresh token
	newUser, exception := s.userRepository.UpdateOneById(tx, *newUserId, inputs.PartialUpdateUserInput{
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
	_, exception = s.userInfoRepository.CreateOneByUserId(tx, *newUserId, inputs.CreateUserInfoInput{})
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user account
	_, exception = s.userAccountRepository.CreateOneByUserId(
		tx,
		*newUserId,
		inputs.CreateUserAccountInput{
			AuthCode:          authCode,
			AuthCodeExpiredAt: authCodeExpiredAt,
		},
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Create user setting
	_, exception = s.userSettingRepository.CreateOneByUserId(tx, *newUserId, inputs.CreateUserSettingInput{})
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
			PublicId:           newUser.PublicId,
			Name:               newUser.Name,
			DisplayName:        newUser.DisplayName,
			Email:              newUser.Email,
			AccessToken:        *accessToken,
			CSRFToken:          *csrfToken,
			Role:               newUser.Role,
			Plan:               newUser.Plan,
			Status:             newUser.Status,
			AvatarURL:          "",
			Language:           enums.Language_English,
			GeneralSettingCode: 0,
			PrivacySettingCode: 0,
			CreatedAt:          newUser.CreatedAt,
			UpdatedAt:          newUser.UpdatedAt,
		},
	)
	if exception != nil {
		exception.Log()
	}

	// send the welcome email to the registered user

	if exception = emails.SyncSendWelcomeEmail(
		newUser.Email,
		newUser.Name,
		newUser.Status.String(),
	); exception != nil {
		exception.Log()
	}

	return &dtos.RegisterResDto{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		CSRFToken:    *csrfToken,
		CreatedAt:    newUser.CreatedAt,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context, reqDto *dtos.LoginReqDto,
) (*dtos.LoginResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	if util.IsAlphaAndNumberString(reqDto.Body.Account) { // if the account field contains user name
		if user, exception = s.userRepository.GetOneByName(db, reqDto.Body.Account, nil); exception != nil {
			return nil, exception
		}
	} else if util.IsEmailString(reqDto.Body.Account) { // if the account field contains email
		if user, exception = s.userRepository.GetOneByEmail(db, reqDto.Body.Account, nil); exception != nil {
			return nil, exception
		}
	}

	if user == nil {
		return nil, exceptions.Auth.InvalidDto()
	}

	if user.BlockLoginUntil.After(time.Now()) {
		return nil, exceptions.Auth.LoginBlockedDueToTryingTooManyTimes(user.BlockLoginUntil)
	}

	if !s.checkPasswordHash(user.Password, reqDto.Body.Password) {
		newLoginCount := user.LoginCount + 1
		blockLoginUntil, exception := util.GetLoginBlockedUntilByLoginCount(newLoginCount)
		if exception != nil {
			return nil, exception
		}
		updateInvalidUserInput := inputs.UpdateUserInput{
			LoginCount: &newLoginCount,
		}
		updateInvalidUserInput.BlockLoginUtil = blockLoginUntil // we don't care if blockLoginUntil is nil or not, since we always set the SetNull to nil

		_, exception = s.userRepository.UpdateOneById(db, user.Id, inputs.PartialUpdateUserInput{
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

	if user.UserAgent != reqDto.Header.UserAgent {
		// send a security email to warn the user
		if exception := emails.SyncSendSecurityAlertEmail(
			user.Email,
			user.Name,
			user.Status.String(),
			"Login in Different Place",
			"Your account has a recent login action in other place",
			time.Now(),
			"",
		); exception != nil {
			exception.Log()
		}
	}

	accessToken, exception := tokens.GenerateAccessToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := tokens.GenerateRefreshToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	csrfToken, exception := tokens.GenerateCSRFToken()
	if exception != nil {
		return nil, exception
	}

	// check if the user data cache exists
	if _, exception := caches.GetUserDataCache(user.Id); exception == nil {
		// then just update the existing user data cache
		if exception = caches.UpdateUserDataCache(
			user.Id,
			caches.UpdateUserDataCacheDto{
				AccessToken: accessToken,
			},
		); exception != nil {
			exception.Log()
		}
	} else { // else if it does not exist
		// then we have to first get the relative data from differenct tables
		// we done this by one custom sql so it's not that slow...
		// once we have the required data, we set it as the user data cache
		output := struct {
			PublicId           string           `gorm:"public_id"`
			Name               string           `gorm:"name"`
			DisplayName        string           `gorm:"display_name"`
			Email              string           `gorm:"email"`
			Role               enums.UserRole   `gorm:"role"`
			Plan               enums.UserPlan   `gorm:"plan"`
			Status             enums.UserStatus `gorm:"status"`
			AvatarURL          *string          `gorm:"avatar_url"`
			Language           enums.Language   `gorm:"language"`
			GeneralSettingCode int64            `gorm:"general_setting_code"`
			PrivacySettingCode int64            `gorm:"privacy_setting_code"`
			CreatedAt          time.Time        `gorm:"created_at"`
			UpdatedAt          time.Time        `gorm:"updated_at"`
		}{}
		result := db.Raw(usersql.GetUserDataCacheByIdSQL, user.Id).Scan(&output)
		if err := result.Error; err != nil || result.RowsAffected == 0 {
			exception := exceptions.User.NotFound().WithError(err).Log()
			if err != nil {
				exception.WithError(err)
			}
			// should be exists in database
			return nil, exception
		}

		newUserDataCache := caches.UserDataCache{
			PublicId:           output.PublicId,
			Name:               output.Name,
			DisplayName:        output.DisplayName,
			Email:              output.Email,
			AccessToken:        *accessToken,
			CSRFToken:          *csrfToken,
			Role:               output.Role,
			Plan:               output.Plan,
			Status:             output.Status,
			AvatarURL:          "",
			Language:           output.Language,
			GeneralSettingCode: output.GeneralSettingCode,
			PrivacySettingCode: output.PrivacySettingCode,
			CreatedAt:          output.CreatedAt,
			UpdatedAt:          output.UpdatedAt,
		}
		if output.AvatarURL != nil {
			newUserDataCache.AvatarURL = *output.AvatarURL
		}
		exception := caches.SetUserDataCache(
			user.Id,
			newUserDataCache,
		)
		if exception != nil {
			return nil, exception.Log()
		}
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := s.userRepository.UpdateOneById(
		db,
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &user.PrevStatus,
				RefreshToken: refreshToken,
				UserAgent:    &reqDto.Header.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		})
	if exception != nil {
		return nil, exception
	}

	return &dtos.LoginResDto{
		AccessToken:  *accessToken,
		RefreshToken: updatedUser.RefreshToken,
		CSRFToken:    *csrfToken,
		UpdatedAt:    updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context, reqDto *dtos.LogoutReqDto,
) (*dtos.LogoutResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Auth.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	offlineStatus := enums.UserStatus_Offline
	emptyString := ""
	updatedUser, exception := s.userRepository.UpdateOneById(
		db,
		reqDto.ContextFields.UserId,
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

	exception = caches.DeleteUserDataCache(reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.LogoutResDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) SendAuthCode(
	ctx context.Context, reqDto *dtos.SendAuthCodeReqDto,
) (*dtos.SendAuthCodeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	blockAuthCodeUntil := util.GetAuthCodeBlockUntil()
	output := struct {
		Name               string    `json:"name"`
		UserAgent          string    `json:"userAgent"`
		BlockAuthCodeUntil time.Time `json:"blockAuthCodeUntil"`
		Now                time.Time `json:"now"`
	}{}
	result := db.Raw(authsql.UpdateAuthCodeSQL, authCode, authCodeExpiredAt, blockAuthCodeUntil, reqDto.Body.Email).Scan(&output)
	if err := result.Error; err != nil || result.RowsAffected == 0 {
		exception := exceptions.Auth.AuthCodeBlockedDueToTryingTooManyTimes(output.BlockAuthCodeUntil)
		if err != nil {
			exception.WithError(err)
		}
		return nil, exception
	}

	if exception := emails.SyncSendValidationEmail(
		reqDto.Body.Email,
		output.Name,
		authCode,
		output.UserAgent,
		authCodeExpiredAt,
	); exception != nil {
		return nil, exception
	}

	return &dtos.SendAuthCodeResDto{
		AuthCodeExpiredAt:  authCodeExpiredAt,
		BlockAuthCodeUntil: blockAuthCodeUntil,
		UpdatedAt:          time.Now(),
	}, nil
}

func (s *AuthService) ValidateEmail(
	ctx context.Context, reqDto *dtos.ValidateEmailReqDto,
) (*dtos.ValidateEmailResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	output := struct {
		UpdatedAt time.Time `json:"updatedAt"`
	}{}
	result := db.Raw(authsql.ValidateEmailSQL, reqDto.ContextFields.UserId, reqDto.Body.AuthCode).Scan(&output)
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

func (s *AuthService) ResetEmail(
	ctx context.Context, reqDto *dtos.ResetEmailReqDto,
) (*dtos.ResetEmailResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	output := struct {
		UpdatedAt time.Time `json:"updatedAt"`
	}{}
	result := db.Raw(authsql.ResetEmailSQL, reqDto.Body.NewEmail, reqDto.Body.AuthCode, reqDto.ContextFields.UserId).Scan(&output)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}

	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}

	authCode := util.GenerateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	_, exception := s.userAccountRepository.UpdateOneByUserId(
		db,
		reqDto.ContextFields.UserId,
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

func (s *AuthService) ForgetPassword(
	ctx context.Context, reqDto *dtos.ForgetPasswordReqDto,
) (*dtos.ForgetPasswordResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	var preloads = []schemas.UserRelation{schemas.UserRelation_UserAccount, schemas.UserRelation_UserInfo, schemas.UserRelation_UserSetting}
	if util.IsEmailString(reqDto.Body.Account) { // if the account field contains email
		if user, exception = s.userRepository.GetOneByEmail(db, reqDto.Body.Account, preloads); exception != nil {
			return nil, exception
		}
	} else if util.IsAlphaAndNumberString(reqDto.Body.Account) { // if the account field contains user name
		if user, exception = s.userRepository.GetOneByName(db, reqDto.Body.Account, preloads); exception != nil {
			return nil, exception
		}
	} else {
		return nil, exceptions.Auth.InvalidDto()
	}

	if reqDto.Body.AuthCode != user.UserAccount.AuthCode {
		return nil, exceptions.Auth.WrongAuthCode()
	}

	accessToken, exception := tokens.GenerateAccessToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	refreshToken, exception := tokens.GenerateRefreshToken(user.Id.String(), user.Name, user.Email, user.UserAgent)
	if exception != nil {
		return nil, exception
	}
	csrfToken, exception := tokens.GenerateCSRFToken()
	if exception != nil {
		return nil, exception
	}

	// update the access token of the user
	exception = caches.UpdateUserDataCache(user.Id, caches.UpdateUserDataCacheDto{AccessToken: accessToken})
	if exception != nil {
		exception.Log() // if the cache does not exist the user, then just skip this update operation
		// and also try to set the new user cache data
		exception = caches.SetUserDataCache(user.Id, caches.UserDataCache{
			PublicId:           user.PublicId,
			Name:               user.Name,
			DisplayName:        user.DisplayName,
			Email:              user.Email,
			AccessToken:        *accessToken,
			CSRFToken:          *csrfToken,
			Role:               user.Role,
			Plan:               user.Plan,
			Status:             user.Status,
			AvatarURL:          *user.UserInfo.AvatarURL,
			Language:           user.UserSetting.Language,
			GeneralSettingCode: user.UserSetting.GeneralSettingCode,
			PrivacySettingCode: user.UserSetting.PrivacySettingCode,
			CreatedAt:          user.CreatedAt,
			UpdatedAt:          user.UpdatedAt,
		})
		if exception != nil {
			exception.Log() // if the set operation also failed, then just log it without abort the following
		}
	}

	hashedPassword, exception := s.hashPassword(reqDto.Body.NewPassword)
	if exception != nil {
		return nil, exception
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := s.userRepository.UpdateOneById(
		db,
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Password:     &hashedPassword,
				RefreshToken: refreshToken,
				UserAgent:    &reqDto.Header.UserAgent,
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

func (s *AuthService) DeleteMe(
	ctx context.Context, reqDto *dtos.DeleteMeReqDto,
) (*dtos.DeleteMeResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	db := s.db.WithContext(ctx)

	output := struct {
		DeletedAt time.Time `json:"deletedAt" gorm:"deleted_at"`
	}{}
	result := db.Raw(authsql.DeleteMeSQL, reqDto.ContextFields.UserId, reqDto.Body.AuthCode).Scan(&output)
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
