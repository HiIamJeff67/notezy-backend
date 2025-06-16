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

type UserInfoRepository interface {
	GetOneByUserId(userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInput) *exceptions.Exception
	UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserInput) (*schemas.UserInfo, *exceptions.Exception)
}

type userInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) *userInfoRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &userInfoRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *userInfoRepository) GetOneByUserId(userId uuid.UUID) (*schemas.UserInfo, *exceptions.Exception) {
	userInfo := schemas.UserInfo{}
	result := r.db.Table(schemas.UserInfo{}.TableName()).Where("user_id = ?", userId).First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func (r *userInfoRepository) CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserInfoInput) *exceptions.Exception {
	if err := models.Validator.Struct(input); err != nil {
		return exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var newUserInfo schemas.UserInfo
	newUserInfo.UserId = userId
	util.CopyNonNilFields(&newUserInfo, input)

	result := r.db.Table(schemas.UserInfo{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return nil
}

// TODO: Partial Update
func (r *userInfoRepository) UpdateOneByUserId(userId uuid.UUID, input inputs.UpdateUserInfoInput) (*schemas.UserInfo, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var updatedUserInfo schemas.UserInfo
	util.CopyNonNilFields(&updatedUserInfo, input)
	result := r.db.Table(schemas.UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Updates(&updatedUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}

	return &updatedUserInfo, nil
}
