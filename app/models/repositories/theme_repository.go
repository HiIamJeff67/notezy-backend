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

type ThemeRepository interface {
	GetOneById(id uuid.UUID) (*schemas.Theme, *exceptions.Exception)
	GetAll() (*[]schemas.Theme, *exceptions.Exception)
	CreateOneByAuthorId(authorId uuid.UUID, input inputs.CreateThemeInput) *exceptions.Exception
	UpdateOneById(id uuid.UUID, authorId uuid.UUID, input inputs.UpdateThemeInput) *exceptions.Exception
	DeleteOneById(id uuid.UUID, authorId uuid.UUID) *exceptions.Exception
}

type themeRepository struct {
	db *gorm.DB
}

func NewThemeRepository(db *gorm.DB) *themeRepository {
	if db == nil {
		db = models.NotezyDB
	}
	return &themeRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *themeRepository) GetOneById(id uuid.UUID) (*schemas.Theme, *exceptions.Exception) {
	theme := schemas.Theme{}
	result := r.db.Table(schemas.Theme{}.TableName()).
		Where("id = ?", id).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &theme, nil
}

func (r *themeRepository) GetAll() (*[]schemas.Theme, *exceptions.Exception) {
	themes := []schemas.Theme{}
	result := r.db.Table(schemas.Theme{}.TableName()).
		Find(&themes)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &themes, nil
}

func (r *themeRepository) CreateOneByAuthorId(authorId uuid.UUID, input inputs.CreateThemeInput) *exceptions.Exception {
	if err := models.Validator.Struct(input); err != nil {
		return exceptions.Theme.FailedToCreate().WithError(err)
	}

	var newTheme schemas.Theme
	newTheme.AuthorId = authorId
	util.CopyNonNilFields(&newTheme, input)
	result := r.db.Table(schemas.Theme{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newTheme)
	if err := result.Error; err != nil {
		return exceptions.Theme.FailedToCreate().WithError(err)
	}
	return nil
}

func (r *themeRepository) UpdateOneById(id uuid.UUID, authorId uuid.UUID, input inputs.UpdateThemeInput) *exceptions.Exception {
	if err := models.Validator.Struct(input); err != nil {
		return exceptions.Theme.FailedToCreate().WithError(err)
	}

	var updatedTheme schemas.Theme
	util.CopyNonNilFields(&updatedTheme, input)
	result := r.db.Table(schemas.Theme{}.TableName()).
		Where("id = ? AND author_id = ?", id, authorId).
		Clauses(clause.Returning{}).
		Updates(&updatedTheme)
	if err := result.Error; err != nil {
		return exceptions.Theme.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return exceptions.Theme.FailedToUpdate()
	}
	return nil
}

func (r *themeRepository) DeleteOneById(id uuid.UUID, authorId uuid.UUID) *exceptions.Exception {
	var deletedTheme schemas.Theme
	result := r.db.Table(schemas.Theme{}.TableName()).
		Where("id = ? AND author_id = ?", id, authorId).
		Clauses(clause.Returning{}).
		Delete(&deletedTheme)
	if err := result.Error; err != nil {
		return exceptions.Theme.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Theme.FailedToDelete()
	}
	return nil
}
