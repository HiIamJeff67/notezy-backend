package services

import "gorm.io/gorm"

/* ============================== Interface & Instance ============================== */

type ThemeServiceInterface interface {
}

type ThemeService struct {
	db *gorm.DB
}

func NewThemeService(db *gorm.DB) ThemeServiceInterface {
	return &ThemeService{
		db: db,
	}
}

/* ============================== Services for Themes ============================== */

func (s *ThemeService) GetMyThemes() {}
