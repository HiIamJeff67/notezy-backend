DROP INDEX IF EXISTS "block_group_idx_name_block_pack_id_prev_block_group_id"; 

-- ============================== SQL Seperator ==============================

ALTER TABLE "BlockGroupTable"
ADD CONSTRAINT block_group_idx_name_block_pack_id_prev_block_group_id
UNIQUE (block_pack_id, prev_block_group_id)
DEFERRABLE INITIALLY DEFERRED;