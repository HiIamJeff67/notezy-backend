package repositories

import (
	"gorm.io/gorm/clause"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/scopes"
	durablejobexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/exceptions"
	partialupdate "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/partialupdate"
)

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
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.Theme.NotFound().WithOrigin(err)
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
		return nil, durablejobexceptions.Theme.NotFound().WithOrigin(err)
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
		return nil, durablejobexceptions.Theme.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newTheme)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.Theme.FailedToCreate().WithOrigin(err)
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

	updates, err := partialupdate.PartialUpdatePreprocess(input.Values, input.SetNull, *existingTheme)
	if err != nil {
		return nil, exceptions.New("FailedToPreprocessPartialUpdate", "Repository", "Update", "Failed to preprocess partial update", http.StatusInternalServerError, true)
	}

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Where("id = ? AND author_id = ?", id, authorId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.Theme.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, durablejobexceptions.Theme.NoChanges()
	}

	return &updates, nil
}

func (r *ThemeRepository) DeleteOneById(
	id uuid.UUID,
	authorId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	// * If you need to use the functionality of RETURNING from PostgreSQL
	// var deletedTheme schemas.Theme

	// result := r.db.Table(schemas.Theme{}.TableName()).
	// 	Where("id = ? AND author_id = ?", id, authorId).
	// 	Clauses(clause.Returning{}).
	// 	Delete(&deletedTheme)

	result := parsedOptions.DB.Model(&schemas.Theme{}).
		Where("id = ? AND author_id = ?", id, authorId).
		Delete(&schemas.Theme{})
	if err := result.Error; err != nil {
		return durablejobexceptions.Theme.FailedToDelete().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return durablejobexceptions.Theme.NotFound()
	}

	return nil
}
