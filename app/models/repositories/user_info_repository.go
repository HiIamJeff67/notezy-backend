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

type UserInfoRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserInfo, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInfoInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserInfoInput, opts ...options.RepositoryOptions) (*schemas.UserInfo, *exceptions.Exception)
}

type UserInfoRepository struct{}

func NewUserInfoRepository() UserInfoRepositoryInterface {
	return &UserInfoRepository{}
}

/* ============================== Implementations ============================== */

func (r *UserInfoRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserInfo, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	userInfo := schemas.UserInfo{}
	result := parsedOptions.DB.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func (r *UserInfoRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateUserInfoInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUserInfo schemas.UserInfo
	newUserInfo.UserId = userId
	if err := copier.Copy(&newUserInfo, &input); err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserInfo{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return &newUserInfo.Id, nil
}

func (r *UserInfoRepository) UpdateOneByUserId(
	userId uuid.UUID,
	input inputs.PartialUpdateUserInfoInput,
	opts ...options.RepositoryOptions,
) (*schemas.UserInfo, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUserInfo, exception := r.GetOneByUserId(
		userId,
		opts...,
	)
	if exception != nil || existingUserInfo == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserInfo)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserInfo)
	}

	result := parsedOptions.DB.Model(&schemas.UserInfo{}).
		Where("user_id = ?", userId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserInfo.NoChanges()
	}

	return &updates, nil
}
