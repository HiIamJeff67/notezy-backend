package options

import (
	"gorm.io/gorm"

	models "notezy-backend/app/models"
	types "notezy-backend/shared/types"
)

type RepositoryOptionFields struct {
	DB                  *gorm.DB
	OnlyDeleted         types.Ternary
	SkipPermissionCheck bool
	BatchSize           int
}

type RepositoryOptions func(*RepositoryOptionFields)

func WithDB(db *gorm.DB) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.DB = db
	}
}

func WithOnlyDeleted(onlyDeleted types.Ternary) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.OnlyDeleted = onlyDeleted
	}
}

func WithSkipPermissionCheck() RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.SkipPermissionCheck = true
	}
}

func WithBatchSize(batchSize int) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.BatchSize = batchSize
	}
}

func GetDefaultOptions() RepositoryOptionFields {
	return RepositoryOptionFields{
		DB:                  models.NotezyDB,
		OnlyDeleted:         types.Ternary_Neutral,
		SkipPermissionCheck: false,
		BatchSize:           1000,
	}
}

func ParseRepositoryOptions(opts ...RepositoryOptions) RepositoryOptionFields {
	ros := GetDefaultOptions()
	for _, opt := range opts {
		opt(&ros)
	}
	return ros
}
