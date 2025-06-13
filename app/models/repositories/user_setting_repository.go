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

func GetUserSettingByUserId(db *gorm.DB, userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	userSetting := schemas.UserSetting{}
	result := db.Table(schemas.UserSetting{}.TableName()).Where("user_id = ?", userId).First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func CreateUserSettingByUserId(db *gorm.DB, userId uuid.UUID, input inputs.CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var newUserSetting schemas.UserSetting
	newUserSetting.UserId = userId
	util.CopyNonNilFields(&newUserSetting, input)
	result := db.Table(schemas.UserSetting{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return &newUserSetting.Id, nil
}

func UpdateUserSettingByUserId(db *gorm.DB, userId uuid.UUID, input inputs.UpdateUserSettingInput) (*schemas.UserSetting, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var updatedUserSetting schemas.UserSetting
	updatedUserSetting.UserId = userId
	util.CopyNonNilFields(&updatedUserSetting, input)
	result := db.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Create(&updatedUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToUpdate().WithError(err)
	}

	return &updatedUserSetting, nil
}
