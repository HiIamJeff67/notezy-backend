package services

import (
	"context"

	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	models "notezy-backend/app/models"
	schemas "notezy-backend/app/models/schemas"
)

/* ============================== Interface & Instance ============================== */

type BadgeServiceInterface interface {
	GetPublicBadgeByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception)
	GetPublicBadgeByPublicUserId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception)
	GetPublicBadgesByPublicUserIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicBadge, *exceptions.Exception)
}

type BadgeService struct {
	db *gorm.DB
}

func NewBadgeService(db *gorm.DB) BadgeServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &BadgeService{db: db}
}

/* ============================== Service Methods for Badge ============================== */

/* ============================== SErvices for Public Badge ============================== */

func (s *BadgeService) GetPublicBadgeByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception) {
	badge := schemas.Badge{}
	result := s.db.Table(schemas.Badge{}.TableName()).
		Where("public_id = ?", publicId).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return badge.ToPublicBadge(), nil
}

func (s *BadgeService) GetPublicBadgeByPublicUserId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception) {
	badge := schemas.Badge{}
	result := s.db.Table(schemas.Badge{}.TableName()+" b").
		Select("b.*, utb.user_id").
		Joins("LEFT JOIN \"UsersToBadgesTable\" ON utb.badge_id = b.id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = utb.user_id").
		Where("u.public_id = ?", publicId).
		First(&badge)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	return badge.ToPublicBadge(), nil
}

func (s *BadgeService) GetPublicBadgesByPublicUserIds(ctx context.Context, publicIds []string) ([]*gqlmodels.PublicBadge, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicBadge{}, nil
	}

	var badges []*schemas.Badge
	result := s.db.Table(schemas.Badge{}.TableName()+" b").
		Select("b.*, utb.user_id").
		Joins("LEFT JOIN \"UsersToBadgesTable\" ON utb.badge_id = b.id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = utb.user_id").
		Where("u.public_id IN ?", publicIds).
		Find(&badges)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	publicBadges := make([]*gqlmodels.PublicBadge, len(badges))
	for index, badge := range badges {
		publicBadges[index] = badge.ToPublicBadge()
	}

	return publicBadges, nil
}
