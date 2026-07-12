export const internalRealtimePath = "/internal/realtime/v1";
export const internalFrameHeaderSize = 39;

export enum InternalFrameType {
  InternalFrameType_Attach = 1,
  InternalFrameType_Detach = 2,
  InternalFrameType_YjsDocument = 3,
  InternalFrameType_Awareness = 4,
  InternalFrameType_ResyncRequired = 5,
  InternalFrameType_PermissionRevoked = 6,
}

export enum InternalChannelType {
  InternalChannelType_BlockPack = 1,
}
