package adapterscontract

const (
	YjsBlockPackRoomPrefix                    = "block-pack"
	YjsBlockPackFragmentName                  = "document-store"
	YjsBlockPackSchemaId                      = "notegic.blocknote"
	YjsBlockPackSchemaVersion                 = 1
	BlockPackDocumentQuotaPolicyVersion       = 1
	YjsCompactionUpdateThreshold        int64 = 500
	YjsMaintenanceMaximumPayloadBytes         = 64 * 1024 * 1024
	YjsDocumentMaximumLoadPayloadBytes        = 64 * 1024 * 1024
	InternalFrameHeaderSize             int64 = 39
)
