import type { Block } from "@blocknote/core";
import { yXmlFragmentToBlocks } from "@blocknote/core/yjs";
import type * as Y from "yjs";

import { notezyBlockNoteEditor } from "./blocknote_schema.js";
import { YjsBlockPackFragmentName } from "./constants/fragment_name.js";

export class BlockNoteProjector {
  projectYjsDocument(document: Y.Doc): Block[] {
    return yXmlFragmentToBlocks(
      notezyBlockNoteEditor,
      document.getXmlFragment(YjsBlockPackFragmentName),
    );
  }
}
