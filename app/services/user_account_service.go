package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
)

/* ============================== Interface & Instance ============================== */

type UserAccountServiceInterface interface {
	GetMyAccount(ctx context.Context, reqDto *dtos.GetMyAccountReqDto) (*dtos.GetMyAccountResDto, *exceptions.Exception)
	UpdateMyAccount(ctx context.Context, reqDto *dtos.UpdateMyAccountReqDto) (*dtos.UpdateMyAccountResDto, *exceptions.Exception)
}

type UserAccountService struct {
	db                    *gorm.DB
	userAccountRepository repositories.UserAccountRepositoryInterface
}

func NewUserAccountService(
	db *gorm.DB,
	userAccountRepository repositories.UserAccountRepositoryInterface,
) UserAccountServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserAccountService{
		db:                    db,
		userAccountRepository: userAccountRepository,
	}
}

/* ============================== Service Methods for UserAccount ============================== */

func (s *UserAccountService) GetMyAccount(
	ctx context.Context, reqDto *dtos.GetMyAccountReqDto,
) (*dtos.GetMyAccountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.UserAccount.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	userAccount, exception := s.userAccountRepository.GetOneByUserId(
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyAccountResDto{
		CountryCode:       *userAccount.CountryCode,
		PhoneNumber:       userAccount.PhoneNumber,
		GoogleCredential:  userAccount.GoogleCredential,
		DiscordCredential: userAccount.DiscordCredential,
	}, nil
}

func (s *UserAccountService) UpdateMyAccount(
	ctx context.Context, reqDto *dtos.UpdateMyAccountReqDto,
) (*dtos.UpdateMyAccountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.UserAccount.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	result := db.Model(&schemas.UserAccount{}).
		Where("user_id = ? AND auth_code = ?", reqDto.ContextFields.UserId, reqDto.Body.AuthCode).
		First(&schemas.UserAccount{})
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}

	_, exception := s.userAccountRepository.UpdateOneByUserId(
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateUserAccountInput{
			Values: inputs.UpdateUserAccountInput{
				BackupEmail: reqDto.Body.Values.BackupEmail,
				CountryCode: reqDto.Body.Values.CountryCode,
				PhoneNumber: reqDto.Body.Values.PhoneNumber,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyAccountResDto{
		UpdatedAt: time.Now(),
	}, nil
}
