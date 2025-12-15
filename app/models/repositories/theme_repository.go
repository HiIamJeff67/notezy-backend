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

type ThemeRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads []schemas.ThemeRelation, opts ...options.RepositoryOptions) (*schemas.Theme, *exceptions.Exception)
	GetAll(opts ...options.RepositoryOptions) ([]schemas.Theme, *exceptions.Exception)
	CreateOneByAuthorId(authorId uuid.UUID, input inputs.CreateThemeInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, authorId uuid.UUID, input inputs.PartialUpdateThemeInput, opts ...options.RepositoryOptions) (*schemas.Theme, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, authorId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type ThemeRepository struct{}

func NewThemeRepository() ThemeRepositoryInterface {
	return &ThemeRepository{}
}

/* ============================== Implementations ============================== */

func (r *ThemeRepository) GetOneById(
	id uuid.UUID,
	preloads []schemas.ThemeRelation,
	opts ...options.RepositoryOptions,
) (*schemas.Theme, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var theme schemas.Theme

	query := parsedOptions.DB.Table(schemas.Theme{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Where("id = ?", id).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &theme, nil
}

func (r *ThemeRepository) GetAll(
	opts ...options.RepositoryOptions,
) ([]schemas.Theme, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var themes []schemas.Theme
	result := parsedOptions.DB.Table(schemas.Theme{}.TableName()).
		Find(&themes)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return themes, nil
}

func (r *ThemeRepository) CreateOneByAuthorId(
	authorId uuid.UUID,
	input inputs.CreateThemeInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newTheme schemas.Theme
	newTheme.AuthorId = authorId

	if err := copier.Copy(&newTheme, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newTheme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	return &newTheme.Id, nil
}

func (r *ThemeRepository) UpdateOneById(
	id uuid.UUID,
	authorId uuid.UUID,
	input inputs.PartialUpdateThemeInput,
	opts ...options.RepositoryOptions,
) (*schemas.Theme, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingTheme, exception := r.GetOneById(
		id,
		nil,
		opts...,
	)
	if exception != nil || existingTheme == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingTheme)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingTheme)
	}

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Where("id = ? AND author_id = ?", id, authorId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Theme.NoChanges()
	}

	return &updates, nil
}

func (r *ThemeRepository) DeleteOneById(
	id uuid.UUID,
	authorId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	// * If you need to use the funcionality of RETURNING from PostgreSQL
	// var deletedTheme schemas.Theme

	// result := r.db.Table(schemas.Theme{}.TableName()).
	// 	Where("id = ? AND author_id = ?", id, authorId).
	// 	Clauses(clause.Returning{}).
	// 	Delete(&deletedTheme)

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Where("id = ? AND author_id = ?", id, authorId).
		Delete(&schemas.Theme{})
	if err := result.Error; err != nil {
		return exceptions.Theme.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Theme.NotFound()
	}

	return nil
}
