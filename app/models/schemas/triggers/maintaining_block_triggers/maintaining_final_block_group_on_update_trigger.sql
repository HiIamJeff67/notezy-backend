CREATE OR REPLACE FUNCTION trigger_function_maintaining_final_block_group_on_update()
RETURNS TRIGGER AS $$
BEGIN
    WITH new_final_candidates AS (
        SELECT 
            u.block_pack_id,
            u.id AS new_final_id
        FROM updated_rows u
        WHERE NOT EXISTS (
            SELECT 1 
            FROM updated_rows u2 
            WHERE u2.block_pack_id = u.block_pack_id 
              AND u2.prev_block_group_id = u.id
        )
    )

    UPDATE "BlockPackTable" bp
    SET final_block_group_id = nfc.new_final_id
    FROM new_final_candidates nfc
    WHERE bp.id = nfc.block_pack_id
      AND (
          bp.final_block_group_id IS NULL 
          OR 
          EXISTS ( 
              SELECT 1 
              FROM updated_rows u 
              WHERE u.block_pack_id = bp.id 
                AND u.prev_block_group_id = bp.final_block_group_id
          )
      );
      
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

DROP TRIGGER IF EXISTS trigger_maintaining_final_block_group_on_update ON "BlockGroupTable";

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_maintaining_final_block_group_on_update
    AFTER UPDATE OF prev_block_group_id ON "BlockGroupTable" -- 只監聽 prev_block_group_id 的變更
    REFERENCING NEW TABLE AS updated_rows
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_maintaining_final_block_group_on_update();