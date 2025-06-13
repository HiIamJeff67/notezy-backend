package repositories

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

func GetUserInfoByUserId(db *gorm.DB, userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	userInfo := schemas.UserInfo{}
	result := db.Table(schemas.UserInfo{}.TableName()).Where("user_id = ?", userId).First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func CreateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input inputs.CreateUserInfoInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var newUserInfo schemas.UserInfo
	newUserInfo.UserId = userId
	util.CopyNonNilFields(&newUserInfo, input)

	result := db.Table(schemas.UserInfo{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return &newUserInfo.Id, nil
}

func UpdateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input inputs.UpdateUserInfoInput) (*schemas.UserInfo, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var updatedUserInfo schemas.UserInfo
	util.CopyNonNilFields(&updatedUserInfo, input)
	result := db.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Updates(&updatedUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}

	return &updatedUserInfo, nil
}
