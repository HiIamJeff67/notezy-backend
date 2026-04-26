CREATE OR REPLACE FUNCTION maintaining_descendants_block_group_id_on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Process all changed blocks
    WITH changed_blocks AS (
        SELECT DISTINCT n.id, n.block_group_id as new_block_group_id, o.block_group_id as old_block_group_id
        FROM NEW n JOIN OLD o ON n.id = o.id
        WHERE n.block_group_id IS DISTINCT FROM o.block_group_id
    ),
    all_descendants AS (
        SELECT d.id, cb.new_block_group_id
        FROM changed_blocks cb
        JOIN (
            WITH RECURSIVE desc AS (
                SELECT id, new_block_group_id FROM changed_blocks
                UNION ALL
                SELECT b.id, d.new_block_group_id FROM "BlockTable" b JOIN desc d ON b.parent_block_id = d.id
                WHERE b.deleted_at IS NULL
            )
            SELECT * FROM desc
        ) d ON d.id != cb.id
    )
    UPDATE "BlockTable" 
    SET block_group_id = all_descendants.new_block_group_id
    FROM all_descendants
    WHERE "BlockTable".id = all_descendants.id;
    
    RETURN NULL;  -- Return value ignored for AFTER triggers
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_maintaining_descendants_block_group_id_on_update
    AFTER UPDATE 
    ON "BlockTable"
    FOR EACH STATEMENT 
    EXECUTE FUNCTION maintaining_descendants_block_group_id_on_update();