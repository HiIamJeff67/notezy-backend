package constants

import "time"

const (
	YjsBlockPackRoomPrefix                    = "block-pack"
	YjsBlockPackFragmentName                  = "document-store"
	YjsBlockPackSchemaId                      = "notezy.blocknote"
	YjsBlockPackSchemaVersion                 = 1
	YjsCompactionUpdateThreshold        int64 = 500
	YjsMaintenanceScanInterval                = 5 * time.Minute
	YjsMaintenanceMaxDocumentsPerRun          = 20
	YjsCompactedUpdateRetention               = 24 * time.Hour
	YjsCleanupMaxUpdatesPerRun                = 1_000
	YjsMaintenanceWorkerRequestTimeout        = 30 * time.Second
	YjsMaintenanceWorkerMaxPayloadBytes       = 64 * 1024 * 1024
)
