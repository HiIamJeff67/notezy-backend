package services

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	usersettingsdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/user-settings"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
)

type UserSettingServiceInterface interface {
	GetMySetting(ctx context.Context, requestDto *usersettingsdto.GetMySettingRequestDto) (*usersettingsdto.GetMySettingResponseDto, *exceptions.Exception)
	UpdateMySetting(ctx context.Context, requestDto *usersettingsdto.UpdateMySettingRequestDto) (*usersettingsdto.UpdateMySettingResponseDto, *exceptions.Exception)
}

type UserSettingService struct {
	db                    *gorm.DB
	userSettingRepository repositories.UserSettingRepositoryInterface
}

func NewUserSettingService(
	db *gorm.DB,
	userSettingRepository repositories.UserSettingRepositoryInterface,
) UserSettingServiceInterface {
	if db == nil {
		db = data.NotezyDB
	}
	return &UserSettingService{
		db:                    db,
		userSettingRepository: userSettingRepository,
	}
}

/* ============================== Service Methods for UserSetting ============================== */

func (s *UserSettingService) GetMySetting(
	ctx context.Context,
	requestDto *usersettingsdto.GetMySettingRequestDto,
) (*usersettingsdto.GetMySettingResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"UserSetting",
			"GetMySetting",
			"User setting request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	userSetting, exception := s.userSettingRepository.GetOneByUserId(
		actorUserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &usersettingsdto.GetMySettingResponseDto{
		Language:           userSetting.Language,
		GeneralSettingCode: userSetting.GeneralSettingCode,
		PrivacySettingCode: userSetting.PrivacySettingCode,
	}, nil
}

func (s *UserSettingService) UpdateMySetting(
	ctx context.Context,
	requestDto *usersettingsdto.UpdateMySettingRequestDto,
) (*usersettingsdto.UpdateMySettingResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"UserSetting",
			"UpdateMySetting",
			"User setting request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	db := s.db.WithContext(ctx)

	updatedUserSetting, exception := s.userSettingRepository.UpdateOneByUserId(
		actorUserId,
		inputs.PartialUpdateUserSettingInput{
			Values: inputs.UpdateUserSettingInput{
				Language:           requestDto.Body.Values.Language,
				GeneralSettingCode: requestDto.Body.Values.GeneralSettingCode,
				PrivacySettingCode: requestDto.Body.Values.PrivacySettingCode,
			},
			SetNull: requestDto.Body.SetNull,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &usersettingsdto.UpdateMySettingResponseDto{
		UpdatedAt: updatedUserSetting.UpdatedAt,
	}, nil
}
