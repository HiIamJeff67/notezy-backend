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

type UserInfoRepositoryInterface interface {
	GetOneByUserId(db *gorm.DB, userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception)
	CreateOneByUserId(db *gorm.DB, userId uuid.UUID, input inputs.CreateUserInfoInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(db *gorm.DB, userId uuid.UUID, input inputs.PartialUpdateUserInfoInput) (*schemas.UserInfo, *exceptions.Exception)
}

type UserInfoRepository struct{}

func NewUserInfoRepository() UserInfoRepositoryInterface {
	return &UserInfoRepository{}
}

/* ============================== CRUD operations ============================== */

func (r *UserInfoRepository) GetOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
) (*schemas.UserInfo, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	userInfo := schemas.UserInfo{}
	result := db.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func (r *UserInfoRepository) CreateOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
	input inputs.CreateUserInfoInput,
) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var newUserInfo schemas.UserInfo
	newUserInfo.UserId = userId
	if err := copier.Copy(&newUserInfo, &input); err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	result := db.Model(&schemas.UserInfo{}).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return &newUserInfo.Id, nil
}

func (r *UserInfoRepository) UpdateOneByUserId(
	db *gorm.DB,
	userId uuid.UUID,
	input inputs.PartialUpdateUserInfoInput,
) (*schemas.UserInfo, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	existingUserInfo, exception := r.GetOneByUserId(db, userId)
	if exception != nil || existingUserInfo == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserInfo)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserInfo)
	}

	result := db.Model(&schemas.UserInfo{}).
		Where("user_id = ?", userId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserInfo.NoChanges()
	}

	return &updates, nil
}
