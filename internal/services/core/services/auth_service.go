package services

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	emaildto "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	authdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/auth"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata/inputs"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	emailtransport "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/email"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	snowflake "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/snowflake"
	stringutil "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/stringutil"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"

	authsql "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/sqls/auth"
	badgesql "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/sqls/badge"
	usersql "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/sqls/user"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, requestDto *authdto.RegisterRequestDto) (*authdto.RegisterResponseDto, *exceptions.Exception)
	RegisterViaGoogle(ctx context.Context, requestDto *authdto.RegisterViaGoogleRequestDto) (*authdto.RegisterViaGoogleResponseDto, *exceptions.Exception)
	Login(ctx context.Context, requestDto *authdto.LoginRequestDto) (*authdto.LoginResponseDto, *exceptions.Exception)
	LoginViaGoogle(ctx context.Context, requestDto *authdto.LoginViaGoogleRequestDto) (*authdto.LoginViaGoogleResponseDto, *exceptions.Exception)
	Logout(ctx context.Context, requestDto *authdto.LogoutRequestDto) (*authdto.LogoutResponseDto, *exceptions.Exception)
	SendAuthCode(ctx context.Context, requestDto *authdto.SendAuthCodeRequestDto) (*authdto.SendAuthCodeResponseDto, *exceptions.Exception)
	ValidateEmail(ctx context.Context, requestDto *authdto.ValidateEmailRequestDto) (*authdto.ValidateEmailResponseDto, *exceptions.Exception)
	ResetEmail(ctx context.Context, requestDto *authdto.ResetEmailRequestDto) (*authdto.ResetEmailResponseDto, *exceptions.Exception)
	ForgetPassword(ctx context.Context, requestDto *authdto.ForgetPasswordRequestDto) (*authdto.ForgetPasswordResponseDto, *exceptions.Exception)
	ResetMe(ctx context.Context, requestDto *authdto.ResetMeRequestDto) (*authdto.ResetMeResponseDto, *exceptions.Exception)
	DeleteMe(ctx context.Context, requestDto *authdto.DeleteMeRequestDto) (*authdto.DeleteMeResponseDto, *exceptions.Exception)
}

type AuthService struct {
	db                    *gorm.DB
	userRepository        repositories.UserRepositoryInterface
	userInfoRepository    repositories.UserInfoRepositoryInterface
	userAccountRepository repositories.UserAccountRepositoryInterface
	userSettingRepository repositories.UserSettingRepositoryInterface
	rootShelfRepository   repositories.RootShelfRepositoryInterface
	oauthService          OAuthServiceInterface
	emailClient           emailtransport.ClientInterface
}

func NewAuthService(
	db *gorm.DB,
	userRepository repositories.UserRepositoryInterface,
	userInfoRepository repositories.UserInfoRepositoryInterface,
	userAccountRepository repositories.UserAccountRepositoryInterface,
	userSettingRepository repositories.UserSettingRepositoryInterface,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
	oauthService OAuthServiceInterface,
	emailClient emailtransport.ClientInterface,
) AuthServiceInterface {
	if db == nil {
		db = data.NotezyDB
	}
	return &AuthService{
		db:                    db,
		userRepository:        userRepository,
		userInfoRepository:    userInfoRepository,
		userAccountRepository: userAccountRepository,
		userSettingRepository: userSettingRepository,
		rootShelfRepository:   rootShelfRepository,
		oauthService:          oauthService,
		emailClient:           emailClient,
	}
}

/* ============================== Auxiliary Functions ============================== */

var loginCountToBlockDurationMap = map[int32]time.Duration{
	3:  5 * time.Minute,
	5:  15 * time.Minute,
	7:  30 * time.Minute,
	10: 1 * time.Hour,
	15: 6 * time.Hour,
	20: 24 * time.Hour,
	30: 7 * 24 * time.Hour,
}

func (s *AuthService) generateAuthCode() string {
	randomNumber := rand.IntN(constants.MaxAuthCode + 1)
	code := strconv.Itoa(randomNumber)
	for len(code) < constants.MaxLengthOfAuthCode {
		code = "0" + code
	}
	return code
}

func (s *AuthService) generateRandomFakeDisplayName() string {
	gofakeit.Seed(0)
	return fmt.Sprintf("%s%s%d", gofakeit.AdjectiveDescriptive(), gofakeit.LastName(), gofakeit.Number(100000, 999999))
}

func (s *AuthService) getLoginBlockedUntilByLoginCount(loginCount int32) (*time.Time, *exceptions.Exception) {
	if loginCount < 0 {
		return nil, exceptions.New("InvalidLoginCount", "Auth", "GetLoginBlockedUntil", "Login count is invalid", http.StatusInternalServerError, true)
	}

	var blockDuration *time.Duration
	for count, duration := range loginCountToBlockDurationMap {
		if loginCount >= count {
			blockDuration = &duration
		}
	}
	if blockDuration == nil {
		return nil, nil
	}
	blockedUntil := time.Now().Add(*blockDuration)
	return &blockedUntil, nil
}

func (s *AuthService) getAuthCodeBlockUntil() time.Time {
	return time.Now().Add(60 * time.Second)
}

func (s *AuthService) hashPassword(password string) (string, *exceptions.Exception) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", exceptions.New(
			"FailedToGenerateHashValue",
			"Auth",
			"Hash",
			"Failed to generate a hash value",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	return string(bytes), nil
}

func (s *AuthService) comparePassword(hashedPassword string, password string) bool {
	return (bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))) == nil
}

func (s *AuthService) getOAuthFakeName() (string, *exceptions.Exception) {
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		return "", apiexceptions.Auth.FailedToCompileRegularExpression().WithOrigin(err)
	}
	fakeName := strings.ToLower(uuid.New().String())
	fakeName = reg.ReplaceAllString(fakeName, "")
	if len(fakeName) < 6 {
		fakeName = fakeName + snowflake.GenerateRepeatableID()
	}
	return fakeName, nil
}

func (s *AuthService) getOAuthFakePassword() (string, *exceptions.Exception) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
	if err != nil {
		return "", exceptions.New("FailedToGenerateHashValue", "Auth", "Hash", "Failed to generate a hash value", http.StatusInternalServerError, true).WithOrigin(err)
	}
	return string(bytes), nil
}

func (s *AuthService) generateAccessToken(userPublicId uuid.UUID, name string, email string, userAgent string) (*string, *exceptions.Exception) {
	token, err := sharedtokens.GenerateAccessToken(
		userPublicId.String(),
		sharedtokens.AccessTokenClaims{
			Name:      name,
			Email:     email,
			UserAgent: userAgent,
		},
	)
	if err != nil {
		return nil, exceptions.New(
			"GenerationFailed",
			"Token",
			"GenerateAccessToken",
			"Failed to generate the access token",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return token, nil
}

func (s *AuthService) generateRefreshToken(userPublicId uuid.UUID, name string, email string, userAgent string) (*string, *exceptions.Exception) {
	token, err := sharedtokens.GenerateRefreshToken(
		userPublicId.String(),
		sharedtokens.RefreshTokenClaims{
			Name:      name,
			Email:     email,
			UserAgent: userAgent,
		},
	)
	if err != nil {
		return nil, exceptions.New(
			"GenerationFailed",
			"Token",
			"GenerateRefreshToken",
			"Failed to generate the refresh token",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return token, nil
}

func (s *AuthService) generateCSRFToken() (*string, *exceptions.Exception) {
	token, err := sharedtokens.GenerateCSRFToken(sharedtokens.CSRFTokenClaims{})
	if err != nil {
		return nil, exceptions.New(
			"GenerationFailed",
			"Token",
			"GenerateCSRFToken",
			"Failed to generate the CSRF token",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return token, nil
}

/* ============================== Service Methods for Authentication ============================== */

func (s *AuthService) Register(
	ctx context.Context, reqDto *authdto.RegisterRequestDto,
) (*authdto.RegisterResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.Auth.InvalidDto().WithOrigin(err)
	}

	// put the hash part outside the transaction to decrease its blocking time from heavily hashing operation
	hashedPassword, exception := s.hashPassword(reqDto.Body.Password)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	createUserInput := inputs.CreateUserInput{
		Name:        reqDto.Body.Name,
		DisplayName: s.generateRandomFakeDisplayName(), // we generate a default display name for the new user
		Email:       reqDto.Body.Email,
		Password:    hashedPassword,
		UserAgent:   reqDto.Header.UserAgent,
	}
	newUserId, exception := s.userRepository.CreateOne(
		createUserInput,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	createdUser, exception := s.userRepository.GetOneById(
		*newUserId,
		nil,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newAccessToken, exception := s.generateAccessToken(
		createdUser.PublicId,
		createUserInput.Name,
		createUserInput.Email,
		createUserInput.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newRefreshToken, exception := s.generateRefreshToken(
		createdUser.PublicId,
		createUserInput.Name,
		createUserInput.Email,
		createUserInput.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newCSRFToken, exception := s.generateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	authCode := s.generateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)

	newUser, exception := s.userRepository.UpdateOneById(
		*newUserId,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				RefreshToken: newRefreshToken,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userInfoRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserInfoInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userAccountRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserAccountInput{
			AuthCode:          authCode,
			AuthCodeExpiredAt: authCodeExpiredAt,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userSettingRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserSettingInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	exception = userdata.NewUserDataCacheClient().Set(
		newUser.Name,
		userdata.UserDataCache{
			Id:                 *newUserId,
			PublicId:           newUser.PublicId,
			Name:               newUser.Name,
			DisplayName:        newUser.DisplayName,
			Email:              newUser.Email,
			AccessToken:        *newAccessToken,
			CSRFToken:          *newCSRFToken,
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
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
	}

	if exception = s.emailClient.SendWelcomeEmail(ctx, emaildto.SendWelcomeEmailRequestDto{
		To:       newUser.Email,
		UserName: newUser.Name,
		Status:   newUser.Status.String(),
	}); exception != nil {
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
	}

	return &authdto.RegisterResponseDto{
		PublicId:     newUser.PublicId,
		Name:         newUser.Name,
		DisplayName:  newUser.DisplayName,
		Email:        newUser.Email,
		AccessToken:  *newAccessToken,
		RefreshToken: *newRefreshToken,
		CSRFToken:    *newCSRFToken,
		CreatedAt:    newUser.CreatedAt,
	}, nil
}

func (s *AuthService) RegisterViaGoogle(
	ctx context.Context, reqDto *authdto.RegisterViaGoogleRequestDto,
) (*authdto.RegisterViaGoogleResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.Auth.InvalidDto().WithOrigin(err)
	}

	userInfo, exception := s.oauthService.GetGoogleUserInfo(ctx, reqDto.Body.AuthorizationCode)
	if exception != nil {
		return nil, exception
	}

	fakePassword, exception := s.getOAuthFakePassword()
	if exception != nil {
		return nil, exception
	}
	hashedPassword, exception := s.hashPassword(fakePassword)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	// try to generate fake name at most 5 times
	var newUserId *uuid.UUID
	createUserInput := inputs.CreateUserInput{
		Name:        "",
		DisplayName: s.generateRandomFakeDisplayName(), // we generate a default display name for the new user
		Email:       userInfo.Email,
		Password:    hashedPassword,
		UserAgent:   reqDto.Header.UserAgent,
	}
	for i := 0; i < constants.MaxRetriesOfGeneratingFakeName; i++ {
		fakeName, exception := s.getOAuthFakeName()
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

		if len(fakeName) > constants.MaxNameLength {
			fakeName = fakeName[:constants.MaxNameLength]
		}
		createUserInput.Name = fakeName

		newUserId, exception = s.userRepository.CreateOne(
			createUserInput,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		)
		if exception == nil {
			break
		}
	}
	if newUserId == nil {
		tx.Rollback()
		return nil, exception
	}
	createdUser, exception := s.userRepository.GetOneById(
		*newUserId,
		nil,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newAccessToken, exception := s.generateAccessToken(
		createdUser.PublicId,
		createUserInput.Name,
		createUserInput.Email,
		createUserInput.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newRefreshToken, exception := s.generateRefreshToken(
		createdUser.PublicId,
		createUserInput.Name,
		createUserInput.Email,
		createUserInput.UserAgent,
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newCSRFToken, exception := s.generateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	authCode := s.generateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)

	newUser, exception := s.userRepository.UpdateOneById(
		*newUserId,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				RefreshToken: newRefreshToken,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userInfoRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserInfoInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userAccountRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserAccountInput{
			AuthCode:          authCode,
			AuthCodeExpiredAt: authCodeExpiredAt,
			GoogleCredential:  &userInfo.Id,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	_, exception = s.userSettingRepository.CreateOneByUserId(
		*newUserId,
		inputs.CreateUserSettingInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	exception = userdata.NewUserDataCacheClient().Set(
		newUser.Name,
		userdata.UserDataCache{
			Id:                 *newUserId,
			PublicId:           newUser.PublicId,
			Name:               newUser.Name,
			DisplayName:        newUser.DisplayName,
			Email:              newUser.Email,
			AccessToken:        *newAccessToken,
			CSRFToken:          *newCSRFToken,
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
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	// send the welcome email to the registered user
	if exception = s.emailClient.SendWelcomeEmail(ctx, emaildto.SendWelcomeEmailRequestDto{
		To:       newUser.Email,
		UserName: newUser.Name,
		Status:   newUser.Status.String(),
	}); exception != nil {
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
	}

	return &authdto.RegisterViaGoogleResponseDto{
		PublicId:     newUser.PublicId,
		Name:         newUser.Name,
		DisplayName:  newUser.DisplayName,
		Email:        newUser.Email,
		AccessToken:  *newAccessToken,
		RefreshToken: *newRefreshToken,
		CSRFToken:    *newCSRFToken,
		CreatedAt:    newUser.CreatedAt,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context, reqDto *authdto.LoginRequestDto,
) (*authdto.LoginResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	if stringutil.IsAlphaAndNumberString(reqDto.Body.Account) { // if the account field contains user name
		if user, exception = s.userRepository.GetOneByName(
			reqDto.Body.Account,
			nil,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	} else if stringutil.IsEmailString(reqDto.Body.Account) { // if the account field contains email
		if user, exception = s.userRepository.GetOneByEmail(
			reqDto.Body.Account,
			nil,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	} else {
		tx.Rollback()
		return nil, apiexceptions.Auth.InvalidDto()
	}

	if user == nil {
		tx.Rollback()
		return nil, apiexceptions.Auth.InvalidDto()
	}

	if user.BlockLoginUntil.After(time.Now()) {
		tx.Rollback()
		return nil, apiexceptions.Auth.LoginBlockedDueToTryingTooManyTimes(user.BlockLoginUntil)
	}

	if !s.comparePassword(user.Password, reqDto.Body.Password) {
		newLoginCount := user.LoginCount + 1
		blockLoginUntil, exception := s.getLoginBlockedUntilByLoginCount(newLoginCount)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

		_, exception = s.userRepository.UpdateOneById(
			user.Id,
			inputs.PartialUpdateUserInput{
				Values: inputs.UpdateUserInput{
					LoginCount:     &newLoginCount,
					BlockLoginUtil: blockLoginUntil,
				},
				SetNull: nil,
			},
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

		if blockLoginUntil != nil {
			tx.Rollback()
			return nil, apiexceptions.Auth.LoginBlockedDueToTryingTooManyTimes(*blockLoginUntil)
		}

		tx.Rollback()
		return nil, apiexceptions.Auth.WrongPassword() // login procedure early ends here
	}

	if user.UserAgent != reqDto.Header.UserAgent {
		// send a security email to warn the user
		if exception := s.emailClient.SendSecurityAlertEmail(ctx, emaildto.SendSecurityAlertEmailRequestDto{
			To:               user.Email,
			UserName:         user.Name,
			Status:           user.Status.String(),
			AlertType:        "Login in Different Place",
			Reason:           "Your account has a recent login action in other place",
			TimeOfOccurrence: time.Now(),
			OtherDetails:     "",
		}); exception != nil {
			_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		}
	}

	newAccessToken, exception := s.generateAccessToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newRefreshToken, exception := s.generateRefreshToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newCSRFToken, exception := s.generateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// check if the user data cache exists
	if _, exception := userdata.NewUserDataCacheClient().Get(user.Name); exception == nil {
		// then just update the existing user data cache
		if exception = userdata.NewUserDataCacheClient().Update(
			user.Name,
			cacheinputs.UpdateUserDataCacheInput{
				AccessToken: newAccessToken,
				CSRFToken:   newCSRFToken,
			},
		); exception != nil {
			_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		}
	} else { // else if it does not exist
		// then we have to first get the relative data from different tables
		// we done this by one custom sql so it's not that slow...
		// once we have the required data, we set it as the user data cache
		output := struct {
			Id                 uuid.UUID        `gorm:"id"`
			PublicId           uuid.UUID        `gorm:"public_id"`
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
		err := tx.Raw(usersql.GetUserDataCacheByIdSQL, user.Id).
			Row().
			Scan(
				&output.Id,
				&output.PublicId,
				&output.Name,
				&output.DisplayName,
				&output.Email,
				&output.Role,
				&output.Plan,
				&output.Status,
				&output.AvatarURL,
				&output.Language,
				&output.GeneralSettingCode,
				&output.PrivacySettingCode,
				&output.CreatedAt,
				&output.UpdatedAt,
			)
		if err != nil {
			tx.Rollback()
			return nil, apiexceptions.User.NotFound().WithOrigin(err)
		}

		newUserDataCache := userdata.UserDataCache{
			Id:                 user.Id,
			PublicId:           output.PublicId,
			Name:               output.Name,
			DisplayName:        output.DisplayName,
			Email:              output.Email,
			AccessToken:        *newAccessToken,
			CSRFToken:          *newCSRFToken,
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
		exception := userdata.NewUserDataCacheClient().Set(
			user.Name,
			newUserDataCache,
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := s.userRepository.UpdateOneById(
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &user.PrevStatus,
				RefreshToken: newRefreshToken,
				UserAgent:    &reqDto.Header.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	return &authdto.LoginResponseDto{
		PublicId:     user.PublicId,
		Name:         user.Name,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		AccessToken:  *newAccessToken,
		RefreshToken: updatedUser.RefreshToken,
		CSRFToken:    *newCSRFToken,
		UpdatedAt:    updatedUser.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (s *AuthService) LoginViaGoogle(
	ctx context.Context, reqDto *authdto.LoginViaGoogleRequestDto,
) (*authdto.LoginViaGoogleResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.Auth.InvalidDto().WithOrigin(err)
	}

	userInfo, exception := s.oauthService.GetGoogleUserInfo(ctx, reqDto.Body.AuthorizationCode)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	// otherwise, the user should provide their account and password
	var user *schemas.User = nil
	if user, exception = s.userRepository.GetOneByEmail(
		userInfo.Email,
		[]schemas.UserRelation{schemas.UserRelation_UserAccount},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if user == nil {
		tx.Rollback()
		return nil, apiexceptions.Auth.InvalidDto()
	}

	if user.BlockLoginUntil.After(time.Now()) {
		tx.Rollback()
		return nil, apiexceptions.Auth.LoginBlockedDueToTryingTooManyTimes(user.BlockLoginUntil)
	}

	if user.UserAccount.GoogleCredential == nil || userInfo.Id != *user.UserAccount.GoogleCredential {
		newLoginCount := user.LoginCount + 1
		blockLoginUntil, exception := s.getLoginBlockedUntilByLoginCount(newLoginCount)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

		_, exception = s.userRepository.UpdateOneById(
			user.Id,
			inputs.PartialUpdateUserInput{
				Values: inputs.UpdateUserInput{
					LoginCount:     &newLoginCount,
					BlockLoginUtil: blockLoginUntil,
				},
				SetNull: nil,
			},
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}

		if blockLoginUntil != nil {
			tx.Rollback()
			return nil, apiexceptions.Auth.LoginBlockedDueToTryingTooManyTimes(*blockLoginUntil)
		}

		tx.Rollback()
		return nil, apiexceptions.Auth.WrongPassword() // login via google procedure early ends here
	}

	if user.UserAgent != reqDto.Header.UserAgent {
		// send a security email to warn the user
		if exception := s.emailClient.SendSecurityAlertEmail(ctx, emaildto.SendSecurityAlertEmailRequestDto{
			To:               user.Email,
			UserName:         user.Name,
			Status:           user.Status.String(),
			AlertType:        "Login in Different Place",
			Reason:           "Your account has a recent login action in other place",
			TimeOfOccurrence: time.Now(),
			OtherDetails:     "",
		}); exception != nil {
			_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		}
	}

	newAccessToken, exception := s.generateAccessToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newRefreshToken, exception := s.generateRefreshToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newCSRFToken, exception := s.generateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// check if the user data cache exists
	if _, exception := userdata.NewUserDataCacheClient().Get(user.Name); exception == nil {
		// then just update the existing user data cache
		if exception = userdata.NewUserDataCacheClient().Update(
			user.Name,
			cacheinputs.UpdateUserDataCacheInput{
				AccessToken: newAccessToken,
				CSRFToken:   newCSRFToken,
			},
		); exception != nil {
			_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		}
	} else { // else if it does not exist
		// then we have to first get the relative data from different tables
		// we done this by one custom sql so it's not that slow...
		// once we have the required data, we set it as the user data cache
		output := struct {
			Id                 uuid.UUID        `gorm:"id"`
			PublicId           uuid.UUID        `gorm:"public_id"`
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
		err := tx.Raw(usersql.GetUserDataCacheByIdSQL, user.Id).
			Row().
			Scan(
				&output.Id,
				&output.PublicId,
				&output.Name,
				&output.DisplayName,
				&output.Email,
				&output.Role,
				&output.Plan,
				&output.Status,
				&output.AvatarURL,
				&output.Language,
				&output.GeneralSettingCode,
				&output.PrivacySettingCode,
				&output.CreatedAt,
				&output.UpdatedAt,
			)
		if err != nil {
			tx.Rollback()
			return nil, apiexceptions.User.NotFound().WithOrigin(err)
		}

		newUserDataCache := userdata.UserDataCache{
			Id:                 user.Id,
			PublicId:           output.PublicId,
			Name:               output.Name,
			DisplayName:        output.DisplayName,
			Email:              output.Email,
			AccessToken:        *newAccessToken,
			CSRFToken:          *newCSRFToken,
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
		exception := userdata.NewUserDataCacheClient().Set(
			user.Name,
			newUserDataCache,
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := s.userRepository.UpdateOneById(
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &user.PrevStatus,
				RefreshToken: newRefreshToken,
				UserAgent:    &reqDto.Header.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	return &authdto.LoginViaGoogleResponseDto{
		PublicId:     user.PublicId,
		Name:         user.Name,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		AccessToken:  *newAccessToken,
		RefreshToken: updatedUser.RefreshToken,
		CSRFToken:    *newCSRFToken,
		UpdatedAt:    updatedUser.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context, reqDto *authdto.LogoutRequestDto,
) (*authdto.LogoutResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.Auth.InvalidDto().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserName, exception := contexts.GetActorUserName(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	offlineStatus := enums.UserStatus_Offline
	emptyString := ""
	updatedUser, exception := s.userRepository.UpdateOneById(
		actorUserId,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Status:       &offlineStatus,
				RefreshToken: &emptyString,
			},
			SetNull: nil,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	exception = userdata.NewUserDataCacheClient().Delete(actorUserName)
	if exception != nil {
		return nil, exception
	}

	return &authdto.LogoutResponseDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) SendAuthCode(
	ctx context.Context, reqDto *authdto.SendAuthCodeRequestDto,
) (*authdto.SendAuthCodeResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	authCode := s.generateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	blockAuthCodeUntil := s.getAuthCodeBlockUntil()
	output := struct {
		Name               string    `json:"name" gorm:"column:name;"`
		UserAgent          string    `json:"userAgent" gorm:"column:user_agent;"`
		BlockAuthCodeUntil time.Time `json:"blockAuthCodeUntil" gorm:"column:block_auth_code_until;"`
		Now                time.Time `json:"now" gorm:"column:now;"`
	}{}
	err := db.Raw(authsql.UpdateAuthCodeSQL,
		authCode, authCodeExpiredAt, blockAuthCodeUntil, reqDto.Body.Email,
	).Row().
		Scan(&output.Name, &output.UserAgent, &output.BlockAuthCodeUntil, &output.Now)
	if err != nil {
		return nil, apiexceptions.Auth.AuthCodeBlockedDueToTryingTooManyTimes(output.BlockAuthCodeUntil).WithOrigin(err)
	}

	if exception := s.emailClient.SendValidationEmail(ctx, emaildto.SendValidationEmailRequestDto{
		To:        reqDto.Body.Email,
		UserName:  output.Name,
		AuthCode:  authCode,
		UserAgent: output.UserAgent,
		ExpiredAt: authCodeExpiredAt,
	}); exception != nil {
		return nil, exception
	}

	return &authdto.SendAuthCodeResponseDto{
		AuthCodeExpiredAt:  authCodeExpiredAt,
		BlockAuthCodeUntil: blockAuthCodeUntil,
		UpdatedAt:          time.Now(),
	}, nil
}

func (s *AuthService) ValidateEmail(
	ctx context.Context, reqDto *authdto.ValidateEmailRequestDto,
) (*authdto.ValidateEmailResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	var updatedAt time.Time
	err := db.Raw(authsql.ValidateEmailSQL, actorUserId, reqDto.Body.AuthCode).
		Row().
		Scan(&updatedAt)
	if err != nil {
		return nil, apiexceptions.User.FailedToUpdate().WithOrigin(err)
	}

	return &authdto.ValidateEmailResponseDto{
		UpdatedAt: updatedAt,
	}, nil
}

func (s *AuthService) ResetEmail(
	ctx context.Context, reqDto *authdto.ResetEmailRequestDto,
) (*authdto.ResetEmailResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	var updatedAt time.Time
	err := tx.Raw(authsql.ResetEmailSQL, reqDto.Body.NewEmail, reqDto.Body.AuthCode, actorUserId).
		Row().
		Scan(&updatedAt)
	if err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToUpdate().WithOrigin(err)
	}

	authCode := s.generateAuthCode()
	authCodeExpiredAt := time.Now().Add(constants.ExpirationTimeOfAuthCode)
	_, exception = s.userAccountRepository.UpdateOneByUserId(
		actorUserId,
		inputs.PartialUpdateUserAccountInput{
			Values: inputs.UpdateUserAccountInput{
				AuthCode:          &authCode,
				AuthCodeExpiredAt: &authCodeExpiredAt,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	return &authdto.ResetEmailResponseDto{
		UpdatedAt: updatedAt,
	}, nil
}

func (s *AuthService) ForgetPassword(
	ctx context.Context, reqDto *authdto.ForgetPasswordRequestDto,
) (*authdto.ForgetPasswordResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	var user *schemas.User = nil
	var exception *exceptions.Exception = nil
	var preloads = []schemas.UserRelation{schemas.UserRelation_UserAccount, schemas.UserRelation_UserInfo, schemas.UserRelation_UserSetting}
	if stringutil.IsEmailString(reqDto.Body.Account) { // if the account field contains email
		if user, exception = s.userRepository.GetOneByEmail(
			reqDto.Body.Account,
			preloads,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	} else if stringutil.IsAlphaAndNumberString(reqDto.Body.Account) { // if the account field contains user name
		if user, exception = s.userRepository.GetOneByName(
			reqDto.Body.Account,
			preloads,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	} else {
		tx.Rollback()
		return nil, apiexceptions.Auth.InvalidDto()
	}

	if reqDto.Body.AuthCode != user.UserAccount.AuthCode {
		tx.Rollback()
		return nil, apiexceptions.Auth.WrongAuthCode()
	}

	newAccessToken, exception := s.generateAccessToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newRefreshToken, exception := s.generateRefreshToken(user.PublicId, user.Name, user.Email, user.UserAgent)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	newCSRFToken, exception := s.generateCSRFToken()
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// update the access token of the user
	exception = userdata.NewUserDataCacheClient().Update(user.Name, cacheinputs.UpdateUserDataCacheInput{AccessToken: newAccessToken})
	if exception != nil {
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		// and also try to set the new user cache data
		exception = userdata.NewUserDataCacheClient().Set(user.Name, userdata.UserDataCache{
			Id:                 user.Id,
			PublicId:           user.PublicId,
			Name:               user.Name,
			DisplayName:        user.DisplayName,
			Email:              user.Email,
			AccessToken:        *newAccessToken,
			CSRFToken:          *newCSRFToken,
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
			_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
		}
	}

	hashedPassword, exception := s.hashPassword(reqDto.Body.NewPassword)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// update the refresh token and the status of the user
	var zeroLoginCount int32 = 0 // reset the login count if the login procedure is valid
	updatedUser, exception := s.userRepository.UpdateOneById(
		user.Id,
		inputs.PartialUpdateUserInput{
			Values: inputs.UpdateUserInput{
				Password:     &hashedPassword,
				RefreshToken: newRefreshToken,
				UserAgent:    &reqDto.Header.UserAgent,
				LoginCount:   &zeroLoginCount,
			},
			SetNull: nil,
		},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithOrigin(err)
	}

	return &authdto.ForgetPasswordResponseDto{
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *AuthService) ResetMe(
	ctx context.Context, reqDto *authdto.ResetMeRequestDto,
) (*authdto.ResetMeResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	// Instead of deleting the user, we recreate their relative data in the database
	// and make sure not to update the access token and refresh token, and csrf token in the reset logic
	// Note that the user will not logged out after the reset operation

	// try to retrieve the target user to reset and validate his/her auth code first
	var resetUserAccount schemas.UserAccount
	result := tx.Model(&resetUserAccount).
		Where("user_id = ? AND auth_code = ?", actorUserId, reqDto.Body.AuthCode).
		First(&resetUserAccount)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.UserAccount.NotFound().WithOrigin(err)
	}

	// delete the user info
	if err := tx.Where("user_id = ?", actorUserId).Delete(&schemas.UserInfo{}).Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.UserInfo.FailedToDelete().WithOrigin(err)
	}
	// and then re-create a new user info
	if _, exception := s.userInfoRepository.CreateOneByUserId(
		resetUserAccount.UserId,
		inputs.CreateUserInfoInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// delete the user setting
	if err := tx.Where("user_id = ?", actorUserId).Delete(&schemas.UserSetting{}).Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.UserSetting.FailedToDelete().WithOrigin(err)
	}
	// and then re-create a new user setting
	if _, exception := s.userSettingRepository.CreateOneByUserId(
		resetUserAccount.UserId,
		inputs.CreateUserSettingInput{},
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// delete all the badges of the user
	if err := tx.Exec(badgesql.DeleteAllMyBadgesSQL, actorUserId).Error; err != nil {
		// skip if there's no users to badges to delete
	}

	// soft delete all the root shelves of the user
	if exception := s.rootShelfRepository.SoftDeleteManyByUserId(
		actorUserId,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	); exception != nil {
		// skip if there's no root shelves to soft delete
	} else {
		// then hard delete all the root shelves of the user
		if exception := s.rootShelfRepository.HardDeleteManyByUserId(
			actorUserId,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		); exception != nil {
			// skip if there's no root shelves to hard delete
		}
	}

	// delete other stuff in the future...

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apiexceptions.User.FailedToCommitTransaction().WithDetails(err)
	}

	return &authdto.ResetMeResponseDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *AuthService) DeleteMe(
	ctx context.Context, reqDto *authdto.DeleteMeRequestDto,
) (*authdto.DeleteMeResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, apiexceptions.User.InvalidInput().WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserName, exception := contexts.GetActorUserName(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	if err := db.Exec(authsql.DeleteMeSQL, actorUserId, reqDto.Body.AuthCode).Error; err != nil {
		return nil, apiexceptions.User.FailedToDelete().WithOrigin(err)
	}

	exception = userdata.NewUserDataCacheClient().Delete(actorUserName)
	if exception != nil {
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, exception.String(), exception)
	}

	return &authdto.DeleteMeResponseDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *AuthService) RegisterViaMeta() {}

func (s *AuthService) RegisterViaGithub() {}

func (s *AuthService) LoginViaMeta() {}

func (s *AuthService) LoginViaGithub() {}
