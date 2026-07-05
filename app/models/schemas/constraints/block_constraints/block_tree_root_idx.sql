DROP INDEX IF EXISTS "block_idx_tree_root";
DROP INDEX IF EXISTS "block_idx_single_root_head";
DROP INDEX IF EXISTS "block_idx_single_child_head";
DROP INDEX IF EXISTS "block_idx_unique_prev_block";
DROP INDEX IF EXISTS "block_idx_unique_next_block";

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX CONCURRENTLY block_idx_single_root_head
ON "BlockTable" (block_pack_id)
WHERE parent_block_id IS NULL AND prev_block_id IS NULL AND deleted_at IS NULL;

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX CONCURRENTLY block_idx_single_child_head
ON "BlockTable" (block_pack_id, parent_block_id)
WHERE parent_block_id IS NOT NULL AND prev_block_id IS NULL AND deleted_at IS NULL;

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX CONCURRENTLY block_idx_unique_prev_block
ON "BlockTable" (prev_block_id)
WHERE prev_block_id IS NOT NULL AND deleted_at IS NULL;

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX CONCURRENTLY block_idx_unique_next_block
ON "BlockTable" (next_block_id)
WHERE next_block_id IS NOT NULL AND deleted_at IS NULL;
