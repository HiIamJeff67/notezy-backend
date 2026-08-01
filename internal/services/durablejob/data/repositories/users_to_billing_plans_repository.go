package repositories

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/scopes"
	durablejobexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/exceptions"
	partialupdate "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/partialupdate"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type UsersToBillingPlansRepositoryInterface interface {
	GetOnyById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UsersToBillingPlans, *exceptions.Exception)
	GetAllByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.UsersToBillingPlans, *exceptions.Exception)
	CreateOne(userId uuid.UUID, input inputs.CreateUsersToBillingPlansInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateUsersToBillingPlansInput, opts ...options.RepositoryOptions) (*schemas.UsersToBillingPlans, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	DeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UsersToBillingPlansRepository struct{}

func NewUsersToBillingPlansRepository() UsersToBillingPlansRepositoryInterface {
	return &UsersToBillingPlansRepository{}
}

func (r *UsersToBillingPlansRepository) GetOnyById(
	id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions,
) (*schemas.UsersToBillingPlans, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var usersToBillingPlans schemas.UsersToBillingPlans
	result := parsedOptions.DB.Table(schemas.UsersToBillingPlans{}.TableName()).
		Where("id = ? and user_id = ?", id, userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&usersToBillingPlans)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.NotFound().WithOrigin(result.Error)},
		{First: usersToBillingPlans.Id == uuid.Nil, Second: durablejobexceptions.UsersToBillingPlans.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &usersToBillingPlans, nil
}

func (r *UsersToBillingPlansRepository) GetAllByUserId(
	userId uuid.UUID, opts ...options.RepositoryOptions,
) ([]schemas.UsersToBillingPlans, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var usersToBillingPlans []schemas.UsersToBillingPlans
	result := parsedOptions.DB.Table(schemas.UsersToBillingPlans{}.TableName()).
		Where("user_id = ?", userId).
		Find(&usersToBillingPlans)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.NotFound().WithOrigin(result.Error)},
		{First: len(usersToBillingPlans) == 0, Second: durablejobexceptions.UsersToBillingPlans.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return usersToBillingPlans, nil
}

func (r *UsersToBillingPlansRepository) CreateOne(
	userId uuid.UUID,
	input inputs.CreateUsersToBillingPlansInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUsersToBillingPlans schemas.UsersToBillingPlans
	newUsersToBillingPlans.UserId = userId

	if err := copier.Copy(&newUsersToBillingPlans, &input); err != nil {
		return nil, durablejobexceptions.UsersToBillingPlans.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUsersToBillingPlans)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.FailedToCreate().WithOrigin(result.Error)},
		{First: newUsersToBillingPlans.Id == uuid.Nil, Second: durablejobexceptions.UsersToBillingPlans.FailedToCreate()},
	}); exception != nil {
		return nil, exception
	}

	return &newUsersToBillingPlans.Id, nil
}

func (r *UsersToBillingPlansRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateUsersToBillingPlansInput,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToBillingPlans, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUsersToBillingPlans, exception := r.GetOnyById(
		id,
		userId,
		opts...,
	)
	if exception := exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: existingUsersToBillingPlans == nil, Second: durablejobexceptions.UsersToBillingPlans.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	updates, err := partialupdate.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUsersToBillingPlans)
	if err != nil {
		return nil, exceptions.New("FailedToPreprocessPartialUpdate", "Repository", "Update", "Failed to preprocess partial update", http.StatusInternalServerError, true)
	}

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Where("id = ? and user_id = ?", id, userId).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: durablejobexceptions.UsersToBillingPlans.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}

func (r *UsersToBillingPlansRepository) DeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Where("id = ? and user_id = ?", id, userId).
		Delete(&schemas.UsersToBillingPlans{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: durablejobexceptions.UsersToBillingPlans.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *UsersToBillingPlansRepository) DeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return durablejobexceptions.UsersToBillingPlans.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Where("ids IN ? and user_id = ?", ids, userId).
		Delete(&schemas.UsersToBillingPlans{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: durablejobexceptions.UsersToBillingPlans.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: durablejobexceptions.UsersToBillingPlans.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
