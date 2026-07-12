import * as Y from "yjs";
import type WebSocket from "ws";

import type { Room } from "./room.js";

export class RoomRegistry {
  private readonly rooms = new Map<string, Room>();

  getOrCreate(blockPackId: string): Room {
    const existingRoom = this.rooms.get(blockPackId);
    if (existingRoom !== undefined) {
      existingRoom.lastActiveAt = new Date();

      return existingRoom;
    }

    const room: Room = {
      document: new Y.Doc(),
      dirtyUpdateCount: 0,
      lastActiveAt: new Date(),
      subscribers: new Map(),
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

  detachAll(webSocket: WebSocket): void {
    for (const room of this.rooms.values()) {
      for (const [key, subscriber] of room.subscribers) {
        if (subscriber.webSocket === webSocket) {
          room.subscribers.delete(key);
        }
      }
    }
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

  private getSubscriberKey(connectionId: string, connectorChannelId: number): string {
    return `${connectionId}:${connectorChannelId}`;
  }

  get size(): number {
    return this.rooms.size;
  }
}
