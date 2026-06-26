package options

import (
	"gorm.io/gorm"

	models "github.com/HiIamJeff67/notezy-backend/app/models"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

const (
	LockingStrengthUpdate      = "UPDATE"
	LockingStrengthNoKeyUpdate = "NO KEY UPDATE"
	LockingStrengthShare       = "SHARE"
)

type RepositoryOptionFields struct {
	DB                   *gorm.DB
	IsTransactionStarted bool
	OnlyDeleted          types.Ternary
	SkipPermissionCheck  bool
	BatchSize            int
	LockingStrength      *string
}

type RepositoryOptions func(*RepositoryOptionFields)

func WithDB(db *gorm.DB) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.DB = db
	}
}

func WithIsTransactionStarted(isTransactionStarted bool) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.IsTransactionStarted = isTransactionStarted
	}
}

func WithTransactionDB(db *gorm.DB) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.DB = db
		ros.IsTransactionStarted = true
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

func WithLockingStrength(lockingStrength string) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.LockingStrength = &lockingStrength
	}
}

func GetDefaultOptions() RepositoryOptionFields {
	return RepositoryOptionFields{
		DB:                   models.NotezyDB,
		OnlyDeleted:          types.Ternary_Neutral,
		SkipPermissionCheck:  false,
		BatchSize:            1000,
		IsTransactionStarted: false,
		LockingStrength:      nil,
	}
}

func ParseRepositoryOptions(opts ...RepositoryOptions) RepositoryOptionFields {
	ros := GetDefaultOptions()
	for _, opt := range opts {
		opt(&ros)
	}
	return ros
}
