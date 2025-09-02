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

type UserSettingRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput) (*schemas.UserSetting, *exceptions.Exception)
}

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) UserSettingRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserSettingRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *UserSettingRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserSetting, *exceptions.Exception) {
	userSetting := schemas.UserSetting{}
	result := r.db.Table(schemas.UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func (r *UserSettingRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception) {
	var newUserSetting schemas.UserSetting
	newUserSetting.UserId = userId
	if err := copier.Copy(&newUserSetting, &input); err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	result := r.db.Model(&schemas.UserSetting{}).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return &newUserSetting.Id, nil
}

func (r *UserSettingRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserSettingInput) (*schemas.UserSetting, *exceptions.Exception) {
	existingUserSetting, exception := r.GetOneByUserId(userId)
	if exception != nil || existingUserSetting == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserSetting)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserSetting)
	}

	result := r.db.Model(&schemas.UserSetting{}).
		Where("user_id = ?").
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserSetting.NoChanges()
	}

	return &updates, nil
}
