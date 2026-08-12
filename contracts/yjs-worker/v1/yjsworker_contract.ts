// Versioned TypeScript contract corresponding to constants.go in this package.
// The planned contract generator will own this file once it is introduced.

export const YjsBlockPackRoomPrefix = "block-pack";
export const YjsBlockPackFragmentName = "document-store";
export const YjsBlockPackSchemaId = "notezy.blocknote";
export const YjsBlockPackSchemaVersion = 1;
export const BlockPackDocumentQuotaPolicyVersion = 1;
export const YjsCompactionUpdateThreshold = 500;
export const YjsMaintenanceMaximumPayloadBytes = 64 * 1024 * 1024;
export const YjsDocumentMaximumLoadPayloadBytes = 64 * 1024 * 1024;
export const InternalFrameHeaderSize = 39;

export type BlockPackQuotaPolicy = {
  version: number;
  maximumBlockCount: number;
};

export const YjsWorkerCoreCommandTopic = "notezy.adapters.core.command.v1";
export const CoreYjsWorkerReplyTopic = "notezy.core.adapters.reply.v1";
export const CoreYjsWorkerMaintenanceCommandTopic =
  "notezy.core.adapters.maintenance-command.v1";
export const YjsWorkerCoreMaintenanceResultTopic =
  "notezy.adapters.core.maintenance-result.v1";
