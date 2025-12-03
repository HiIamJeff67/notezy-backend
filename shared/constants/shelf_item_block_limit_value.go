package constants

/* ============================== Limitation of Shelf and Items ============================== */

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

	MaxShelfNameLength int = 128
	MaxItemNameLength  int = 128

	PeekFileSize            int64 = 256 * Byte
	MaxTextbookFileSize     int64 = 5 * MB
	MaxNotebookFileSize     int64 = 5 * MB
	MaxLearningCardFileSize int64 = 1 * MB
	MaxWorkFlowFileSize     int64 = 10 * MB
)

/* ============================== Limitation of BlockContent or Props JSON ============================== */

const (
	MaxFileBlockCaptionLength = 512
	MaxFileBlockNameLength    = 256
	MaxHeadingLevel           = 6
)
