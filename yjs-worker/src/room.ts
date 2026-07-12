import type { Doc } from "yjs";

export type Room = {
  document: Doc;
  dirtyUpdateCount: number;
  lastActiveAt: Date;
};
