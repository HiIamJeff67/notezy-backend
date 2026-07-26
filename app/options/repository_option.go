package options

import (
	"gorm.io/gorm"

	models "github.com/HiIamJeff67/notezy-backend/app/models"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
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
	AllowedPermissions   []enums.AccessControlPermission
	OnlyDeleted          types.Ternary
	LockingStrength      *string
	BatchSize            int
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

func WithAllowedPermissions(allowedPermissions []enums.AccessControlPermission) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.AllowedPermissions = append([]enums.AccessControlPermission{}, allowedPermissions...)
	}
}

func WithOnlyDeleted(onlyDeleted types.Ternary) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.OnlyDeleted = onlyDeleted
	}
}

func WithLockingStrength(lockingStrength string) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.LockingStrength = &lockingStrength
	}
}

func WithBatchSize(batchSize int) RepositoryOptions {
	return func(ros *RepositoryOptionFields) {
		ros.BatchSize = batchSize
	}
}

func (ros RepositoryOptionFields) HasAllowedPermissions() bool {
	return ros.AllowedPermissions != nil
}

func GetDefaultOptions() RepositoryOptionFields {
	return RepositoryOptionFields{
		DB:                   models.NotezyDB,
		IsTransactionStarted: false,
		AllowedPermissions:   nil,
		OnlyDeleted:          types.Ternary_Neutral,
		LockingStrength:      nil,
		BatchSize:            1000,
	}
}

func ParseRepositoryOptions(opts ...RepositoryOptions) RepositoryOptionFields {
	ros := GetDefaultOptions()
	for _, opt := range opts {
		opt(&ros)
	}
	return ros
}
