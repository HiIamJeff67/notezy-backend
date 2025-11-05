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
	GetOneByUserId(db *gorm.DB, userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception)
	CreateOneByUserId(db *gorm.DB, userId uuid.UUID, input inputs.CreateUserAccountInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(db *gorm.DB, userId uuid.UUID, input inputs.PartialUpdateUserAccountInput) (*schemas.UserAccount, *exceptions.Exception)
}

type UserAccountRepository struct{}

func NewUserAccountRepository() UserAccountRepositoryInterface {
	return &UserAccountRepository{}
}

/* ============================== CRUD operations ============================== */

func (r *UserAccountRepository) GetOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
) (*schemas.UserAccount, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	userAccount := schemas.UserAccount{}

	result := db.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}

	return &userAccount, nil
}

func (r *UserAccountRepository) CreateOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
	input inputs.CreateUserAccountInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId

	if err := copier.Copy(&newUserAccount, &input); err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}

	result := db.Model(&schemas.UserAccount{}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}
	return &newUserAccount.Id, nil
}

func (r *UserAccountRepository) UpdateOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
	input inputs.PartialUpdateUserAccountInput,
) (*schemas.UserAccount, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	existingUserAccount, exception := r.GetOneByUserId(db, userId)
	if exception != nil || existingUserAccount == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserAccount)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserAccount)
	}

	result := db.Model(&schemas.UserAccount{}).
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
