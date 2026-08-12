import { blocksToYXmlFragment } from "@blocknote/core/yjs";
import * as Y from "yjs";
import { YjsBlockPackFragmentName } from "../../../contracts/yjs-worker/v1/yjsworker_contract.js";
import {
  type NotezyBlock,
  notezyBlockNoteEditor,
} from "../types/blocknote_schema.js";

export class YjsDocumentInitializationService {
  initialize(blocks: NotezyBlock[]): {
    snapshot: Buffer;
    stateVector: Buffer;
  } {
    const document = new Y.Doc();
    try {
      blocksToYXmlFragment(
        notezyBlockNoteEditor,
        blocks,
        document.getXmlFragment(YjsBlockPackFragmentName)
      );

      return {
        snapshot: Buffer.from(Y.encodeStateAsUpdate(document)),
        stateVector: Buffer.from(Y.encodeStateVector(document)),
      };
    } finally {
      document.destroy();
    }
  }
}
