package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	platformpostgres "github.com/HiIamJeff67/notezy-backend/shared/platform/postgres"

	schemas "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/schemas"
)

func Connect(config platformpostgres.Config) (*gorm.DB, error) {
	db, err := platformpostgres.Connect(config)
	if err != nil {
		return nil, err
	}
	if err := platformpostgres.MigrateTablesToDatabase(db, schemas.MigratingTables); err != nil {
		_ = platformpostgres.Disconnect(db)
		return nil, fmt.Errorf("migrate notification database: %w", err)
	}

	return db, nil
}

func Disconnect(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if err := platformpostgres.Disconnect(db); err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(context.Background(), err, "Failed to close Notification database")
		}
		return err
	}

	return nil
}
