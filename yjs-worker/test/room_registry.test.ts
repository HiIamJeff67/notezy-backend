import assert from "node:assert/strict";
import test from "node:test";

import type WebSocket from "ws";
import * as Y from "yjs";

import { RoomRegistry } from "../src/room_registry.js";

test("RoomRegistry tracks subscribers before a Yjs document is materialized", () => {
  const registry = new RoomRegistry();
  const blockPackId = "7bc6ae1a-b1b3-47a7-9fab-42f34f48f7ca";
  const connectionId = "a7577a40-a86d-4fa9-9233-49c0b3e80385";
  const room = registry.attach(blockPackId, {} as WebSocket, connectionId, 1);
  assert.equal(registry.size, 1);
  assert.equal(room.subscribers.size, 1);
  assert.equal(room.document, null);
  assert.equal(registry.getSubscriber(blockPackId, connectionId, 1), room);

  const document = new Y.Doc();
  document.getMap("document-store").set("title", "Notezy");
  room.document = document;

  assert.equal(room.document.getMap("document-store").get("title"), "Notezy");

  registry.detach(blockPackId, connectionId, 1);

  assert.equal(room.subscribers.size, 0);
  assert.equal(registry.getSubscriber(blockPackId, connectionId, 1), undefined);
});
