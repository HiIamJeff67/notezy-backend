CREATE OR REPLACE FUNCTION trigger_function_cascading_soft_delete_root_shelf()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        WITH updated_sub_shelves AS (
            UPDATE "SubShelfTable"
            SET deleted_at = NEW.deleted_at
            WHERE root_shelf_id = NEW.id 
            AND deleted_at IS NULL
            RETURNING id
        )
        UPDATE "MaterialTable"
        SET deleted_at = NEW.deleted_at
        FROM updated_sub_shelves
        WHERE "MaterialTable".parent_sub_shelf_id = updated_sub_shelves.id
        AND "MaterialTable".deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_cascading_soft_delete_root_shelf
    AFTER UPDATE ON "RootShelfTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_cascading_soft_delete_root_shelf();
