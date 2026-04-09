CREATE OR REPLACE FUNCTION trigger_function_maintaining_final_block_group_on_delete()
RETURNS TRIGGER AS $$
DECLARE
    new_final_block_group_id UUID;
    target_block_pack_id UUID;
BEGIN
    -- Determine the block_pack_id.
    -- If it's a hard delete, OLD holds the data.
    -- If it's a soft delete (UPDATE set deleted_at), NEW holds the data.
    IF (TG_OP = 'DELETE') THEN
        target_block_pack_id := OLD.block_pack_id;
    ELSE
        target_block_pack_id := NEW.block_pack_id;
    END IF;

    -- Find the true tail of the list using the Set Difference logic:
    -- The Final Block is the one that exists (not deleted) AND is not anyone's "prev".
    -- Logic: {All Active IDs in Pack} - {All Active Prev IDs in Pack} = {Final ID}
    SELECT id INTO new_final_block_group_id
    FROM "BlockGroupTable"
    WHERE block_pack_id = target_block_pack_id
      AND deleted_at IS NULL
      AND id NOT IN (
          SELECT prev_block_group_id
          FROM "BlockGroupTable"
          WHERE block_pack_id = target_block_pack_id
            AND prev_block_group_id IS NOT NULL
            AND deleted_at IS NULL
      )
    LIMIT 1;

    -- Update the BlockPack to point to the calculated tail
    -- If the pack is empty (new_final_block_group_id is NULL), we set it to NULL.
    UPDATE "BlockPackTable"
    SET final_block_group_id = new_final_block_group_id
    WHERE id = target_block_pack_id;

    RETURN NULL; -- Return value ignored for AFTER triggers
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_maintaining_final_block_group_on_delete ON "BlockGroupTable";

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_maintaining_final_block_group_on_delete
    AFTER DELETE OR UPDATE OF deleted_at
    ON "BlockGroupTable"
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_maintaining_final_block_group_on_delete();