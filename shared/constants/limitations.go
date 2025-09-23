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
	// limitation of a root shelf
	MaxSubShelvesOfRootShelf int32 = 1e+5 // max number of the sub folders
	MaxMaterialsOfRootShelf  int32 = 1e+5 // max number of files
	// limitation of a sub shelf
	MaxSubShelvesOfSubShelf int32 = 1e+3
	MaxMaterialsOfSubShelf  int32 = 1e+3
)

const (
	MaxNonVideoFileSize        int64 = 5 * MB
	MaxInMemoryStorageFileSize int64 = 5 * MB
	MaxS3StorageFileSize       int64 = 5 * MB
)
