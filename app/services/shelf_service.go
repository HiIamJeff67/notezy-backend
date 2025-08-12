package services

import "gorm.io/gorm"

/* ============================== Interface & Instance ============================== */

type ShelfServiceInterface interface{}

type ShelfService struct {
	db *gorm.DB
}

func NewShelfService(db *gorm.DB) ShelfServiceInterface {
	return &ShelfService{
		db: db,
	}
}

/* ============================== Service Methods for Shelves ============================== */
