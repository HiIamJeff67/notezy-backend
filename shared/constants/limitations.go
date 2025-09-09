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
	CheckPointPerShelfTreeTraverse  = 1e+3
	MaxShelfTreeTraverseTimeout     = 8 * time.Second
	MaxNumOfShelfTreeTraversedNodes = 1e+5
	MaxShelfTreeWidth               = 1e+5
	MaxShelfTreeDepth               = 100 // note that the maximum width is bounded by the msgpack encoding algorithm

	MaxShelvesToSynchronize = 20
)

const (
	MaxNonVideoFileSize        int64 = 5 * MB
	MaxInMemoryStorageFileSize int64 = 5 * MB
	MaxS3StorageFileSize       int64 = 5 * MB
)
