CREATE OR REPLACE FUNCTION trigger_function_accounting_deleted_routine_task()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_routine_task: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH station_deltas AS (
        SELECT
            station_id,
            count(*) as total_delta
        FROM old_table ot
        GROUP BY station_id
    ),
    updated_stations AS (
        UPDATE "StationTable" s
        SET
            routine_task_count = GREATEST(0, routine_task_count - sd.total_delta),
            updated_at = NOW()
        FROM station_deltas sd
        WHERE s.id = sd.station_id
        RETURNING s.id, s.owner_id
    ),
    owner_deltas AS (
        SELECT
            us.owner_id,
            sum(sd.total_delta) as total_delta
        FROM updated_stations us
        JOIN station_deltas sd ON sd.station_id = us.id
        GROUP BY us.owner_id
    )
    UPDATE "UserAccountTable" ua
    SET
        routine_task_count = GREATEST(0, routine_task_count - od.total_delta),
        updated_at = NOW()
    FROM owner_deltas od
    WHERE ua.user_id = od.owner_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_deleted_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_deleted_routine_task
    AFTER DELETE
    ON "RoutineTaskTable"
    REFERENCING OLD TABLE AS old_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_deleted_routine_task();
