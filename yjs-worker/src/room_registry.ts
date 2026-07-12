import * as Y from "yjs";
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
    };
    this.rooms.set(blockPackId, room);

    return room;
  }

  get size(): number {
    return this.rooms.size;
  }
}
