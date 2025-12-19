CREATE OR REPLACE FUNCTION trigger_function_maintaining_final_block_group_on_delete()
RETURNS TRIGGER AS $$
DECLARE
    current_final_id UUID;
BEGIN
    SELECT final_block_group_id INTO current_final_id
    FROM "BlockPackTable"
    WHERE id = OLD.block_pack_id;

    IF (OLD.id = current_final_id) THEN
        UPDATE "BlockPackTable"
        SET final_block_group_id = OLD.prev_block_group_id
        WHERE id = OLD.block_pack_id;
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

DROP TRIGGER IF EXISTS trigger_maintaining_final_block_group_on_delete ON "BlockGroupTable";

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_maintaining_final_block_group_on_delete
    AFTER DELETE
    ON "BlockGroupTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_maintaining_final_block_group_on_delete();