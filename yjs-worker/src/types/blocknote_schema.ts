import {
  type Block,
  BlockNoteEditor,
  BlockNoteSchema,
  defaultBlockSpecs,
} from "@blocknote/core";

const { divider: _, ...notezyBlockSpecs } = defaultBlockSpecs;

export const notezyBlockNoteSchema = BlockNoteSchema.create({
  blockSpecs: notezyBlockSpecs,
});

export const notezyBlockNoteEditor = BlockNoteEditor.create({
  schema: notezyBlockNoteSchema,
});

export type NotezyBlock = Block<
  typeof notezyBlockNoteSchema.blockSchema,
  typeof notezyBlockNoteSchema.inlineContentSchema,
  typeof notezyBlockNoteSchema.styleSchema
>;
