import type { Doc } from "yjs";
import type WebSocket from "ws";

export type RoomSubscriber = {
  webSocket: WebSocket;
  connectionId: string;
  connectorChannelId: number;
};

export type Room = {
  document: Doc;
  dirtyUpdateCount: number;
  lastActiveAt: Date;
  subscribers: Map<string, RoomSubscriber>;
};
