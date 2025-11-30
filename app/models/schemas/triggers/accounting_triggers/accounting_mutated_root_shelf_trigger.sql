CREATE OR REPLACE FUNCTION trigger_function_accounting_mutated_root_shelf()
RETURNS TRIGGER AS $$
DECLARE
    current_count INTEGER;
    max_count INTEGER;
    plan_name TEXT;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        SELECT
            pl.max_root_shelf_count,
            u.plan::TEXT
        INTO
            max_count,
            plan_name
        FROM "UserTable" u
        JOIN "PlanLimitationTable" pl ON u.plan = pl.key
        WHERE u.id = NEW.owner_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Data integrity: Cannot find owner for RootShelf. Possible orphan record.'
            USING ERRCODE = 'data_exception';
        END IF;

        UPDATE "UserAccountTable" ua
        SET 
            root_shelf_count = root_shelf_count + 1,
            updated_at = NOW()
        WHERE ua.user_id = NEW.owner_id
        RETURNING root_shelf_count INTO current_count;

        IF current_count > max_count THEN
            RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % root shelves. Current count: %.', 
                plan_name, max_count, current_count
            USING ERRCODE = 'check_violation';
        END IF;

        RETURN NEW;

    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE "UserAccountTable"
        SET
            root_shelf_count = GREATEST(0, root_shelf_count - 1),
            updated_at = NOW()
        WHERE user_id = OLD.owner_id;

        IF NOT FOUND THEN
             RAISE EXCEPTION 'Data integrity: Cannot find UserAccount of the owner for RootShelf (Owner ID: %).', OLD.owner_id
             USING ERRCODE = 'integrity_constraint_violation';
        END IF;

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Seperator ==============================

CREATE TRIGGER trigger_accounting_mutated_root_shelf
    BEFORE INSERT OR DELETE ON "RootShelfTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_mutated_root_shelf();