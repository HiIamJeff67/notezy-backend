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

	PeekFileSize            int64 = 256 * Byte
	MaxTextbookFileSize     int64 = 5 * MB
	MaxNotebookFileSize     int64 = 5 * MB
	MaxLearningCardFileSize int64 = 1 * MB
	MaxWorkFlowFileSize     int64 = 10 * MB
)

const (
	MaxNonVideoFileSize        int64 = 10 * MB
	MaxInMemoryStorageFileSize int64 = 10 * MB
	MaxS3StorageFileSize       int64 = 10 * MB
)
