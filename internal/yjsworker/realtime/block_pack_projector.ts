import type { Block } from "@blocknote/core";
import { yXmlFragmentToBlocks } from "@blocknote/core/yjs";
import type * as Y from "yjs";
import { YjsBlockPackFragmentName } from "../constants/fragment_name.js";
import { notezyBlockNoteEditor } from "../types/blocknote_schema.js";

export class BlockPackProjector {
  projectYjsDocument(document: Y.Doc): Block[] {
    return yXmlFragmentToBlocks(
      notezyBlockNoteEditor,
      document.getXmlFragment(YjsBlockPackFragmentName)
    );
  }
}
