package constants

import "time"

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
	MaxNonVideoFileSize        int64 = 10 * MB
	MaxInMemoryStorageFileSize int64 = 10 * MB
	MaxS3StorageFileSize       int64 = 10 * MB
)

const (
	MaxUserAgentLength           int = 2048
	MaxURLLength                 int = 2048
	MinPasswordLength            int = 8
	MaxPasswordLength            int = 1024
	MaxHexCodeLength             int = 20
	MaxProgrammingLanguageLength int = 50
)
