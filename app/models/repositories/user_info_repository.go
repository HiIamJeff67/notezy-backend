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

type UserInfoRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserInfo, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInfoInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserInfoInput, opts ...options.RepositoryOptions) (*schemas.UserInfo, *exceptions.Exception)
}

type UserInfoRepository struct{}

func NewUserInfoRepository() UserInfoRepositoryInterface {
	return &UserInfoRepository{}
}

func (r *UserInfoRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserInfo, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	userInfo := schemas.UserInfo{}
	result := parsedOptions.DB.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&userInfo)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserInfo.NotFound().WithOrigin(result.Error)},
		{First: userInfo.Id == uuid.Nil, Second: exceptions.UserInfo.NotFound()},
	}); exception != nil {
		return nil, exception
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
		return nil, exceptions.UserInfo.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserInfo{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserInfo)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserInfo.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UserInfo.NoChanges()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UserInfo.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UserInfo.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}
