import {
  type Block,
  BlockNoteEditor,
  BlockNoteSchema,
  defaultBlockSpecs,
} from "@blocknote/core";

const { divider: _, ...notegicBlockSpecs } = defaultBlockSpecs;

export const notegicBlockNoteSchema = BlockNoteSchema.create({
  blockSpecs: notegicBlockSpecs,
});

export const notegicBlockNoteEditor = BlockNoteEditor.create({
  schema: notegicBlockNoteSchema,
});

export type NotegicBlock = Block<
  typeof notegicBlockNoteSchema.blockSchema,
  typeof notegicBlockNoteSchema.inlineContentSchema,
  typeof notegicBlockNoteSchema.styleSchema
>;
