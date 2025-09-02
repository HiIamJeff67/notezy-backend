package repositories

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type ThemeRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads *[]schemas.ThemeRelation) (*schemas.Theme, *exceptions.Exception)
	GetAll() (*[]schemas.Theme, *exceptions.Exception)
	CreateOneByAuthorId(authorId uuid.UUID, input inputs.CreateThemeInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, authorId uuid.UUID, input inputs.PartialUpdateThemeInput) (*schemas.Theme, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, authorId uuid.UUID) *exceptions.Exception
}

type ThemeRepository struct {
	db *gorm.DB
}

func NewThemeRepository(db *gorm.DB) ThemeRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &ThemeRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *ThemeRepository) GetOneById(id uuid.UUID, preloads *[]schemas.ThemeRelation) (*schemas.Theme, *exceptions.Exception) {
	theme := schemas.Theme{}
	db := r.db.Table(schemas.Theme{}.TableName())
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("id = ?", id).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &theme, nil
}

func (r *ThemeRepository) GetAll() (*[]schemas.Theme, *exceptions.Exception) {
	themes := []schemas.Theme{}

	result := r.db.Table(schemas.Theme{}.TableName()).
		Find(&themes)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &themes, nil
}

func (r *ThemeRepository) CreateOneByAuthorId(authorId uuid.UUID, input inputs.CreateThemeInput) (*uuid.UUID, *exceptions.Exception) {
	var newTheme schemas.Theme
	newTheme.AuthorId = authorId
	if err := copier.Copy(&newTheme, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := r.db.Model(&schemas.Theme{}).
		Create(&newTheme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	return &newTheme.Id, nil
}

func (r *ThemeRepository) UpdateOneById(id uuid.UUID, authorId uuid.UUID, input inputs.PartialUpdateThemeInput) (*schemas.Theme, *exceptions.Exception) {
	existingTheme, exception := r.GetOneById(id, nil)
	if exception != nil || existingTheme == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingTheme)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingTheme)
	}

	result := r.db.Model(&schemas.Theme{}).
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

func (r *ThemeRepository) DeleteOneById(id uuid.UUID, authorId uuid.UUID) *exceptions.Exception {
	// * If you need to use the funcionality of RETURNING from PostgreSQL
	// var deletedTheme schemas.Theme

	// result := r.db.Table(schemas.Theme{}.TableName()).
	// 	Where("id = ? AND author_id = ?", id, authorId).
	// 	Clauses(clause.Returning{}).
	// 	Delete(&deletedTheme)

	result := r.db.Model(&schemas.Theme{}).
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
