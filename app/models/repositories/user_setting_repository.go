package repositories

import (
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type UserSettingRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserSetting, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput, opts ...options.RepositoryOptions) (*schemas.UserSetting, *exceptions.Exception)
}

type UserSettingRepository struct{}

func NewUserSettingRepository() UserSettingRepositoryInterface {
	return &UserSettingRepository{}
}

func (r *UserSettingRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserSetting, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var userSetting schemas.UserSetting
	result := parsedOptions.DB.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&userSetting)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserSetting.NotFound().WithOrigin(result.Error)},
		{First: userSetting.Id == uuid.Nil, Second: exceptions.UserSetting.NotFound()},
	}); exception != nil {
		return nil, exception
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
		return nil, exceptions.UserSetting.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserSetting{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserSetting)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserSetting.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UserSetting.NoChanges()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserSetting.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UserSetting.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}
