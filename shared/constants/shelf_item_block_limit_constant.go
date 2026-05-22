package constants

import types "notezy-backend/shared/types"

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

	PeekFileSize            types.ByteType = 256 * types.Byte
	MaxMaterialTextFileSize  types.ByteType = 5 * types.MB
	MaxMaterialImageFileSize types.ByteType = 20 * types.MB
	MaxMaterialVideoFileSize types.ByteType = 100 * types.MB
	MaxMaterialAudioFileSize types.ByteType = 20 * types.MB
)
