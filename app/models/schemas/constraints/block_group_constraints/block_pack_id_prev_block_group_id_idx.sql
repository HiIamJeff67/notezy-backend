DROP INDEX IF EXISTS "block_group_idx_name_block_pack_id_prev_block_group_id"; 

-- ============================== SQL Seperator ==============================

CREATE UNIQUE INDEX block_group_idx_name_block_pack_id_prev_block_group_id
ON "BlockGroupTable" (block_pack_id, prev_block_group_id)
WHERE deleted_at IS NULL;