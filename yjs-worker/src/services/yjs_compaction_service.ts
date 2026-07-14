import * as Y from "yjs";
import type {
  YjsCompactionInput,
  YjsCompactionResult,
} from "../types/yjs_compaction.js";
import {
  createYjsCompactionBatchResult,
  parseYjsCompactionBatchInput,
} from "../types/yjs_compaction_batch.js";

export class YjsCompactionService {
  compact(input: YjsCompactionInput): YjsCompactionResult {
    const document = new Y.Doc();
    try {
      if (input.snapshot.length > 0) {
        Y.applyUpdate(document, input.snapshot);
      }
      for (const update of input.updates) {
        Y.applyUpdate(document, update.payload);
      }

      return {
        baseCompactedUntilSequence: input.baseCompactedUntilSequence,
        cutoffSequence: input.cutoffSequence,
        snapshot: Buffer.from(Y.encodeStateAsUpdate(document)),
        stateVector: Buffer.from(Y.encodeStateVector(document)),
      };
    } finally {
      document.destroy();
    }
  }

  compactBatch(payload: Buffer): Buffer {
    if (payload.length < 4) {
      throw new Error("invalid yjs compaction batch");
    }

    const inputCount = payload.readUInt32BE(0);
    if (inputCount === 0) {
      throw new Error("empty yjs compaction batch");
    }

    const outputPayloads: Buffer[] = [];
    const blockPackIds = new Set<string>();
    let offset = 4;
    for (let index = 0; index < inputCount; index += 1) {
      if (payload.length - offset < 4) {
        throw new Error("invalid yjs compaction batch");
      }

      const inputLength = payload.readUInt32BE(offset);
      offset += 4;
      if (inputLength > payload.length - offset) {
        throw new Error("invalid yjs compaction batch");
      }

      const batchInput = parseYjsCompactionBatchInput(
        payload.subarray(offset, offset + inputLength)
      );
      if (batchInput === null || blockPackIds.has(batchInput.blockPackId)) {
        throw new Error("invalid yjs compaction batch");
      }
      blockPackIds.add(batchInput.blockPackId);
      offset += inputLength;

      const result = this.compact(batchInput.input);
      outputPayloads.push(
        createYjsCompactionBatchResult(
          batchInput.blockPackId,
          batchInput.input,
          result.snapshot,
          result.stateVector
        )
      );
    }
    if (offset !== payload.length) {
      throw new Error("invalid yjs compaction batch");
    }

    const result = Buffer.alloc(
      4 + outputPayloads.reduce((size, item) => size + 4 + item.length, 0)
    );
    result.writeUInt32BE(outputPayloads.length, 0);

    offset = 4;
    for (const outputPayload of outputPayloads) {
      result.writeUInt32BE(outputPayload.length, offset);
      offset += 4;
      outputPayload.copy(result, offset);
      offset += outputPayload.length;
    }

    return result;
  }
}
