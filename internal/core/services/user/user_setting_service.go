package user

import (
	"context"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-settings"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/core/data/database"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/repositories"
	enums "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas/enums"
)

type UserSettingServiceInterface interface {
	GetMySetting(ctx context.Context, requestDto *apicontract.GetMySettingRequestDto) (*apicontract.GetMySettingResponseDto, *exceptions.Exception)
	UpdateMySetting(ctx context.Context, requestDto *apicontract.UpdateMySettingRequestDto) (*apicontract.UpdateMySettingResponseDto, *exceptions.Exception)
}

type UserSettingService struct {
	validator             *validator.Validate
	db                    *gorm.DB
	userSettingRepository repositories.UserSettingRepositoryInterface
}

func NewUserSettingService(
	validator *validator.Validate,
	db *gorm.DB,
	userSettingRepository repositories.UserSettingRepositoryInterface,
) UserSettingServiceInterface {
	if db == nil {
		db = data.NotezyDB
	}
	return &UserSettingService{
		validator:             validator,
		db:                    db,
		userSettingRepository: userSettingRepository,
	}
}

/* ============================== Service Methods for UserSetting ============================== */

func (s *UserSettingService) GetMySetting(
	ctx context.Context,
	requestDto *apicontract.GetMySettingRequestDto,
) (*apicontract.GetMySettingResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
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

	return &apicontract.GetMySettingResponseDto{
		Language:           *userSetting.Language.ToContractable(),
		GeneralSettingCode: userSetting.GeneralSettingCode,
		PrivacySettingCode: userSetting.PrivacySettingCode,
	}, nil
}

func (s *UserSettingService) UpdateMySetting(
	ctx context.Context,
	requestDto *apicontract.UpdateMySettingRequestDto,
) (*apicontract.UpdateMySettingResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
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
	var language *enums.Language
	if requestDto.Body.Values.Language != nil {
		value := enums.Language(*requestDto.Body.Values.Language)
		language = &value
	}

	updatedUserSetting, exception := s.userSettingRepository.UpdateOneByUserId(
		actorUserId,
		inputs.PartialUpdateUserSettingInput{
			Values: inputs.UpdateUserSettingInput{
				Language:           language,
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

	return &apicontract.UpdateMySettingResponseDto{
		UpdatedAt: updatedUserSetting.UpdatedAt,
	}, nil
}
