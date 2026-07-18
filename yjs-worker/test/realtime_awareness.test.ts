import assert from "node:assert/strict";
import { EventEmitter } from "node:events";
import test from "node:test";

import { WebSocket } from "ws";
import {
  Awareness,
  applyAwarenessUpdate,
  encodeAwarenessUpdate,
} from "y-protocols/awareness";
import * as Y from "yjs";

import { RealtimeGateway } from "../src/realtime/gateway.js";
import { RoomRegistry } from "../src/realtime/room_registry.js";
import { YjsCompactionService } from "../src/services/yjs_compaction_service.js";
import { Telemetry } from "../src/telemetry.js";
import {
  createInternalFrame,
  parseInternalFrame,
} from "../src/types/internal_frame.js";
import { InternalFrameType } from "../src/types/internal_frame_type.js";

class TestWebSocket extends EventEmitter {
  readonly sentFrames: Buffer[] = [];
  readyState = WebSocket.OPEN;

  send(frame: Buffer): void {
    this.sentFrames.push(frame);
  }

  close(): void {
    this.readyState = WebSocket.CLOSED;
  }
}

const telemetry = Telemetry.initialize();

test.after(async () => {
  await telemetry.shutdown();
});

test("RealtimeGateway isolates awareness client IDs and removes them on detach", () => {
  const roomRegistry = new RoomRegistry(telemetry);
  const gateway = new RealtimeGateway(
    roomRegistry,
    new YjsCompactionService(telemetry),
    telemetry
  );
  const firstSocket = new TestWebSocket();
  const secondSocket = new TestWebSocket();
  const firstConnectionId = "0f5e3ec2-9211-4f62-8e57-25fd5d8104ec";
  const secondConnectionId = "97553d3a-c805-4372-b624-c7c30aad5f10";
  const blockPackId = "4bb4cc0e-44e5-4c2f-a1a1-26e7b5150da6";

  gateway.handleConnection(firstSocket as unknown as WebSocket);
  gateway.handleConnection(secondSocket as unknown as WebSocket);

  const room = roomRegistry.attach(
    blockPackId,
    firstSocket as unknown as WebSocket,
    firstConnectionId,
    1
  );
  roomRegistry.attach(
    blockPackId,
    secondSocket as unknown as WebSocket,
    secondConnectionId,
    2
  );
  roomRegistry.initializeAwareness(room, new Y.Doc());
  for (const subscriber of room.subscribers.values()) {
    subscriber.isReady = true;
  }

  const sourceDocument = new Y.Doc();
  const sourceAwareness = new Awareness(sourceDocument);
  sourceAwareness.setLocalState({ user: "first" });
  const clientId = sourceAwareness.clientID;
  const update = Buffer.from(
    encodeAwarenessUpdate(sourceAwareness, [clientId])
  );

  firstSocket.emit(
    "message",
    createInternalFrame(
      InternalFrameType.InternalFrameType_Awareness,
      firstConnectionId,
      1,
      blockPackId,
      update
    ),
    true
  );

  assert.deepEqual(room.awareness?.getStates().get(clientId), {
    user: "first",
  });
  const awarenessFrame = parseInternalFrame(secondSocket.sentFrames.at(-1)!);
  assert.notEqual(awarenessFrame, null);
  assert.equal(
    awarenessFrame.type,
    InternalFrameType.InternalFrameType_Awareness
  );

  secondSocket.emit(
    "message",
    createInternalFrame(
      InternalFrameType.InternalFrameType_Awareness,
      secondConnectionId,
      2,
      blockPackId,
      update
    ),
    true
  );

  const resyncFrame = parseInternalFrame(secondSocket.sentFrames.at(-1)!);
  assert.notEqual(resyncFrame, null);
  assert.equal(
    resyncFrame.type,
    InternalFrameType.InternalFrameType_ResyncRequired
  );
  assert.deepEqual(room.awareness?.getStates().get(clientId), {
    user: "first",
  });

  const observedDocument = new Y.Doc();
  const observedAwareness = new Awareness(observedDocument);
  observedAwareness.setLocalState(null);
  applyAwarenessUpdate(observedAwareness, awarenessFrame.payload, gateway);
  assert.deepEqual(observedAwareness.getStates().get(clientId), {
    user: "first",
  });

  firstSocket.emit(
    "message",
    createInternalFrame(
      InternalFrameType.InternalFrameType_Detach,
      firstConnectionId,
      1,
      blockPackId
    ),
    true
  );

  assert.equal(room.awareness?.getStates().has(clientId), false);
  const removalFrame = parseInternalFrame(secondSocket.sentFrames.at(-1)!);
  assert.notEqual(removalFrame, null);
  assert.equal(
    removalFrame.type,
    InternalFrameType.InternalFrameType_Awareness
  );
  applyAwarenessUpdate(observedAwareness, removalFrame.payload, gateway);
  assert.equal(observedAwareness.getStates().has(clientId), false);

  sourceAwareness.destroy();
  sourceDocument.destroy();
  observedAwareness.destroy();
  observedDocument.destroy();
  roomRegistry.clearAwareness(room);
  room.document?.destroy();
});
