package options

import (
	"gorm.io/gorm"

	models "notezy-backend/app/models"
	types "notezy-backend/shared/types"
)

type RepositoryOptions struct {
	DB                  *gorm.DB
	OnlyDeleted         types.Ternary
	SkipPermissionCheck bool
	BatchSize           int
}

type RepositoryOption func(*RepositoryOptions)

func WithDB(db *gorm.DB) RepositoryOption {
	return func(ros *RepositoryOptions) {
		ros.DB = db
	}
}

func WithOnlyDeleted(onlyDeleted types.Ternary) RepositoryOption {
	return func(ros *RepositoryOptions) {
		ros.OnlyDeleted = onlyDeleted
	}
}

func WithSkipPermissionCheck() RepositoryOption {
	return func(ros *RepositoryOptions) {
		ros.SkipPermissionCheck = true
	}
}

func WithBatchSize(batchSize int) RepositoryOption {
	return func(ros *RepositoryOptions) {
		ros.BatchSize = batchSize
	}
}

func GetDefaultOptions() RepositoryOptions {
	return RepositoryOptions{
		DB:                  models.NotezyDB,
		OnlyDeleted:         types.Ternary_Neutral,
		SkipPermissionCheck: false,
	}
}

func ParseRepositoryOptions(opts ...RepositoryOption) RepositoryOptions {
	ros := GetDefaultOptions()
	for _, opt := range opts {
		opt(&ros)
	}
	return ros
}
