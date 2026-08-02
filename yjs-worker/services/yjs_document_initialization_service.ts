import { blocksToYXmlFragment } from "@blocknote/core/yjs";
import * as Y from "yjs";

import { YjsBlockPackFragmentName } from "../constants/fragment_name.js";
import {
  type NotezyBlock,
  notezyBlockNoteEditor,
} from "../types/blocknote_schema.js";

export type YjsDocumentInitializationResult = {
  snapshot: Buffer;
  stateVector: Buffer;
};

export class YjsDocumentInitializationService {
  initialize(blocks: NotezyBlock[]): YjsDocumentInitializationResult {
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
