package repositories

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

type UsersToBillingPlansRepositoryInterface interface {
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
		First(&usersToBillingPlans)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.NotFound().WithOrigin(result.Error)},
		{First: usersToBillingPlans.Id == uuid.Nil, Second: exceptions.UsersToBillingPlans.NotFound()},
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
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.NotFound().WithOrigin(result.Error)},
		{First: len(usersToBillingPlans) == 0, Second: exceptions.UsersToBillingPlans.NotFound()},
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
		return nil, exceptions.UsersToBillingPlans.FailedToCreate().WithOrigin(err)
	}

	var createdUsersToBillingPlans schemas.UsersToBillingPlans
	result := parsedOptions.DB.Model(&createdUsersToBillingPlans).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUsersToBillingPlans)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.FailedToCreate().WithOrigin(result.Error)},
		{First: createdUsersToBillingPlans.Id == uuid.Nil, Second: exceptions.UsersToBillingPlans.FailedToCreate()},
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
		{First: existingUsersToBillingPlans == nil, Second: exceptions.UsersToBillingPlans.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUsersToBillingPlans)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUsersToBillingPlans)
	}

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Where("id = ? and user_id = ?", id, userId).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UsersToBillingPlans.NoChanges()},
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
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UsersToBillingPlans.NoChanges()},
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
		return exceptions.UsersToBillingPlans.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.UsersToBillingPlans{}).
		Where("ids IN ? and user_id = ?", ids, userId).
		Delete(&schemas.UsersToBillingPlans{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.UsersToBillingPlans.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.UsersToBillingPlans.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
