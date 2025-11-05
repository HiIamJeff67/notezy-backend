package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	schemas "notezy-backend/app/models/schemas"
)

/* ============================== Definitions ============================== */

type BadgeRepositoryInterface interface {
	GetOneById(db *gorm.DB, id uuid.UUID, preloads []schemas.BadgeRelation) (*schemas.Badge, *exceptions.Exception)
}

type BadgeRepository struct{}

func NewBadgeRepository() BadgeRepositoryInterface {
	return &BadgeRepository{}
}

/* ============================== CRUD operations ============================== */

func (r *BadgeRepository) GetOneById(
	db *gorm.DB,
	id uuid.UUID,
	preloads []schemas.BadgeRelation,
) (*schemas.Badge, *exceptions.Exception) {
	if db == nil {
		db = models.NotezyDB
	}

	badge := schemas.Badge{}

	query := db.Table(schemas.Badge{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
	}

	result := query.Where("id = ?", id).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return &badge, nil
}
