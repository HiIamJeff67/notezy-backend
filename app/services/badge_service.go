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
	// services for public badges
	GetPublicBadgeByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception)
	GetPublicBadgeByUserPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception)
	GetPublicBadgesByUserPublicIds(ctx context.Context, publicIds []string, requiredStatic bool) ([]*gqlmodels.PublicBadge, *exceptions.Exception)
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

func (s *BadgeService) GetPublicBadgeByUserPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicBadge, *exceptions.Exception) {
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

func (s *BadgeService) GetPublicBadgesByUserPublicIds(ctx context.Context, publicIds []string, requiredStatic bool) ([]*gqlmodels.PublicBadge, *exceptions.Exception) {
	if len(publicIds) == 0 {
		return []*gqlmodels.PublicBadge{}, nil
	}

	var badgesWithPublicUserIds []*struct {
		schemas.Badge
		UserPublicId string `gorm:"column:public_id"`
	}
	result := s.db.Table(schemas.Badge{}.TableName()+" b").
		Select("b.*, u.public_id").
		Joins("LEFT JOIN \"UsersToBadgesTable\" utb ON utb.badge_id = b.id").
		Joins("LEFT JOIN \"UserTable\" u ON u.id = utb.user_id").
		Where("u.public_id IN ?", publicIds).
		Find(&badgesWithPublicUserIds)
	if err := result.Error; err != nil {
		return nil, exceptions.Badge.NotFound().WithError(err)
	}

	if requiredStatic {
		publicIdToIndexMap := make(map[string]int)
		for index, publicId := range publicIds {
			publicIdToIndexMap[publicId] = index
		}

		publicBadges := make([]*gqlmodels.PublicBadge, len(badgesWithPublicUserIds))
		for _, badgeWithPublicUserId := range badgesWithPublicUserIds {
			index := publicIdToIndexMap[badgeWithPublicUserId.PublicId]
			publicBadges[index] = badgeWithPublicUserId.Badge.ToPublicBadge()
		}

		return publicBadges, nil
	}

	publicBadges := make([]*gqlmodels.PublicBadge, 0)
	for _, badgeWithPublicUserId := range badgesWithPublicUserIds {
		publicBadges = append(publicBadges, badgeWithPublicUserId.Badge.ToPublicBadge())
	}

	return publicBadges, nil
}
