CREATE OR REPLACE FUNCTION trigger_function_cascading_restore_soft_deleted_root_shelf()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
        WITH updated_sub_shelves AS (
            UPDATE "SubShelfTable"
            SET deleted_at = NULL
            WHERE root_shelf_id = NEW.id 
            AND deleted_at = OLD.deleted_at
            RETURNING id
        )
        UPDATE "MaterialTable"
        SET deleted_at = NULL
        FROM updated_sub_shelves
        WHERE "MaterialTable".parent_sub_shelf_id = updated_sub_shelves.id
        AND "MaterialTable".deleted_at = OLD.deleted_at;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

DROP TRIGGER IF EXISTS trigger_cascading_restore_soft_deleted_root_shelf ON "RootShelfTable"

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_cascading_restore_soft_deleted_root_shelf
    AFTER UPDATE 
    ON "RootShelfTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_cascading_restore_soft_deleted_root_shelf();