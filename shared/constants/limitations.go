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
	MaxSubShelvesOfRootShelf int32 = 1e+2 // max number of the sub folders
	MaxContentOfRootShelf    int32 = 1e+2 // max number of all types of content under a root shelf
	MaxMaterialsOfRootShelf  int32 = 1e+2 // max number of materials(files)
	MaxBlockPackOfRootShelf  int32 = 1e+2 // max number of block packs

	// limitation of a sub shelf
	MaxSubShelvesOfSubShelf int32 = 1e+2 // max number of sub folders
	MaxContentOfSubShelf    int32 = 1e+2 // max number of all types of content under a sub shelf
	MaxMaterialsOfSubShelf  int32 = 1e+2 // max number of materials(files)
	MaxBlockPackOfSubShelf  int32 = 1e+2 // max number of block packs

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
