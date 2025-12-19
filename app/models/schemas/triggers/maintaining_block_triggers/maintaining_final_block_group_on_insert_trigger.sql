CREATE OR REPLACE FUNCTION trigger_function_maintaining_final_block_group_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- find the new tail in the block groups list of the block pack
    WITH new_final_candidates AS (
        SELECT 
            i.block_pack_id,
            i.id AS new_final_id
        FROM inserted_rows i
        WHERE NOT EXISTS (
            -- if there's no data point to it, then it is the tail of the block groups list
            SELECT 1 
            FROM inserted_rows i2 
            WHERE i2.block_pack_id = i.block_pack_id 
              AND i2.prev_block_group_id = i.id
        )
    )

    UPDATE "BlockPackTable" bp
    SET final_block_group_id = nfc.new_final_id
    FROM new_final_candidates nfc
    WHERE bp.id = nfc.block_pack_id
      -- only update while the final block group id is conflicted with the prev block group id of the new block group
      AND (
          bp.final_block_group_id IS NULL -- condition 1: if the original final block group id is null
          OR 
          EXISTS ( -- condition 2: if the original final block group id is not null
              SELECT 1 
              FROM inserted_rows i 
              WHERE i.block_pack_id = bp.id 
                AND i.prev_block_group_id = bp.final_block_group_id
          )
      );
      
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

DROP TRIGGER IF EXISTS trigger_maintaining_final_block_group_on_insert ON "BlockGroupTable";

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_maintaining_final_block_group_on_insert
    AFTER INSERT ON "BlockGroupTable"
    REFERENCING NEW TABLE AS inserted_rows
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_maintaining_final_block_group_on_insert();