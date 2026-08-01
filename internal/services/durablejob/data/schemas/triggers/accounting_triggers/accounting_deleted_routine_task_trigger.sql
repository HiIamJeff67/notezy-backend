CREATE OR REPLACE FUNCTION trigger_function_accounting_deleted_routine_task()
RETURNS TRIGGER AS $$
DECLARE
    station_owner_id uuid;
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_routine_task: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    SELECT owner_id
    INTO station_owner_id
    FROM "StationTable" s
    JOIN "RoutineTable" r ON r.station_id = s.id
    WHERE r.id = OLD.routine_id;

    UPDATE "UserAccountTable" ua
    SET
        routine_task_cost_unit_count = GREATEST(0, routine_task_cost_unit_count - OLD.cost_unit),
        updated_at = NOW()
    WHERE ua.user_id = station_owner_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_deleted_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_deleted_routine_task
    AFTER DELETE
    ON "RoutineTaskTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_deleted_routine_task();
