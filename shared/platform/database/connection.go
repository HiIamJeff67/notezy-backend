package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionString(config Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Name,
		config.Password,
	)
}

func Connect(config Config) (*gorm.DB, error) {
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
