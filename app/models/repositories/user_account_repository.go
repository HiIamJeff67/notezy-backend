package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type UserAccountRepository interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput) *exceptions.Exception
	UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserAccountInput) *exceptions.Exception
}

type userAccountRepository struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) *userAccountRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &userAccountRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *userAccountRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserAccount, *exceptions.Exception) {
	userAccount := schemas.UserAccount{}
	result := r.db.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithError(err)
	}

	return &userAccount, nil
}

func (r *userAccountRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput) *exceptions.Exception {
	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId
	util.CopyNonNilFields(&newUserAccount, input)
	result := r.db.Table(schemas.UserAccount{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return exceptions.UserAccount.FailedToCreate().WithError(err)
	}
	return nil
}

func (r *userAccountRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserAccountInput) *exceptions.Exception {
	var updatedUserAccount schemas.UserAccount
	util.CopyNonNilFields(&updatedUserAccount, input)
	result := r.db.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Updates(&input)
	if err := result.Error; err != nil {
		return exceptions.UserAccount.FailedToUpdate().WithError(err)
	}

	return nil
}

// We do not allow to just delete the userAccount,
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
