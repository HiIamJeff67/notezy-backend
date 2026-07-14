import type WebSocket from "ws";

import { YjsRoomIdleEvictionMilliseconds } from "../constants/eviction.js";
import type { Room } from "../types/room.js";

export class RoomRegistry {
  private readonly rooms = new Map<string, Room>();

  getOrCreate(blockPackId: string): Room {
    const existingRoom = this.rooms.get(blockPackId);
    if (existingRoom !== undefined) {
      existingRoom.lastActiveAt = new Date();

      return existingRoom;
    }

    const room: Room = {
      document: null,
      dirtyUpdateCount: 0,
      lastActiveAt: new Date(),
      subscribers: new Map(),
      isLoading: false,
      lastUpdateSequence: 0,
      compactedUntilSequence: 0,
      projectedUntilSequence: -1,
      pendingYjsUpdates: [],
      pendingPersistenceUpdates: [],
      pendingPersistencePayloadBytes: 0,
      idleEvictionTimer: null,
      persistenceDebounceTimer: null,
      persistenceMaximumWaitTimer: null,
      persistenceRetryTimer: null,
      inFlightPersistenceBatch: null,
      isCompacting: false,
      projectionTimer: null,
      inFlightProjection: null,
    };
    this.rooms.set(blockPackId, room);

    return room;
  }

  attach(
    blockPackId: string,
    webSocket: WebSocket,
    connectionId: string,
    connectorChannelId: number
  ): Room {
    const room = this.getOrCreate(blockPackId);
    this.cancelEviction(room);
    room.subscribers.set(
      this.getSubscriberKey(connectionId, connectorChannelId),
      {
        webSocket,
        connectionId,
        connectorChannelId,
        isReady: false,
      }
    );

    return room;
  }

  detach(
    blockPackId: string,
    connectionId: string,
    connectorChannelId: number
  ): Room | undefined {
    const room = this.rooms.get(blockPackId);
    if (room === undefined) {
      return undefined;
    }

    room.subscribers.delete(
      this.getSubscriberKey(connectionId, connectorChannelId)
    );
    room.lastActiveAt = new Date();

    return room;
  }

  detachAll(webSocket: WebSocket): Array<[string, Room]> {
    const detachedRooms: Array<[string, Room]> = [];
    for (const [blockPackId, room] of this.rooms) {
      let didDetach = false;
      for (const [key, subscriber] of room.subscribers) {
        if (subscriber.webSocket === webSocket) {
          room.subscribers.delete(key);
          didDetach = true;
        }
      }

      if (didDetach) {
        room.lastActiveAt = new Date();
        detachedRooms.push([blockPackId, room]);
      }
    }

    return detachedRooms;
  }

  getSubscriber(
    blockPackId: string,
    connectionId: string,
    connectorChannelId: number
  ): Room | undefined {
    const room = this.rooms.get(blockPackId);
    if (
      room === undefined ||
      !room.subscribers.has(
        this.getSubscriberKey(connectionId, connectorChannelId)
      )
    ) {
      return undefined;
    }

    room.lastActiveAt = new Date();

    return room;
  }

  get(blockPackId: string): Room | undefined {
    return this.rooms.get(blockPackId);
  }

  private isEvictable(blockPackId: string, room: Room): boolean {
    return (
      this.rooms.get(blockPackId) === room &&
      room.subscribers.size === 0 &&
      !room.isLoading &&
      room.pendingYjsUpdates.length === 0 &&
      room.pendingPersistenceUpdates.length === 0 &&
      room.inFlightPersistenceBatch === null &&
      !room.isCompacting &&
      room.persistenceDebounceTimer === null &&
      room.persistenceMaximumWaitTimer === null &&
      room.persistenceRetryTimer === null &&
      room.inFlightProjection === null
    );
  }

  scheduleRoomEviction(blockPackId: string): void {
    const room = this.rooms.get(blockPackId);
    if (room === undefined) {
      return;
    }

    if (room.subscribers.size > 0) {
      this.cancelEviction(room);

      return;
    }

    if (room.projectionTimer !== null) {
      clearTimeout(room.projectionTimer);
      room.projectionTimer = null;
    }

    if (
      !this.isEvictable(blockPackId, room) ||
      room.idleEvictionTimer !== null
    ) {
      return;
    }

    room.idleEvictionTimer = setTimeout(() => {
      room.idleEvictionTimer = null;
      if (this.isEvictable(blockPackId, room)) {
        this.evict(blockPackId, room);
      }
    }, YjsRoomIdleEvictionMilliseconds);
  }

  cancelEviction(room: Room): void {
    if (room.idleEvictionTimer === null) {
      return;
    }

    clearTimeout(room.idleEvictionTimer);
    room.idleEvictionTimer = null;
  }

  evict(blockPackId: string, room: Room): boolean {
    if (!this.isEvictable(blockPackId, room)) {
      return false;
    }

    this.cancelEviction(room);
    if (room.persistenceDebounceTimer !== null) {
      clearTimeout(room.persistenceDebounceTimer);
    }
    if (room.persistenceMaximumWaitTimer !== null) {
      clearTimeout(room.persistenceMaximumWaitTimer);
    }
    if (room.persistenceRetryTimer !== null) {
      clearTimeout(room.persistenceRetryTimer);
    }
    if (room.projectionTimer !== null) {
      clearTimeout(room.projectionTimer);
    }

    room.document?.destroy();
    room.document = null;
    room.pendingYjsUpdates = [];
    room.pendingPersistenceUpdates = [];
    room.pendingPersistencePayloadBytes = 0;
    room.persistenceDebounceTimer = null;
    room.persistenceMaximumWaitTimer = null;
    room.persistenceRetryTimer = null;
    room.projectionTimer = null;
    this.rooms.delete(blockPackId);

    return true;
  }

  private getSubscriberKey(
    connectionId: string,
    connectorChannelId: number
  ): string {
    return `${connectionId}:${connectorChannelId}`;
  }

  get size(): number {
    return this.rooms.size;
  }

  entries(): IterableIterator<[string, Room]> {
    return this.rooms.entries();
  }
}
