package repositories

import (
	"gorm.io/gorm/clause"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/options"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/scopes"
	apiexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/core/exceptions"
	partialupdate "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/partialupdate"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type UserRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
	GetOneByPublicId(publicId uuid.UUID, preloads []schemas.UserRelation, opts ...options.RepositoryOptions) (*schemas.User, *exceptions.Exception)
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&user)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: apiexceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: apiexceptions.User.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &user, nil
}

func (r *UserRepository) GetOneByPublicId(
	publicId uuid.UUID,
	preloads []schemas.UserRelation,
	opts ...options.RepositoryOptions,
) (*schemas.User, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	user := schemas.User{}
	query := parsedOptions.DB.Table(schemas.User{}.TableName())
	for _, preload := range preloads {
		query = query.Preload(string(preload))
	}

	result := query.
		Where("public_id = ?", publicId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&user)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: apiexceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: apiexceptions.User.NotFound()},
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&user)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: apiexceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: apiexceptions.User.NotFound()},
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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&user)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: apiexceptions.User.NotFound().WithOrigin(result.Error)},
		{First: user.Id == uuid.Nil, Second: apiexceptions.User.NotFound()},
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
		{First: result.Error != nil, Second: apiexceptions.User.NotFound().WithOrigin(result.Error)},
		{First: len(users) == 0, Second: apiexceptions.User.NotFound()},
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
		return nil, apiexceptions.User.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.User{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUser)
	if err := result.Error; err != nil {
		// instead of using exceptions.Cover(), we can just get the error string and switch on it to return the corresponded exceptions
		// this approach is faster and more straight forward
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_name\" (SQLSTATE 23505)":
			return nil, apiexceptions.User.DuplicateName(input.Name)
		case "ERROR: duplicate key value violates unique constraint \"uni_UserTable_email\" (SQLSTATE 23505)":
			return nil, apiexceptions.User.DuplicateEmail(input.Email)
		default:
			return nil, apiexceptions.User.FailedToCreate() // .WithOrigin(err) <- don't show the database error to outside
		}
	}
	if result.RowsAffected == 0 {
		// check the remaining condition here,
		// since there's only 1 more condition to check,
		// there's no need to use exceptions.Cover() to map all the it
		return nil, apiexceptions.User.NoChanges()
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

	updates, err := partialupdate.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUser)
	if err != nil {
		return nil, exceptions.New("FailedToPreprocessPartialUpdate", "Repository", "Update", "Failed to preprocess partial update", http.StatusInternalServerError, true)
	}

	result := parsedOptions.DB.Model(&schemas.User{}).
		Where("id = ?", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: apiexceptions.User.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: apiexceptions.User.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}
