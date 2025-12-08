CREATE OR REPLACE FUNCTION trigger_function_accounting_inserted_block()
RETURNS TRIGGER AS $$
DECLARE
    r RECORD;
BEGIN
    IF (TG_OP <> 'INSERT') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_inserted_block: %. Expected INSERT.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    FOR r IN -- batch account the user block count, if all the mutated blocks belong to one user, then this will only execute once
        WITH new_blocks_agg AS (
            SELECT 
                block_group_id, 
                count(*) as count_delta
            FROM new_table
            GROUP BY block_group_id
        ),
        owner_deltas AS (
            SELECT 
                bg.owner_id,
                sum(nba.count_delta) as total_delta
            FROM new_blocks_agg nba
            JOIN "BlockGroupTable" bg ON nba.block_group_id = bg.id
            GROUP BY bg.owner_id
        ),
        updated_accounts AS (
            UPDATE "UserAccountTable" ua
            SET 
                block_count = block_count + od.total_delta,
                updated_at = NOW()
            FROM owner_deltas od
            WHERE ua.user_id = od.owner_id
            RETURNING ua.user_id, ua.block_count
        )

        -- finally select the limitation and iteratively check the updated block count of the user is not exceeded the limitation
        SELECT 
            u.id, u.plan, ua.block_count, pl.max_block_count
        FROM updated_accounts ua
        JOIN "UserTable" u ON ua.user_id = u.id
        JOIN "PlanLimitationTable" pl ON u.plan = pl.key
        WHERE ua.block_count > pl.max_block_count
    LOOP
        RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % blocks. Current count: %.', 
            r.plan, r.max_block_count, r.block_count
        USING ERRCODE = 'check_violation';
    END LOOP;

    FOR r IN
        WITH new_blocks_agg AS (
            SELECT 
                block_group_id, 
                count(*) as count_delta
            FROM new_table
            GROUP BY block_group_id
        ),
        bp_deltas AS (
            SELECT 
                bg.block_pack_id, 
                sum(nba.count_delta) as total_delta
            FROM new_blocks_agg nba
            JOIN "BlockGroupTable" bg ON nba.block_group_id = bg.id
            GROUP BY bg.block_pack_id
        ),
        updated_packs AS (
            UPDATE "BlockPackTable" bp
            SET 
                block_count = block_count + bpu.total_delta,
                updated_at = NOW()
            FROM bp_deltas bpu
            WHERE bp.id = bpu.block_pack_id
            RETURNING bp.id, bp.block_count
        )

        SELECT DISTINCT ON (up.id)
            up.id, up.block_count, pl.max_block_count_per_block_pack, u.plan
        FROM updated_packs up
        JOIN "BlockGroupTable" bg ON bg.block_pack_id = up.id
        JOIN "UserTable" u ON bg.owner_id = u.id
        JOIN "PlanLimitationTable" pl ON u.plan = pl.key
        WHERE up.block_count > pl.max_block_count_per_block_pack
    LOOP
            RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % blocks in each block pack. Current count: %.', 
            r.plan, r.max_block_count_per_block_pack, r.block_count
        USING ERRCODE = 'check_violation';
    END LOOP;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_accounting_inserted_block
    AFTER INSERT ON "BlockTable"
    REFERENCING NEW TABLE AS new_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_inserted_block();