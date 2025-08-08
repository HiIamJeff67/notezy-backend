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
	GetOneByUserId(userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInfoInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserInfoInput) (*schemas.UserInfo, *exceptions.Exception)
}

type UserInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) UserInfoRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &UserInfoRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *UserInfoRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception) {
	userInfo := schemas.UserInfo{}
	result := r.db.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func (r *UserInfoRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInfoInput) (*uuid.UUID, *exceptions.Exception) {
	var newUserInfo schemas.UserInfo
	newUserInfo.UserId = userId
	if err := copier.Copy(&newUserInfo, &input); err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.UserInfo{}.TableName()).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return &newUserInfo.Id, nil
}

func (r *UserInfoRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserInfoInput) (*schemas.UserInfo, *exceptions.Exception) {
	existingUserInfo, exception := r.GetOneByUserId(userId)
	if exception != nil || existingUserInfo == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserInfo)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserInfo)
	}

	result := r.db.Table(schemas.UserInfo{}.TableName()).
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
