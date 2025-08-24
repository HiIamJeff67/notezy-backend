package services

import (
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	validation "notezy-backend/app/validation"
)

/* ============================== Interface & Instance ============================== */

type UserAccountServiceInterface interface {
	GetMyAccount(reqDto *dtos.GetMyAccountReqDto) (*dtos.GetMyAccountResDto, *exceptions.Exception)
	UpdateMyAccount(reqDto *dtos.UpdateMyAccountReqDto) (*dtos.UpdateMyAccountResDto, *exceptions.Exception)
}

type UserAccountService struct {
	db *gorm.DB
}

func NewUserAccountService(db *gorm.DB) UserAccountServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserAccountService{db: db}
}

/* ============================== Service Methods for UserAccount ============================== */

func (s *UserAccountService) GetMyAccount(reqDto *dtos.GetMyAccountReqDto) (*dtos.GetMyAccountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userAccountRepository := repositories.NewUserAccountRepository(s.db)

	userAccount, exception := userAccountRepository.GetOneByUserId(reqDto.ContextFields.UserId)
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

func (s *UserAccountService) UpdateMyAccount(reqDto *dtos.UpdateMyAccountReqDto) (*dtos.UpdateMyAccountResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userAccountRepository := repositories.NewUserAccountRepository(s.db)

	updatedUserAccount, exception := userAccountRepository.UpdateOneByUserId(reqDto.ContextFields.UserId, inputs.PartialUpdateUserAccountInput{
		Values: inputs.UpdateUserAccountInput{
			CountryCode:       reqDto.Body.Values.CountryCode,
			PhoneNumber:       reqDto.Body.Values.PhoneNumber,
			GoogleCredential:  reqDto.Body.Values.GoogleCredential,
			DiscordCredential: reqDto.Body.Values.DiscordCredential,
		},
		SetNull: reqDto.Body.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyAccountResDto{
		UpdatedAt: updatedUserAccount.UpdatedAt,
	}, nil
}
