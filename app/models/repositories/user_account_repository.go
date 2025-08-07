package repositories

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type UserAccountRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserAccountInput) (*schemas.UserAccount, *exceptions.Exception)
}

type UserAccountRepository struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) UserAccountRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserAccountRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *UserAccountRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception) {
	userAccount := schemas.UserAccount{}
	result := r.db.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}

	return &userAccount, nil
}

func (r *UserAccountRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput) (*uuid.UUID, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserAccount.InvalidInput().WithError(err).Log()
	}

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId
	if err := copier.Copy(&newUserAccount, &input); err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.UserAccount{}.TableName()).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}
	return &newUserAccount.Id, nil
}

func (r *UserAccountRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserAccountInput) (*schemas.UserAccount, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserAccount.InvalidInput().WithError(err).Log()
	}

	existingUserAccount, exception := r.GetOneByUserId(userId)
	if exception != nil || existingUserAccount == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserAccount)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserAccount)
	}

	result := r.db.Table(schemas.UserAccount{}.TableName()).
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
