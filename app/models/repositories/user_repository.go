package repositories

import (
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type UserRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetOneByName(name string, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetOneByEmail(email string, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetAll(opts ...options.RepositoryOptions) ([]schemas.User, *exceptions.Exception)
	CreateOne(input inputs.CreateUserInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, input inputs.PartialUpdateUserInput, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, input inputs.DeleteUserInput, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UserRepository struct{}

func NewUserRepository() UserRepositoryInterface {
	return &UserRepository{}
}

/* ============================== Implementations ============================== */

func (r *UserRepository) GetOneById(
	id uuid.UUID,
	preloads []schemas.UserRelation,
	opts ...options.RepositoryOptions,
) (*schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	user := schemas.User{}

	db := parsedOptions.DB.Table(schemas.User{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("id = ?", id).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *UserRepository) GetOneByName(
	name string,
	preloads []schemas.UserRelation,
	opts ...options.RepositoryOptions,
) (*schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	user := schemas.User{}

	db := parsedOptions.DB.Table(schemas.User{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("name = ?", name).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *UserRepository) GetOneByEmail(
	email string,
	preloads []schemas.UserRelation,
	opts ...options.RepositoryOptions,
) (*schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	user := schemas.User{}

	query := parsedOptions.DB.Table(schemas.User{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Where("email = ?", email).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func (r *UserRepository) GetAll(
	opts ...options.RepositoryOptions,
) ([]schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	users := []schemas.User{}

	result := parsedOptions.DB.Preload("UserInfo").
		Preload("UserAccount").
		Preload("UserSetting").
		Preload("Badges").
		Preload("Themes").
		Find(&users)

	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(result.Error)
	}

	return users, nil
}

func (r *UserRepository) CreateOne(
	input inputs.CreateUserInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	// note that the create operation in gorm will NOT return anything
	// but the default value we set in gorm field in the above struct will be returned if we specified it in the "returning"
	var newUser schemas.User
	if err := copier.Copy(&newUser, &input); err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.User{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUser)
	if err := result.Error; err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_name\" (SQLSTATE 23505)":
			return nil, exceptions.User.DuplicateName(input.Name)
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_email\" (SQLSTATE 23505)":
			return nil, exceptions.User.DuplicateEmail(input.Email)
		default:
			return nil, exceptions.User.FailedToCreate() // .WithError(err) <- don't show the database error to outside
		}
	}

	return &newUser.Id, nil
}

func (r *UserRepository) UpdateOneById(
	id uuid.UUID,
	input inputs.PartialUpdateUserInput,
	opts ...options.RepositoryOptions,
) (*schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUser, exception := r.GetOneById(
		id,
		nil,
		opts...,
	)
	if exception != nil || existingUser == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUser)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUser)
	}

	result := parsedOptions.DB.Model(&schemas.User{}).
		Where("id = ?", id).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.User.NoChanges()
	}

	return &updates, nil
}

func (r *UserRepository) DeleteOneById(
	id uuid.UUID,
	input inputs.DeleteUserInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.User{}).
		Where("id = ? AND name = ? AND password", id, input.Name, input.Password).
		Delete(&schemas.User{})
	if err := result.Error; err != nil {
		return exceptions.User.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.User.NoChanges()
	}

	return nil
}
