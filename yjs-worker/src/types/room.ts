import type { Doc } from "yjs";
import type WebSocket from "ws";

import type { InternalFrame } from "./internal_frame.js";

export type RoomSubscriber = {
  webSocket: WebSocket;
  connectionId: string;
  connectorChannelId: number;
};

export type PendingYjsUpdate = {
  webSocket: WebSocket;
  frame: InternalFrame;
};

export type InFlightProjection = {
  connectionId: string;
  connectorChannelId: number;
  projectedSequence: number;
};

export type Room = {
  document: Doc | null;
  dirtyUpdateCount: number;
  lastActiveAt: Date;
  subscribers: Map<string, RoomSubscriber>;
  isLoading: boolean;
  lastUpdateSequence: number;
  compactedUntilSequence: number;
  projectedUntilSequence: number;
  pendingYjsUpdates: PendingYjsUpdate[];
  inFlightYjsUpdate: PendingYjsUpdate | null;
  projectionTimer: NodeJS.Timeout | null;
  inFlightProjection: InFlightProjection | null;
};
