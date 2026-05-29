CREATE OR REPLACE FUNCTION trigger_function_accounting_deleted_routine_tag()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_routine_tag: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH owner_deltas AS (
        SELECT
            s.owner_id,
            count(*) as total_delta
        FROM old_table ot
        JOIN "StationTable" s ON s.id = ot.station_id
        GROUP BY s.owner_id
    )
    UPDATE "UserAccountTable" ua
    SET
        routine_tag_count = GREATEST(0, routine_tag_count - od.total_delta),
        updated_at = NOW()
    FROM owner_deltas od
    WHERE ua.user_id = od.owner_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_deleted_routine_tag ON "RoutineTagTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_deleted_routine_tag
    AFTER DELETE
    ON "RoutineTagTable"
    REFERENCING OLD TABLE AS old_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_deleted_routine_tag();
