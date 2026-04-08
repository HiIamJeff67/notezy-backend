package repositories

import (
	"github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	types "notezy-backend/shared/types"
)

type BadgeRepositoryInterface interface {
	GetOneById(id uuid.UUID, preloads []schemas.BadgeRelation, opts ...options.RepositoryOptions) (*schemas.Badge, *exceptions.Exception)
}

type BadgeRepository struct{}

func NewBadgeRepository() BadgeRepositoryInterface {
	return &BadgeRepository{}
}

func (r *BadgeRepository) GetOneById(
	id uuid.UUID,
	preloads []schemas.BadgeRelation,
	opts ...options.RepositoryOptions,
) (*schemas.Badge, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	badge := schemas.Badge{}

	query := parsedOptions.DB.Table(schemas.Badge{}.TableName())
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Where("id = ?", id).
		First(&badge)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Badge.NotFound().WithOrigin(result.Error)},
		{First: badge.Id == uuid.Nil, Second: exceptions.Badge.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &badge, nil
}
