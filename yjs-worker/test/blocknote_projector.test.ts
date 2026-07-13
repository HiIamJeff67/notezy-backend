import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import type { Block } from "@blocknote/core";
import { blocksToYXmlFragment } from "@blocknote/core/yjs";
import * as Y from "yjs";

import { BlockNoteProjector } from "../src/blocknote_projector.js";
import { notezyBlockNoteEditor } from "../src/blocknote_schema.js";

const blockNoteProjector = new BlockNoteProjector();

async function readFixture(name: string): Promise<Block[]> {
  const fixture = await readFile(
    new URL(`../../tmp/${name}`, import.meta.url),
    "utf8"
  );

  return JSON.parse(fixture) as Block[];
}

for (const fixtureName of [
  "temp_deep_block_contents.json",
  "temp_wide_block_contents.json",
]) {
  test(`projects ${fixtureName} through the canonical BlockNote Y.XmlFragment`, async () => {
    const sourceBlocks = await readFixture(fixtureName);
    const document = new Y.Doc();
    blocksToYXmlFragment(
      notezyBlockNoteEditor,
      sourceBlocks,
      document.getXmlFragment("document-store")
    );

    const projectedBlocks = blockNoteProjector.projectYjsDocument(document);
    const rematerializedDocument = new Y.Doc();
    blocksToYXmlFragment(
      notezyBlockNoteEditor,
      projectedBlocks,
      rematerializedDocument.getXmlFragment("document-store")
    );

    assert.deepEqual(
      blockNoteProjector.projectYjsDocument(rematerializedDocument),
      projectedBlocks
    );
    assert.equal(projectedBlocks.length, sourceBlocks.length);
  });
}
