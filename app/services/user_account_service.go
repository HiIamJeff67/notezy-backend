package services

import (
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
)

/* ============================== Interface & Instance ============================== */

type UserAccountServiceInterface interface {
	GetMyAccount(reqDto *dtos.GetMyAccountReqDto) (*dtos.GetMyAccountResDto, *exceptions.Exception)
	UpdateMyAccount(reqDto *dtos.UpdateMyAccountReqDto) (*dtos.UpdateMyAccountResDto, *exceptions.Exception)
}

type userAccountService struct{}

var UserAccountService UserAccountServiceInterface = &userAccountService{}

/* ============================== Services ============================== */

func (s *userAccountService) GetMyAccount(reqDto *dtos.GetMyAccountReqDto) (*dtos.GetMyAccountResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userAccountRepository := repositories.NewUserAccountRepository(nil)

	userAccount, exception := userAccountRepository.GetOneByUserId(reqDto.UserId)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyAccountResDto{
		CountryCode:       userAccount.CountryCode,
		PhoneNumber:       userAccount.PhoneNumber,
		GoogleCredential:  userAccount.GoogleCredential,
		DiscordCredential: userAccount.DiscordCredential,
	}, nil
}

func (s *userAccountService) UpdateMyAccount(reqDto *dtos.UpdateMyAccountReqDto) (*dtos.UpdateMyAccountResDto, *exceptions.Exception) {
	if err := models.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	userAccountRepository := repositories.NewUserAccountRepository(nil)

	updatedUserAccount, exception := userAccountRepository.UpdateOneByUserId(reqDto.UserId, inputs.PartialUpdateUserAccountInput{
		Values: inputs.UpdateUserAccountInput{
			CountryCode:       reqDto.Values.CountryCode,
			PhoneNumber:       reqDto.Values.PhoneNumber,
			GoogleCredential:  reqDto.Values.GoogleCredential,
			DiscordCredential: reqDto.Values.DiscordCredential,
		},
		SetNull: reqDto.SetNull,
	})
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyAccountResDto{
		UpdatedAt: updatedUserAccount.UpdatedAt,
	}, nil
}
