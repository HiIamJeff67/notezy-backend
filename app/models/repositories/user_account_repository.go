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

type UserAccountRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserAccountInput, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
}

type UserAccountRepository struct{}

func NewUserAccountRepository() UserAccountRepositoryInterface {
	return &UserAccountRepository{}
}

/* ============================== Implementations ============================== */

func (r *UserAccountRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var userAccount schemas.UserAccount
	result := parsedOptions.DB.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}

	return &userAccount, nil
}

func (r *UserAccountRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId

	if err := copier.Copy(&newUserAccount, &input); err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}
	return &newUserAccount.Id, nil
}

func (r *UserAccountRepository) UpdateOneByUserId(
	userId uuid.UUID,
	input inputs.PartialUpdateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUserAccount, exception := r.GetOneByUserId(
		userId,
		opts...,
	)
	if exception != nil || existingUserAccount == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserAccount)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserAccount)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where("user_id = ?", userId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserAccount.NoChanges()
	}

	return &updates, nil
}

// We do not allow to just delete the userAccount,
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
