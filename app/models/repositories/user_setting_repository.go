package repositories

import (
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type UserSettingRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserSetting, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput, opts ...options.RepositoryOptions) (*schemas.UserSetting, *exceptions.Exception)
}

type UserSettingRepository struct{}

func NewUserSettingRepository() UserSettingRepositoryInterface {
	return &UserSettingRepository{}
}

/* ============================== Implementations ============================== */

func (r *UserSettingRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserSetting, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var userSetting schemas.UserSetting
	result := parsedOptions.DB.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func (r *UserSettingRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateUserSettingInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUserSetting schemas.UserSetting
	newUserSetting.UserId = userId
	if err := copier.Copy(&newUserSetting, &input); err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserSetting{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return &newUserSetting.Id, nil
}

func (r *UserSettingRepository) UpdateOneByUserId(
	userId uuid.UUID,
	input inputs.PartialUpdateUserSettingInput,
	opts ...options.RepositoryOptions,
) (*schemas.UserSetting, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUserSetting, exception := r.GetOneByUserId(
		userId,
		opts...,
	)
	if exception != nil || existingUserSetting == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserSetting)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserSetting)
	}

	result := parsedOptions.DB.Model(&schemas.UserSetting{}).
		Where("user_id = ?").
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserSetting.NoChanges()
	}

	return &updates, nil
}
