package repositories

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	"notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type UserSettingRepository interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput) (*schemas.UserSetting, *exceptions.Exception)
}

type userSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) UserSettingRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &userSettingRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *userSettingRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception) {
	userSetting := schemas.UserSetting{}
	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func (r *userSettingRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserSetting.InvalidInput().WithError(err)
	}

	var newUserSetting schemas.UserSetting
	newUserSetting.UserId = userId
	if err := copier.Copy(&newUserSetting, &input); err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return &newUserSetting.Id, nil
}

func (r *userSettingRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput) (*schemas.UserSetting, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserSetting.InvalidInput().WithError(err)
	}

	existingUserSetting, exception := r.GetOneByUserId(userId)
	if exception != nil || existingUserSetting == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserSetting)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserSetting)
	}

	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserSetting.NotFound()
	}

	return &updates, nil
}
