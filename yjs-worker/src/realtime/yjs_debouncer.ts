import { randomUUID } from "node:crypto";
import type WebSocket from "ws";
import { WebSocket as WebSocketState } from "ws";
import * as Y from "yjs";

import {
  YjsPersistenceBatchDebounceMilliseconds,
  YjsPersistenceBatchMaximumPayloadBytes,
  YjsPersistenceBatchMaximumUpdateCount,
  YjsPersistenceBatchMaximumWaitMilliseconds,
  YjsPersistenceBatchRetryMilliseconds,
} from "../constants/yjs_persistence_batch.js";
import { InternalFrameType } from "../types/internal_frame_type.js";
import type { Room } from "../types/room.js";
import {
  createYjsPersistenceBatch,
  type InFlightYjsPersistenceBatch,
} from "../types/yjs_persistence_batch.js";
import { YjsPersistenceFailureType } from "../types/yjs_persistence_failure_type.js";
import type { PendingYjsUpdate } from "../types/yjs_update.js";

export class YjsDebouncer {
  private sendInternalFrame!: (
    webSocket: WebSocket,
    type: InternalFrameType,
    connectionId: string,
    connectorChannelId: number,
    blockPackId: string,
    payload?: Buffer
  ) => boolean;
  private resyncRoom!: (room: Room, blockPackId: string) => void;

  bindCallbacks(
    sendInternalFrame: (
      webSocket: WebSocket,
      type: InternalFrameType,
      connectionId: string,
      connectorChannelId: number,
      blockPackId: string,
      payload?: Buffer
    ) => boolean,
    resyncRoom: (room: Room, blockPackId: string) => void
  ): void {
    this.sendInternalFrame = sendInternalFrame;
    this.resyncRoom = resyncRoom;
  }

  scheduleFlush(room: Room, blockPackId: string): void {
    if (room.document === null || room.pendingPersistenceUpdates.length === 0) {
      return;
    }

    if (room.persistenceDebounceTimer !== null) {
      clearTimeout(room.persistenceDebounceTimer);
    }
    room.persistenceDebounceTimer = setTimeout(() => {
      room.persistenceDebounceTimer = null;
      this.flush(room, blockPackId);
    }, YjsPersistenceBatchDebounceMilliseconds);

    if (room.persistenceMaximumWaitTimer !== null) {
      return;
    }

    room.persistenceMaximumWaitTimer = setTimeout(() => {
      room.persistenceMaximumWaitTimer = null;
      this.flush(room, blockPackId);
    }, YjsPersistenceBatchMaximumWaitMilliseconds);
  }

  flush(
    room: Room,
    blockPackId: string,
    webSocket: WebSocket | null = null,
    connectionId: string | null = null,
    connectorChannelId: number | null = null
  ): void {
    if (
      room.document === null ||
      room.inFlightPersistenceBatch !== null ||
      room.pendingPersistenceUpdates.length === 0
    ) {
      return;
    }

    if (room.persistenceDebounceTimer !== null) {
      clearTimeout(room.persistenceDebounceTimer);
      room.persistenceDebounceTimer = null;
    }
    if (room.persistenceMaximumWaitTimer !== null) {
      clearTimeout(room.persistenceMaximumWaitTimer);
      room.persistenceMaximumWaitTimer = null;
    }

    const pendingPersistenceUpdates = room.pendingPersistenceUpdates;
    room.pendingPersistenceUpdates = [];
    room.pendingPersistencePayloadBytes = 0;

    let payload: Buffer;
    try {
      payload = Buffer.from(
        Y.mergeUpdates(
          pendingPersistenceUpdates.map(update => update.frame.payload)
        )
      );
    } catch {
      this.resyncRoom(room, blockPackId);

      return;
    }

    const firstUpdate = pendingPersistenceUpdates[0];
    const originConnectionId = pendingPersistenceUpdates.every(
      update => update.frame.connectionId === firstUpdate.frame.connectionId
    )
      ? firstUpdate.frame.connectionId
      : null;
    room.inFlightPersistenceBatch = {
      persistenceBatchId: randomUUID(),
      originConnectionId,
      payload,
      webSocket: webSocket ?? firstUpdate.webSocket,
      connectionId: connectionId ?? firstUpdate.frame.connectionId,
      connectorChannelId:
        connectorChannelId ?? firstUpdate.frame.connectorChannelId,
      updateCount: pendingPersistenceUpdates.length,
    };

    if (!this.sendInFlight(room, blockPackId)) {
      this.scheduleRetry(room, blockPackId);
    }
  }

  queueUpdate(room: Room, pendingYjsUpdate: PendingYjsUpdate): void {
    if (room.document === null) {
      room.pendingYjsUpdates.push(pendingYjsUpdate);

      return;
    }

    try {
      Y.applyUpdate(room.document, pendingYjsUpdate.frame.payload);
    } catch {
      this.sendInternalFrame(
        pendingYjsUpdate.webSocket,
        InternalFrameType.InternalFrameType_ResyncRequired,
        pendingYjsUpdate.frame.connectionId,
        pendingYjsUpdate.frame.connectorChannelId,
        pendingYjsUpdate.frame.blockPackId
      );

      return;
    }

    room.pendingPersistenceUpdates.push(pendingYjsUpdate);
    room.pendingPersistencePayloadBytes +=
      pendingYjsUpdate.frame.payload.length;
    if (
      room.pendingPersistenceUpdates.length >=
        YjsPersistenceBatchMaximumUpdateCount ||
      room.pendingPersistencePayloadBytes >=
        YjsPersistenceBatchMaximumPayloadBytes
    ) {
      this.flush(room, pendingYjsUpdate.frame.blockPackId);

      return;
    }

    this.scheduleFlush(room, pendingYjsUpdate.frame.blockPackId);
  }

  queueDeferredUpdates(room: Room): void {
    const pendingYjsUpdates = room.pendingYjsUpdates;
    room.pendingYjsUpdates = [];

    for (const pendingYjsUpdate of pendingYjsUpdates) {
      this.queueUpdate(room, pendingYjsUpdate);
    }
  }

  retryInFlight(
    room: Room,
    blockPackId: string,
    webSocket: WebSocket,
    connectionId: string,
    connectorChannelId: number
  ): void {
    if (room.inFlightPersistenceBatch === null) {
      return;
    }

    if (room.persistenceRetryTimer !== null) {
      clearTimeout(room.persistenceRetryTimer);
      room.persistenceRetryTimer = null;
    }

    room.inFlightPersistenceBatch.webSocket = webSocket;
    room.inFlightPersistenceBatch.connectionId = connectionId;
    room.inFlightPersistenceBatch.connectorChannelId = connectorChannelId;
    if (!this.sendInFlight(room, blockPackId)) {
      this.scheduleRetry(room, blockPackId);
    }
  }

  handlePersisted(
    room: Room,
    blockPackId: string,
    connectionId: string,
    connectorChannelId: number,
    updateSequence: number | null
  ): InFlightYjsPersistenceBatch | null {
    if (
      updateSequence === null ||
      room.inFlightPersistenceBatch === null ||
      room.inFlightPersistenceBatch.connectionId !== connectionId ||
      room.inFlightPersistenceBatch.connectorChannelId !== connectorChannelId ||
      updateSequence !== room.lastUpdateSequence + 1
    ) {
      this.resyncRoom(room, blockPackId);

      return null;
    }

    const inFlightPersistenceBatch = room.inFlightPersistenceBatch;
    room.inFlightPersistenceBatch = null;
    if (room.persistenceRetryTimer !== null) {
      clearTimeout(room.persistenceRetryTimer);
      room.persistenceRetryTimer = null;
    }
    room.lastUpdateSequence = updateSequence;
    room.dirtyUpdateCount += inFlightPersistenceBatch.updateCount;

    return inFlightPersistenceBatch;
  }

  handlePersistenceFailure(
    room: Room,
    blockPackId: string,
    payload: Buffer
  ): void {
    if (room.inFlightPersistenceBatch === null) {
      this.resyncRoom(room, blockPackId);

      return;
    }

    if (
      payload[0] ===
      YjsPersistenceFailureType.YjsPersistenceFailureType_Retryable
    ) {
      this.scheduleRetry(room, blockPackId);

      return;
    }

    if (
      payload[0] ===
      YjsPersistenceFailureType.YjsPersistenceFailureType_Terminal
    ) {
      this.resyncRoom(room, blockPackId);
    }
  }

  private sendInFlight(room: Room, blockPackId: string): boolean {
    const batch = room.inFlightPersistenceBatch;
    if (batch === null) {
      return false;
    }

    return this.sendInternalFrame(
      batch.webSocket,
      InternalFrameType.InternalFrameType_AppendYjsUpdateBatch,
      batch.connectionId,
      batch.connectorChannelId,
      blockPackId,
      createYjsPersistenceBatch(
        batch.persistenceBatchId,
        batch.originConnectionId,
        batch.payload
      )
    );
  }

  private scheduleRetry(room: Room, blockPackId: string): void {
    if (
      room.inFlightPersistenceBatch === null ||
      room.persistenceRetryTimer !== null ||
      room.inFlightPersistenceBatch.webSocket.readyState !== WebSocketState.OPEN
    ) {
      return;
    }

    room.persistenceRetryTimer = setTimeout(() => {
      room.persistenceRetryTimer = null;

      if (!this.sendInFlight(room, blockPackId)) {
        this.scheduleRetry(room, blockPackId);
      }
    }, YjsPersistenceBatchRetryMilliseconds);
  }
}
