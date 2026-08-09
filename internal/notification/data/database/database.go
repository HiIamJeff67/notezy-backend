package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	schemas "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/schemas"
)

func Connect(config platformdatabase.Config) (*gorm.DB, error) {
	db, err := platformdatabase.Connect(config)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(schemas.MigratingTables...); err != nil {
		_ = platformdatabase.Disconnect(db)
		return nil, fmt.Errorf("migrate notification database: %w", err)
	}

	return db, nil
}

func Disconnect(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if err := platformdatabase.Disconnect(db); err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(context.Background(), err, "Failed to close Notification database")
		}
		return err
	}

	return nil
}
