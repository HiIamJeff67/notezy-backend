package operations

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/schemas"
	"notezy-backend/app/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetThemeById(db *gorm.DB, id uuid.UUID) (*schemas.Theme, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	theme := schemas.Theme{}
	result := db.Table(schemas.Theme{}.TableName()).
		Where("id = ?", id).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &theme, nil
}

func GetAllThemes(db *gorm.DB) (*[]schemas.Theme, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	themes := []schemas.Theme{}
	result := db.Table(schemas.Theme{}.TableName()).
		Find(&themes)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return &themes, nil
}

func CreateThemeByAuthorId(db *gorm.DB, input inputs.CreateThemeInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	var newTheme schemas.Theme
	util.CopyNonNilFields(&newTheme, input)
	result := db.Table(schemas.Theme{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newTheme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}
	return &newTheme.Id, nil
}

func UpdateThemeById(db *gorm.DB, id uuid.UUID, authorId uuid.UUID, input inputs.UpdateThemeInput) (*schemas.Theme, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	if err := models.Validator.Struct(input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	var updatedTheme schemas.Theme
	util.CopyNonNilFields(&updatedTheme, input)
	result := db.Table(schemas.Theme{}.TableName()).
		Where("id = ? AND author_id = ?", id, authorId).
		Clauses(clause.Returning{}).
		Updates(&updatedTheme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Theme.FailedToUpdate()
	}
	return &updatedTheme, nil
}

func DeleteThemeById(db *gorm.DB, id uuid.UUID, authorId uuid.UUID) (*schemas.Theme, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	var deletedTheme schemas.Theme
	result := db.Table(schemas.Theme{}.TableName()).
		Where("id = ? AND author_id = ?", id, authorId).
		Clauses(clause.Returning{}).
		Delete(&deletedTheme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Theme.FailedToDelete()
	}
	return &deletedTheme, nil
}
