package operations

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/schemas"
	"notezy-backend/app/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserAccountByUserId(db *gorm.DB, userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception) {
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

func GetAllUserAccount(db *gorm.DB) (*[]schemas.UserAccount, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	userAccounts := []schemas.UserAccount{}
	result := db.Table(schemas.UserAccount{}.TableName()).Find(&userAccounts)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}
	return &userAccounts, nil
}

func CreateUserAccountByUserId(db *gorm.DB, userId uuid.UUID, input inputs.CreateUserAccountInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId
	util.CopyNonNilFields(&newUserAccount, input)
	result := db.Table(schemas.UserAccount{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithError(err)
	}
	return &newUserAccount.Id, nil
}

func UpdateUserAccountByUserId(db *gorm.DB, userId uuid.UUID, input inputs.UpdateUserAccountInput) (*schemas.UserAccount, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var updatedUserAccount schemas.UserAccount
	util.CopyNonNilFields(&updatedUserAccount, input)
	result := db.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Updates(&input)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToUpdate().WithError(err)
	}

	return &updatedUserAccount, nil
}

// We do not allow to just delete the userAccount,
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
