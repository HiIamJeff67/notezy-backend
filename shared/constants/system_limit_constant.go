package constants

import (
	"time"

	types "notezy-backend/shared/types"
)

/* ============================== Email Worker limitations ============================== */

const (
	EmailWorkerManagerTickerDuration = 100 * time.Millisecond
)

/* ============================== Storage limitations ============================== */

const (
	MaxNonVideoFileSize        types.ByteType = 10 * types.MB
	MaxInMemoryStorageFileSize types.ByteType = 10 * types.MB
	MaxS3StorageFileSize       types.ByteType = 10 * types.MB
)

/* ============================== Variable constraints ============================== */

const (
	MaxUserAgentLength           int = 2048
	MaxURLLength                 int = 2048
	MinPasswordLength            int = 8
	MaxPasswordLength            int = 1024
	MaxHexCodeLength             int = 20
	MaxProgrammingLanguageLength int = 50

	MaxRetriesOfGeneratingFakeName = 5
	// make sure the below values are as the same as the constraint in the dto while registering or creating the user
	MaxNameLength = 32
	MinNameLenght = 6
)

/* ============================== Database or orm limitation ============================== */

const (
	// the max batch size of the PostgreSQL and Gorm is limited by
	// the formula: Max Batch Size = 65535 / number of columns in the target table
	MaxBatchCreateBlockSize int = 3000
)

const (
	DefaultSearchLimit = 10
	MaxSearchLimit     = 100
)
