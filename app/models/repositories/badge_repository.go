package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	models "notezy-backend/app/models"
	schemas "notezy-backend/app/models/schemas"
)

/* ============================== Definitions ============================== */

type BadgeRepositoryInterface interface {
	GetOneById(id uuid.UUID) (*schemas.Badge, *exceptions.Exception)

	// repository for public badges
	GetPublicOneByEncodedSearchCursor(encodedSearchCursor string) (*gqlmodels.PublicBadge, *exceptions.Exception)
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

func (r *BadgeRepository) GetOneById(id uuid.UUID) (*schemas.Badge, *exceptions.Exception) {
	badge := schemas.Badge{}
	result := r.db.Table(schemas.Badge{}.TableName()).
		Where("id = ?", id).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return &badge, nil
}

func (r *BadgeRepository) GetPublicOneByEncodedSearchCursor(encodedSearchCursor string) (*gqlmodels.PublicBadge, *exceptions.Exception) {
	badge := schemas.Badge{}
	result := r.db.Table(schemas.Badge{}.TableName()).
		Where("encoded_search_cursor = ?", encodedSearchCursor).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return badge.ToPublicBadge(), nil
}
