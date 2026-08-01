package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
)

func ConnectionString(config configs.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.DBName,
		config.Password,
	)
}

func Connect(config configs.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(ConnectionString(config)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func Disconnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
