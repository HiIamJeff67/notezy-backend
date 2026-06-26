CREATE OR REPLACE FUNCTION trigger_function_accounting_updated_routine_task()
RETURNS TRIGGER AS $$
DECLARE
    station_owner_id uuid;
    station_owner_plan "UserPlan";
    max_routine_task_cost_unit_count integer;
    new_cost_unit bigint;
    cost_unit_delta bigint;
    updated_routine_task_cost_unit_count bigint;
BEGIN
    IF (TG_OP <> 'UPDATE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_updated_routine_task: %. Expected UPDATE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    IF NEW.station_id <> OLD.station_id THEN
        RAISE EXCEPTION 'RoutineTask station move is not supported by accounting triggers.'
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    new_cost_unit = (octet_length(COALESCE(NEW.payload::text, ''))::bigint + 1023) / 1024;
    cost_unit_delta = new_cost_unit - OLD.cost_unit;
    NEW.cost_unit = new_cost_unit;

    IF cost_unit_delta = 0 THEN
        RETURN NEW;
    END IF;

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
    WHERE s.id = NEW.station_id;

    UPDATE "UserAccountTable"
    SET
        routine_task_cost_unit_count = GREATEST(0, routine_task_cost_unit_count + cost_unit_delta),
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

DROP TRIGGER IF EXISTS trigger_accounting_updated_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_updated_routine_task
    BEFORE UPDATE OF payload, cost_unit
    ON "RoutineTaskTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_updated_routine_task();
