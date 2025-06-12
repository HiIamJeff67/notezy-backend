package operations

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
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

	result := db.Preload("UserInfo").
		Preload("UserAccount").
		Preload("UserSetting").
		// Preload("Badges").
		Preload("Themes").
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
	if err := copier.Copy(&newUser, &input); err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}

	result := db.Table(schemas.User{}.TableName()).
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
		return nil, exceptions.User.InvalidInput().WithError(err).Log()
	}

	var updatedUser schemas.User
	if err := copier.Copy(&updatedUser, &input); err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}

	result := db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(&updatedUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}
	return &updatedUser, nil
}

func DeleteUserById(db *gorm.DB, id uuid.UUID) (*schemas.User, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	deletedUser := schemas.User{}
	result := db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Delete(&deletedUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.FailedToDelete()
	}
	return &deletedUser, nil

}
