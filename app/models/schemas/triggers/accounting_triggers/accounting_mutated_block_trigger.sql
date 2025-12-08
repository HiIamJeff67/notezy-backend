CREATE TRIGGER trigger_function_accounting_deleted_block();
RETURNS TRIGGER AS $$
DECLARE
    r RECORD;
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_block: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH old_blocks_agg AS (
        SELECT 
            block_group_id, 
            count(*) as count_delta
        FROM old_table
        GROUP BY block_group_id
    ),
    owner_deltas AS (
        SELECT 
            bg.owner_id,
            sum(oba.count_delta) as total_delta
        FROM old_blocks_agg oba
        JOIN "BlockGroupTable" bg ON oba.block_group_id = bg.id
        GROUP BY bg.owner_id
    )
    UPDATE "UserAccountTable" ua
    SET 
        block_count = GREATEST(0, block_count - od.total_delta),
        updated_at = NOW()
    FROM owner_deltas od
    WHERE ua.user_id = od.owner_id;

    WITH old_blocks_agg AS (
        SELECT 
            block_group_id, 
            count(*) as count_delta
        FROM old_table
        GROUP BY block_group_id
    ),
    bp_updates AS (
        SELECT 
            bg.block_pack_id, 
            sum(oba.count_delta) as total_delta
        FROM old_blocks_agg oba
        JOIN "BlockGroupTable" bg ON oba.block_group_id = bg.id
        GROUP BY bg.block_pack_id
    )
    UPDATE "BlockPackTable" bp
    SET 
        block_count = GREATEST(0, block_count - bpu.total_delta),
        updated_at = NOW()
    FROM bp_updates bpu
    WHERE bp.id = bpu.block_pack_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_accounting_deleted_block
    AFTER DELETE ON "BlockTable"
    REFERENCING OLD TABLE AS old_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_deleted_block();