CREATE OR REPLACE FUNCTION trigger_function_accounting_mutated_block()
RETURNS TRIGGER AS $$
DECLARE
    current_count INTEGER;
    max_count INTEGER;
    current_count_per_block_pack INTEGER;
    max_count_per_block_pack INTEGER;
    plan_name TEXT;
    owner_id UUID;
    block_pack_id_val UUID;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        -- 1. Get Owner ID, Plan limits, and BlockPack ID (Read-only info)
        SELECT
            pl.max_block_count,
            pl.max_block_count_per_block_pack,
            u.plan::TEXT,
            u.id,
            bg.block_pack_id
        INTO
            max_count,
            max_count_per_block_pack,
            plan_name,
            owner_id,
            block_pack_id_val
        FROM "BlockGroupTable" bg
        JOIN "BlockPackTable" bp ON bg.block_pack_id = bp.id
        JOIN "UserTable" u ON bg.owner_id = u.id
        JOIN "PlanLimitationTable" pl ON u.plan = pl.key
        WHERE bg.id = NEW.block_group_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Data integrity: Cannot find owner for Block (BlockGroup ID: %). Possible orphan record.', NEW.block_group_id
            USING ERRCODE = 'data_exception';
        END IF;

        UPDATE "UserAccountTable"
        SET
            block_count = block_count + 1,
            updated_at = NOW()
        WHERE user_id = owner_id
        RETURNING block_count INTO current_count;

        IF current_count > max_count THEN
            RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % blocks. Current count: %.', 
                plan_name, max_count, current_count
            USING ERRCODE = 'check_violation';
        END IF;

        UPDATE "BlockPackTable"
        SET
            block_count = block_count + 1,
            updated_at = NOW()
        WHERE id = block_pack_id_val
        RETURNING block_count INTO current_count_per_block_pack;

        IF current_count_per_block_pack > max_count_per_block_pack THEN
            RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % blocks in each block pack. Current count: %.', 
                plan_name, max_count_per_block_pack, current_count_per_block_pack
            USING ERRCODE = 'check_violation';
        END IF;

        RETURN NEW;

    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE "UserAccountTable"
        SET
            block_count = GREATEST(0, block_count - 1),
            updated_at = NOW()
        FROM "BlockGroupTable" bg
        WHERE bg.id = OLD.block_group_id
        AND user_id = bg.owner_id;

        IF NOT FOUND THEN
             RAISE EXCEPTION 'Data integrity: Cannot find UserAccount of the owner for Block (BlockGroup ID: %). Possible orphan record.', OLD.block_group_id
             USING ERRCODE = 'integrity_constraint_violation';
        END IF;

        UPDATE "BlockPackTable"
        SET
            block_count = GREATEST(0, block_count - 1),
            updated_at = NOW()
        FROM "BlockGroupTable" bg
        WHERE bg.id = OLD.block_group_id
        AND "BlockPackTable".id = bg.block_pack_id;

        IF NOT FOUND THEN
             RAISE EXCEPTION 'Data integrity: Cannot find BlockPack for Block (BlockGroup ID: %). Possible orphan record.', OLD.block_group_id
             USING ERRCODE = 'integrity_constraint_violation';
        END IF;

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_accounting_mutated_block
    BEFORE INSERT OR DELETE ON "BlockTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_mutated_block();