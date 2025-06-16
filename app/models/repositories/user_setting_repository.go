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

type UserSettingRepository interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) *exceptions.Exception
	UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserSettingInput) *exceptions.Exception
}

type userSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *userSettingRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &userSettingRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *userSettingRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception) {
	userSetting := schemas.UserSetting{}
	result := r.db.Table(schemas.UserSetting{}.TableName()).Where("user_id = ?", userId).First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func (r *userSettingRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) *exceptions.Exception {
	var newUserSetting schemas.UserSetting
	newUserSetting.UserId = userId
	util.CopyNonNilFields(&newUserSetting, input)
	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return nil
}

func (r *userSettingRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserSettingInput) *exceptions.Exception {
	var updatedUserSetting schemas.UserSetting
	updatedUserSetting.UserId = userId
	util.CopyNonNilFields(&updatedUserSetting, input)
	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Create(&updatedUserSetting)
	if err := result.Error; err != nil {
		return exceptions.UserSetting.FailedToUpdate().WithError(err)
	}

	return nil
}
