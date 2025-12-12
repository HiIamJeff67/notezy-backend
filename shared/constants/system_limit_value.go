package constants

import (
	"time"

	types "notezy-backend/shared/types"
)

const (
	ExtraCapacity = 2
)

const (
	MinIntervalTimeOfLastRequest = time.Microsecond
)

const (
	MaxDatabaseUpdateParameters = 65535
)

const (
	DefaultSearchLimit = 10
	MaxSearchLimit     = 100
)

const (
	MaxNonVideoFileSize        types.ByteType = 10 * types.MB
	MaxInMemoryStorageFileSize types.ByteType = 10 * types.MB
	MaxS3StorageFileSize       types.ByteType = 10 * types.MB
)

const (
	MaxUserAgentLength           int = 2048
	MaxURLLength                 int = 2048
	MinPasswordLength            int = 8
	MaxPasswordLength            int = 1024
	MaxHexCodeLength             int = 20
	MaxProgrammingLanguageLength int = 50
)

const (
	// the max batch size of the PostgreSQL and Gorm is limited by
	// the formula: Max Batch Size = 65535 / number of columns in the target table
	MaxBatchCreateBlockSize int = 3000
)
