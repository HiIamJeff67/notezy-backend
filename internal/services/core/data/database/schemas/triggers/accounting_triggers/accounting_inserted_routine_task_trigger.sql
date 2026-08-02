CREATE OR REPLACE FUNCTION trigger_function_accounting_inserted_routine_task()
RETURNS TRIGGER AS $$
DECLARE
    station_owner_id uuid;
    station_owner_plan "UserPlan";
    max_routine_task_cost_unit_count integer;
    updated_routine_task_cost_unit_count bigint;
BEGIN
    IF (TG_OP <> 'INSERT') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_inserted_routine_task: %. Expected INSERT.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    NEW.cost_unit = (octet_length(COALESCE(NEW.payload::text, ''))::bigint + 1023) / 1024;

    SELECT
        s.owner_id,
        u.plan,
        pl.max_routine_task_cost_unit_count
    INTO
        station_owner_id,
        station_owner_plan,
        max_routine_task_cost_unit_count
    FROM "StationTable" s
    JOIN "UserTable" u ON u.id = s.owner_id
    JOIN "PlanLimitationTable" pl ON pl.key = u.plan
    JOIN "RoutineTable" r ON r.station_id = s.id
    WHERE r.id = NEW.routine_id;

    UPDATE "UserAccountTable"
    SET
        routine_task_cost_unit_count = routine_task_cost_unit_count + NEW.cost_unit,
        updated_at = NOW()
    WHERE user_id = station_owner_id
    RETURNING routine_task_cost_unit_count INTO updated_routine_task_cost_unit_count;

    IF updated_routine_task_cost_unit_count > max_routine_task_cost_unit_count THEN
        RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % routine task cost units. Current cost unit count: %.',
            station_owner_plan, max_routine_task_cost_unit_count, updated_routine_task_cost_unit_count
        USING ERRCODE = 'check_violation';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_inserted_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_inserted_routine_task
    BEFORE INSERT
    ON "RoutineTaskTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_inserted_routine_task();
