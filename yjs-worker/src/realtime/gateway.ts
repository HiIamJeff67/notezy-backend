import { WebSocket } from "ws";
import * as Y from "yjs";
import { YjsCompactionUpdateThreshold } from "../constants/yjs_compaction.js";
import { YjsPersistenceBatchShutdownTimeoutMilliseconds } from "../constants/yjs_persistence_batch.js";
import { YjsCompactionService } from "../services/yjs_compaction_service.js";
import {
  createInternalFrame,
  parseInternalFrame,
} from "../types/internal_frame.js";
import { InternalFrameType } from "../types/internal_frame_type.js";
import type { Room } from "../types/room.js";
import {
  createYjsCompactionResult,
  parseYjsCompactionInput,
} from "../types/yjs_compaction.js";
import {
  parseYjsDocumentState,
  parseYjsUpdateSequence,
} from "../types/yjs_document_state.js";
import { BlockPackProjector } from "./block_pack_projector.js";
import { RoomRegistry } from "./room_registry.js";
import { YjsDebouncer } from "./yjs_debouncer.js";

export class RealtimeGateway {
  private readonly roomRegistry: RoomRegistry;
  private readonly blockPackProjector: BlockPackProjector;
  private readonly yjsCompactionService: YjsCompactionService;
  private readonly webSockets = new Set<WebSocket>();
  private readonly yjsDebouncer: YjsDebouncer;

  constructor(
    roomRegistry: RoomRegistry,
    blockPackProjector: BlockPackProjector,
    yjsCompactionService: YjsCompactionService,
    yjsDebouncer: YjsDebouncer
  ) {
    this.roomRegistry = roomRegistry;
    this.blockPackProjector = blockPackProjector;
    this.yjsCompactionService = yjsCompactionService;
    this.yjsDebouncer = yjsDebouncer;
    this.yjsDebouncer.bindCallbacks(
      this.sendInternalFrame.bind(this),
      this.resyncRoom.bind(this)
    );
  }

  /* ============================== Internal frame delivery ============================== */

  private broadcastInternalFrame(
    room: Room,
    type: InternalFrameType,
    blockPackId: string,
    payload: Buffer
  ): void {
    for (const subscriber of room.subscribers.values()) {
      if (
        !subscriber.isReady &&
        (type === InternalFrameType.InternalFrameType_YjsDocument ||
          type === InternalFrameType.InternalFrameType_Awareness)
      ) {
        continue;
      }

      this.sendInternalFrame(
        subscriber.webSocket,
        type,
        subscriber.connectionId,
        subscriber.connectorChannelId,
        blockPackId,
        payload
      );
    }
  }

  private sendInternalFrame(
    webSocket: WebSocket,
    type: InternalFrameType,
    connectionId: string,
    connectorChannelId: number,
    blockPackId: string,
    payload: Buffer = Buffer.alloc(0)
  ): boolean {
    if (webSocket.readyState !== WebSocket.OPEN) {
      return false;
    }

    webSocket.send(
      createInternalFrame(
        type,
        connectionId,
        connectorChannelId,
        blockPackId,
        payload
      )
    );

    return true;
  }

  private sendRoomInitialState(room: Room, blockPackId: string): void {
    if (room.document === null) {
      return;
    }

    const payload = Buffer.from(Y.encodeStateAsUpdate(room.document));
    for (const subscriber of room.subscribers.values()) {
      if (subscriber.isReady) {
        continue;
      }

      if (
        this.sendInternalFrame(
          subscriber.webSocket,
          InternalFrameType.InternalFrameType_YjsDocument,
          subscriber.connectionId,
          subscriber.connectorChannelId,
          blockPackId,
          payload
        )
      ) {
        subscriber.isReady = true;
      }
    }
  }

  private resyncRoom(room: Room, blockPackId: string): void {
    if (room.projectionTimer !== null) {
      clearTimeout(room.projectionTimer);
    }
    if (room.persistenceDebounceTimer !== null) {
      clearTimeout(room.persistenceDebounceTimer);
    }
    if (room.persistenceMaximumWaitTimer !== null) {
      clearTimeout(room.persistenceMaximumWaitTimer);
    }
    if (room.persistenceRetryTimer !== null) {
      clearTimeout(room.persistenceRetryTimer);
    }

    room.document?.destroy();
    room.document = null;
    room.isLoading = false;
    room.dirtyUpdateCount = 0;
    room.lastUpdateSequence = 0;
    room.compactedUntilSequence = 0;
    room.projectedUntilSequence = -1;
    room.pendingYjsUpdates = [];
    room.pendingPersistenceUpdates = [];
    room.pendingPersistencePayloadBytes = 0;
    room.persistenceDebounceTimer = null;
    room.persistenceMaximumWaitTimer = null;
    room.persistenceRetryTimer = null;
    room.inFlightPersistenceBatch = null;
    room.projectionTimer = null;
    room.inFlightProjection = null;

    this.broadcastInternalFrame(
      room,
      InternalFrameType.InternalFrameType_ResyncRequired,
      blockPackId,
      Buffer.alloc(0)
    );

    this.roomRegistry.scheduleRoomEviction(blockPackId);
  }

  private scheduleBlockProjection(
    room: Room,
    blockPackId: string,
    delayMilliseconds: number = 300
  ): void {
    if (
      room.document === null ||
      room.inFlightPersistenceBatch !== null ||
      room.pendingYjsUpdates.length > 0 ||
      room.pendingPersistenceUpdates.length > 0 ||
      room.inFlightProjection !== null ||
      room.projectionTimer !== null ||
      room.lastUpdateSequence <= room.projectedUntilSequence ||
      room.subscribers.size === 0
    ) {
      return;
    }

    room.projectionTimer = setTimeout(() => {
      room.projectionTimer = null;

      if (
        room.document === null ||
        room.inFlightPersistenceBatch !== null ||
        room.pendingYjsUpdates.length > 0 ||
        room.pendingPersistenceUpdates.length > 0 ||
        room.inFlightProjection !== null ||
        room.lastUpdateSequence <= room.projectedUntilSequence
      ) {
        return;
      }

      const subscriber = room.subscribers.values().next().value;
      if (subscriber === undefined) {
        return;
      }

      const projectedSequence = room.lastUpdateSequence;
      let payload: Buffer;
      try {
        payload = Buffer.from(
          JSON.stringify({
            schemaId: "notezy.blocknote",
            schemaVersion: 1,
            projectedSequence,
            blocks: this.blockPackProjector.projectYjsDocument(room.document),
          })
        );
      } catch (error) {
        console.error("failed to project Yjs document", {
          blockPackId,
          error,
        });
        this.scheduleBlockProjection(room, blockPackId, 1_000);

        return;
      }

      room.inFlightProjection = {
        connectionId: subscriber.connectionId,
        connectorChannelId: subscriber.connectorChannelId,
        projectedSequence,
      };
      if (
        this.sendInternalFrame(
          subscriber.webSocket,
          InternalFrameType.InternalFrameType_ApplyBlockProjection,
          subscriber.connectionId,
          subscriber.connectorChannelId,
          blockPackId,
          payload
        )
      ) {
        return;
      }

      room.inFlightProjection = null;
      this.scheduleBlockProjection(room, blockPackId, 1_000);
    }, delayMilliseconds);
  }

  private requestYjsCompaction(
    room: Room,
    webSocket: WebSocket,
    connectionId: string,
    connectorChannelId: number,
    blockPackId: string
  ): void {
    if (
      room.document === null ||
      room.isCompacting ||
      room.lastUpdateSequence - room.compactedUntilSequence <
        YjsCompactionUpdateThreshold ||
      room.isLoading ||
      room.pendingYjsUpdates.length > 0 ||
      room.pendingPersistenceUpdates.length > 0 ||
      room.inFlightPersistenceBatch !== null ||
      room.persistenceDebounceTimer !== null ||
      room.persistenceMaximumWaitTimer !== null ||
      room.persistenceRetryTimer !== null
    ) {
      return;
    }

    room.isCompacting = this.sendInternalFrame(
      webSocket,
      InternalFrameType.InternalFrameType_LoadCompactableYjsDocument,
      connectionId,
      connectorChannelId,
      blockPackId
    );
  }

  /* ============================== WebSocket connection ============================== */

  handleConnection(webSocket: WebSocket): void {
    this.webSockets.add(webSocket);
    webSocket.on("close", () => {
      this.webSockets.delete(webSocket);
      for (const [blockPackId, room] of this.roomRegistry.detachAll(
        webSocket
      )) {
        if (room.subscribers.size === 0) {
          this.yjsDebouncer.flush(room, blockPackId);
          this.roomRegistry.scheduleRoomEviction(blockPackId);
        }
      }
    });

    webSocket.on("message", (payload, isBinary) => {
      if (!isBinary) {
        webSocket.close(1003, "internal realtime frames must be binary");

        return;
      }

      let framePayload: Buffer;
      if (Buffer.isBuffer(payload)) {
        framePayload = payload;
      } else if (payload instanceof ArrayBuffer) {
        framePayload = Buffer.from(payload);
      } else if (Array.isArray(payload)) {
        framePayload = Buffer.concat(payload);
      } else {
        webSocket.close(1002, "invalid internal realtime frame");

        return;
      }

      const frame = parseInternalFrame(framePayload);
      if (frame === null || frame.version !== 1) {
        webSocket.close(1002, "invalid internal realtime frame");

        return;
      }

      switch (frame.type) {
        case InternalFrameType.InternalFrameType_Attach: {
          const room = this.roomRegistry.attach(
            frame.blockPackId,
            webSocket,
            frame.connectionId,
            frame.connectorChannelId
          );
          if (room.document !== null) {
            if (room.inFlightPersistenceBatch !== null) {
              this.yjsDebouncer.retryInFlight(
                room,
                frame.blockPackId,
                webSocket,
                frame.connectionId,
                frame.connectorChannelId
              );

              return;
            }

            if (room.pendingPersistenceUpdates.length > 0) {
              this.yjsDebouncer.flush(
                room,
                frame.blockPackId,
                webSocket,
                frame.connectionId,
                frame.connectorChannelId
              );

              return;
            }

            this.sendRoomInitialState(room, frame.blockPackId);
            this.scheduleBlockProjection(room, frame.blockPackId);

            return;
          }

          if (room.isLoading) {
            return;
          }

          room.isLoading = true;
          if (
            this.sendInternalFrame(
              webSocket,
              InternalFrameType.InternalFrameType_LoadYjsDocument,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId
            )
          ) {
            return;
          }

          this.resyncRoom(room, frame.blockPackId);

          return;
        }
        case InternalFrameType.InternalFrameType_Detach: {
          const room = this.roomRegistry.detach(
            frame.blockPackId,
            frame.connectionId,
            frame.connectorChannelId
          );
          if (room !== undefined && room.subscribers.size === 0) {
            this.yjsDebouncer.flush(room, frame.blockPackId);
            this.requestYjsCompaction(
              room,
              webSocket,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId
            );
            this.roomRegistry.scheduleRoomEviction(frame.blockPackId);
          }

          return;
        }
        case InternalFrameType.InternalFrameType_YjsDocument: {
          const room = this.roomRegistry.getSubscriber(
            frame.blockPackId,
            frame.connectionId,
            frame.connectorChannelId
          );
          if (room === undefined) {
            this.sendInternalFrame(
              webSocket,
              InternalFrameType.InternalFrameType_ResyncRequired,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId
            );

            return;
          }

          this.yjsDebouncer.queueUpdate(room, { webSocket, frame });

          return;
        }
        case InternalFrameType.InternalFrameType_Awareness: {
          const room = this.roomRegistry.getSubscriber(
            frame.blockPackId,
            frame.connectionId,
            frame.connectorChannelId
          );
          if (room === undefined) {
            this.sendInternalFrame(
              webSocket,
              InternalFrameType.InternalFrameType_ResyncRequired,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId
            );

            return;
          }

          this.broadcastInternalFrame(
            room,
            frame.type,
            frame.blockPackId,
            frame.payload
          );

          return;
        }
        case InternalFrameType.InternalFrameType_YjsDocumentLoaded: {
          const room = this.roomRegistry.get(frame.blockPackId);
          const state = parseYjsDocumentState(frame.payload);
          if (room === undefined || state === null) return;

          try {
            const document = new Y.Doc();
            if (state.snapshot.length > 0) {
              Y.applyUpdate(document, state.snapshot);
            }
            for (const update of state.updates) {
              Y.applyUpdate(document, update.payload);
            }

            room.document = document;
            room.isLoading = false;
            room.lastUpdateSequence = state.lastUpdateSequence;
            room.compactedUntilSequence = state.compactedUntilSequence;
            room.projectedUntilSequence = state.projectedUntilSequence;
            this.yjsDebouncer.queueDeferredUpdates(room);
            if (
              room.inFlightPersistenceBatch === null &&
              room.pendingPersistenceUpdates.length === 0
            ) {
              this.sendRoomInitialState(room, frame.blockPackId);
            }
            this.scheduleBlockProjection(room, frame.blockPackId);
            this.roomRegistry.scheduleRoomEviction(frame.blockPackId);
          } catch {
            this.resyncRoom(room, frame.blockPackId);
          }

          return;
        }
        case InternalFrameType.InternalFrameType_YjsUpdatePersisted: {
          const room = this.roomRegistry.get(frame.blockPackId);
          if (room === undefined) {
            return;
          }

          const inFlightPersistenceBatch = this.yjsDebouncer.handlePersisted(
            room,
            frame.blockPackId,
            frame.connectionId,
            frame.connectorChannelId,
            parseYjsUpdateSequence(frame.payload)
          );
          if (inFlightPersistenceBatch === null) {
            return;
          }

          this.broadcastInternalFrame(
            room,
            InternalFrameType.InternalFrameType_YjsDocument,
            frame.blockPackId,
            inFlightPersistenceBatch.payload
          );
          this.sendRoomInitialState(room, frame.blockPackId);
          this.yjsDebouncer.flush(room, frame.blockPackId);
          this.requestYjsCompaction(
            room,
            inFlightPersistenceBatch.webSocket,
            inFlightPersistenceBatch.connectionId,
            inFlightPersistenceBatch.connectorChannelId,
            frame.blockPackId
          );
          this.scheduleBlockProjection(room, frame.blockPackId);
          this.roomRegistry.scheduleRoomEviction(frame.blockPackId);

          return;
        }
        case InternalFrameType.InternalFrameType_YjsPersistenceFailed: {
          const room = this.roomRegistry.get(frame.blockPackId);
          if (room !== undefined) {
            this.yjsDebouncer.handlePersistenceFailure(
              room,
              frame.blockPackId,
              frame.payload
            );
          }

          return;
        }
        case InternalFrameType.InternalFrameType_CompactableYjsDocumentLoaded: {
          const room = this.roomRegistry.get(frame.blockPackId);
          const input = parseYjsCompactionInput(frame.payload);
          if (input === null || (room !== undefined && !room.isCompacting)) {
            if (room !== undefined) room.isCompacting = false;

            return;
          }

          try {
            const compacted = this.yjsCompactionService.compact(input);

            this.sendInternalFrame(
              webSocket,
              InternalFrameType.InternalFrameType_ApplyCompactedYjsDocument,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId,
              createYjsCompactionResult(
                input,
                compacted.snapshot,
                compacted.stateVector
              )
            );
          } catch {
            if (room !== undefined) room.isCompacting = false;
          }

          return;
        }
        case InternalFrameType.InternalFrameType_YjsDocumentCompacted: {
          const room = this.roomRegistry.get(frame.blockPackId);
          const compactedUntilSequence = parseYjsUpdateSequence(frame.payload);
          if (room === undefined || compactedUntilSequence === null) return;

          room.isCompacting = false;
          room.compactedUntilSequence = compactedUntilSequence;
          this.roomRegistry.scheduleRoomEviction(frame.blockPackId);

          return;
        }
        case InternalFrameType.InternalFrameType_YjsDocumentCompactionFailed: {
          const room = this.roomRegistry.get(frame.blockPackId);
          if (room !== undefined) room.isCompacting = false;

          return;
        }
        case InternalFrameType.InternalFrameType_BlockProjectionApplied: {
          const room = this.roomRegistry.get(frame.blockPackId);
          let projectedUntilSequence: number | null = null;
          try {
            const value: unknown = JSON.parse(frame.payload.toString("utf8"));
            if (
              value !== null &&
              typeof value === "object" &&
              "projectedUntilSequence" in value &&
              typeof value.projectedUntilSequence === "number" &&
              Number.isSafeInteger(value.projectedUntilSequence) &&
              value.projectedUntilSequence >= -1
            ) {
              projectedUntilSequence = value.projectedUntilSequence;
            }
          } catch {}

          if (
            room === undefined ||
            projectedUntilSequence === null ||
            room.inFlightProjection === null ||
            room.inFlightProjection.connectionId !== frame.connectionId ||
            room.inFlightProjection.connectorChannelId !==
              frame.connectorChannelId ||
            projectedUntilSequence <
              room.inFlightProjection.projectedSequence ||
            projectedUntilSequence > room.lastUpdateSequence
          ) {
            if (room !== undefined) {
              this.resyncRoom(room, frame.blockPackId);
            }

            return;
          }

          room.inFlightProjection = null;
          room.projectedUntilSequence = projectedUntilSequence;
          this.scheduleBlockProjection(room, frame.blockPackId);
          this.roomRegistry.scheduleRoomEviction(frame.blockPackId);

          return;
        }
        case InternalFrameType.InternalFrameType_BlockProjectionFailed: {
          const room = this.roomRegistry.get(frame.blockPackId);
          if (
            room === undefined ||
            room.inFlightProjection === null ||
            room.inFlightProjection.connectionId !== frame.connectionId ||
            room.inFlightProjection.connectorChannelId !==
              frame.connectorChannelId
          ) {
            return;
          }

          room.inFlightProjection = null;
          this.scheduleBlockProjection(room, frame.blockPackId, 1_000);
          this.roomRegistry.scheduleRoomEviction(frame.blockPackId);

          return;
        }
        default:
          console.warn(
            "received internal frame before its handler is enabled",
            {
              type: frame.type,
              blockPackId: frame.blockPackId,
            }
          );
      }
    });
  }

  getActiveRoomCount(): number {
    return this.roomRegistry.size;
  }

  async shutdown(): Promise<void> {
    for (const [blockPackId, room] of this.roomRegistry.entries()) {
      this.yjsDebouncer.flush(room, blockPackId);
    }

    const shutdownDeadline =
      Date.now() + YjsPersistenceBatchShutdownTimeoutMilliseconds;
    while (
      Date.now() < shutdownDeadline &&
      [...this.roomRegistry.entries()].some(
        ([, room]) =>
          room.inFlightPersistenceBatch !== null ||
          room.pendingPersistenceUpdates.length > 0
      )
    ) {
      await new Promise(resolve => setTimeout(resolve, 25));
    }

    this.webSockets.forEach(webSocket => {
      webSocket.close(1001, "server shutdown");
    });
  }
}
