CREATE OR REPLACE FUNCTION trigger_function_accounting_inserted_routine_task()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'INSERT') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_inserted_routine_task: %. Expected INSERT.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH owner_deltas AS (
        SELECT
            s.owner_id,
            count(*) as total_delta
        FROM new_table nt
        JOIN "StationTable" s ON s.id = nt.station_id
        GROUP BY s.owner_id
    )
    UPDATE "UserAccountTable" ua
    SET
        routine_task_count = routine_task_count + od.total_delta,
        updated_at = NOW()
    FROM owner_deltas od
    WHERE ua.user_id = od.owner_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_inserted_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_inserted_routine_task
    AFTER INSERT
    ON "RoutineTaskTable"
    REFERENCING NEW TABLE AS new_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_inserted_routine_task();
