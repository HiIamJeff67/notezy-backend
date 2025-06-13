package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	"notezy-backend/app/util"
)

type UserRepository interface {
	GetOneById(id uuid.UUID) (*schemas.User, *exceptions.Exception)
	GetOneByName(name string) (*schemas.User, *exceptions.Exception)
	GetOneByEmail(email string) (*schemas.User, *exceptions.Exception)
	GetAll() (*[]schemas.User, *exceptions.Exception)
	CreateOne(input inputs.CreateUserInput) *exceptions.Exception
	UpdateOneById(id uuid.UUID, input inputs.UpdateUserInput) (*schemas.User, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, input inputs.DeleteUserInput) *exceptions.Exception
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &userRepository{db: db}
}

func (r *userRepository) GetOneById(id uuid.UUID) (*schemas.User, *exceptions.Exception) {
	user := schemas.User{}
	result := r.db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *userRepository) GetOneByName(name string) (*schemas.User, *exceptions.Exception) {
	user := schemas.User{}
	result := r.db.Table(schemas.User{}.TableName()).
		Where("name = ?", name).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *userRepository) GetOneByEmail(email string) (*schemas.User, *exceptions.Exception) {
	user := schemas.User{}
	result := r.db.Table(schemas.User{}.TableName()).
		Where("email = ?", email).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *userRepository) GetAll() (*[]schemas.User, *exceptions.Exception) {
	users := []schemas.User{}

	result := r.db.Preload("UserInfo").
		Preload("UserAccount").
		Preload("UserSetting").
		Preload("Badges").
		Preload("Themes").
		Find(&users)

	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(result.Error)
	}
	return &users, nil
}

func (r *userRepository) CreateOne(input inputs.CreateUserInput) (*uuid.UUID, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	// note that the create operation in gorm will NOT return anything
	// but the default value we set in gorm field in the above struct will be returned if we specified it in the "returning"
	var newUser schemas.User
	if err := copier.Copy(&newUser, &input); err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.User{}.TableName()).
		Create(&newUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}
	return &newUser.Id, nil
}

func (r *userRepository) UpdateOneById(id uuid.UUID, input inputs.PartialUpdateUserInput) (*schemas.User, *exceptions.Exception) {
	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err).Log()
	}

	values := input.Values
	setNull := input.SetNull
	existingUser, exception := r.GetOneById(id)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(values, *setNull, *existingUser)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(values, *setNull, *existingUser)
	}

	result := r.db.Table(schemas.User{}.TableName()).
		Where("id = ?", id).
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NotFound()
	}
	return &updates, nil
}

func (r *userRepository) DeleteOneById(id uuid.UUID, input inputs.DeleteUserInput) *exceptions.Exception {
	deletedUser := schemas.User{}
	result := r.db.Table(schemas.User{}.TableName()).
		Where("id = ? AND name = ? AND password", id, input.Name, input.Password).
		Clauses(clause.Returning{}).
		Delete(&deletedUser)
	if err := result.Error; err != nil {
		return exceptions.User.NotFound().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.User.FailedToDelete()
	}
	return nil

}
