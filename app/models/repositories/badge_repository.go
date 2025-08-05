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
	GetOneById(id uuid.UUID, preloads *[]schemas.BadgeRelation) (*schemas.Badge, *exceptions.Exception)
}

type BadgeRepository struct {
	db *gorm.DB
}

func NewBadgeRepository(db *gorm.DB) BadgeRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &BadgeRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *BadgeRepository) GetOneById(id uuid.UUID, preloads *[]schemas.BadgeRelation) (*schemas.Badge, *exceptions.Exception) {
	badge := schemas.Badge{}
	db := r.db.Table(schemas.Badge{}.TableName())
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("id = ?", id).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return &badge, nil
}
