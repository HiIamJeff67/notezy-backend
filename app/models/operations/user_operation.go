package operations

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/schemas"
	"notezy-backend/app/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserById(db *gorm.DB, id uuid.UUID) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	user := schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetUserByName(db *gorm.DB, name string) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	user := schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Where("name = ?", name).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	user := schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Where("email = ?", email).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetAllUsers(db *gorm.DB) (*[]schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	users := []schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Find(&users)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(result.Error)
	}
	return &users, nil
}

func CreateUser(db *gorm.DB, input inputs.CreateUserInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	// note that the create operation in gorm will NOT return anything
	// but the default value we set in gorm field in the above struct will be returned if we specified it in the "returning"
	var newUser schemas.User
	util.CopyNonNilFields(&newUser, input)
	result := db.Table(schemas.User{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"}, // for the following procedure such as create user info, create user account, generate refresh token etc..
		}}).
		Create(&newUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}
	return &newUser.Id, nil
}

func UpdateUserById(db *gorm.DB, id uuid.UUID, input inputs.UpdateUserInput) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	var updatedUser schemas.User
	updatedUser.UpdatedAt = time.Now()
	util.CopyNonNilFields(&updatedUser, input)
	result := db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(&updatedUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	return &updatedUser, nil
}

func DeleteUserById(db *gorm.DB, id uuid.UUID) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	tx := db.Begin()

	deletedUser := schemas.User{}
	result := tx.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		First(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithError(err)
	}

	result = tx.Table(schemas.User{}.TableName()).
		Delete(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.User.FailedToDelete().WithError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.User.FailedToDelete().WithError(err)
	}

	return &deletedUser, nil

}
