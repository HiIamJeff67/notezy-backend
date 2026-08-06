import type { NotezyBlock } from "./blocknote_schema.js";

export type YjsDocumentInitializationRequest = {
  documents: Array<{
    blocks: NotezyBlock[];
  }>;
};
