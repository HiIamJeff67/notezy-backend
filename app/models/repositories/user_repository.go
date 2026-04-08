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
	types "notezy-backend/shared/types"
)

type UserRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetOneByName(name string, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetOneByEmail(email string, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetAll(opts ...options.RepositoryOptions) ([]schemas.User, *exceptions.Exception)
	CreateOne(input inputs.CreateUserInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, input inputs.PartialUpdateUserInput, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
}

type UserRepository struct{}

func NewUserRepository() UserRepositoryInterface {
	return &UserRepository{}
}

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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: exceptions.User.NotFound()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: exceptions.User.NotFound()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: exceptions.User.NotFound()},
	}); exception != nil {
		return nil, exception
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.User.NotFound().WithOrigin(result.Error)},
		{First: len(users) == 0, Second: exceptions.User.NotFound()},
	}); exception != nil {
		return nil, exception
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
		return nil, exceptions.User.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.User{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUser)
	if err := result.Error; err != nil {
		// instead of using exceptions.Cover(), we can just get the error string and switch on it to return the corresponded exceptions
		// this approach is faster and more straight forward
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_name\" (SQLSTATE 23505)":
			return nil, exceptions.User.DuplicateName(input.Name)
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_email\" (SQLSTATE 23505)":
			return nil, exceptions.User.DuplicateEmail(input.Email)
		default:
			return nil, exceptions.User.FailedToCreate() // .WithOrigin(err) <- don't show the database error to outside
		}
	}
	if result.RowsAffected == 0 {
		// check the remaining condition here,
		// since there's only 1 more condition to check,
		// there's no need to use exceptions.Cover() to map all the it
		return nil, exceptions.User.NoChanges()
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
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.User.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.User.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}
