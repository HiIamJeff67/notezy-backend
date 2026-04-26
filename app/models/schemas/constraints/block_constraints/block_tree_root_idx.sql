DROP INDEX IF EXISTS "block_idx_tree_root";

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX CONCURRENTLY block_idx_tree_root
ON "BlockTable" (block_group_id)
WHERE parent_block_id IS NULL AND deleted_at IS NULL;