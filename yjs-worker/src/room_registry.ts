import type WebSocket from "ws";

import type { Room } from "./types/room.js";

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
      inFlightYjsUpdate: null,
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
    connectorChannelId: number,
  ): Room {
    const room = this.getOrCreate(blockPackId);
    room.subscribers.set(this.getSubscriberKey(connectionId, connectorChannelId), {
      webSocket,
      connectionId,
      connectorChannelId,
    });

    return room;
  }

  detach(blockPackId: string, connectionId: string, connectorChannelId: number): void {
    const room = this.rooms.get(blockPackId);
    if (room === undefined) {
      return;
    }

    room.subscribers.delete(this.getSubscriberKey(connectionId, connectorChannelId));
    room.lastActiveAt = new Date();
  }

  detachAll(webSocket: WebSocket): Array<{ blockPackId: string; room: Room }> {
    const roomsWithPendingUpdates: Array<{ blockPackId: string; room: Room }> = [];

    for (const [blockPackId, room] of this.rooms) {
      for (const [key, subscriber] of room.subscribers) {
        if (subscriber.webSocket === webSocket) {
          room.subscribers.delete(key);
        }
      }

      if (
        room.inFlightYjsUpdate?.webSocket === webSocket ||
        room.pendingYjsUpdates.some((pendingYjsUpdate) => pendingYjsUpdate.webSocket === webSocket)
      ) {
        roomsWithPendingUpdates.push({ blockPackId, room });
      }
    }

    return roomsWithPendingUpdates;
  }

  getSubscriber(
    blockPackId: string,
    connectionId: string,
    connectorChannelId: number,
  ): Room | undefined {
    const room = this.rooms.get(blockPackId);
    if (room === undefined || !room.subscribers.has(this.getSubscriberKey(connectionId, connectorChannelId))) {
      return undefined;
    }

    room.lastActiveAt = new Date();

    return room;
  }

  get(blockPackId: string): Room | undefined {
    return this.rooms.get(blockPackId);
  }

  private getSubscriberKey(connectionId: string, connectorChannelId: number): string {
    return `${connectionId}:${connectorChannelId}`;
  }

  get size(): number {
    return this.rooms.size;
  }
}
