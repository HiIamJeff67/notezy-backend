CREATE OR REPLACE FUNCTION trigger_function_accounting_updated_routine_task()
RETURNS TRIGGER AS $$
DECLARE
    old_station_owner_id uuid;
    new_station_owner_id uuid;
    new_station_owner_plan "UserPlan";
    max_routine_task_cost_unit_count integer;
    new_cost_unit bigint;
    cost_unit_delta bigint;
    updated_routine_task_cost_unit_count bigint;
BEGIN
    IF (TG_OP <> 'UPDATE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_updated_routine_task: %. Expected UPDATE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    new_cost_unit = (octet_length(COALESCE(NEW.payload::text, ''))::bigint + 1023) / 1024;
    cost_unit_delta = new_cost_unit - OLD.cost_unit;
    NEW.cost_unit = new_cost_unit;

    SELECT
        s.owner_id,
        u.plan,
        pl.max_routine_task_cost_unit_count
    INTO
        new_station_owner_id,
        new_station_owner_plan,
        max_routine_task_cost_unit_count
    FROM "StationTable" s
    JOIN "UserTable" u ON u.id = s.owner_id
    JOIN "PlanLimitationTable" pl ON pl.key = u.plan
    JOIN "RoutineTable" r ON r.station_id = s.id
    WHERE r.id = NEW.routine_id;

    SELECT s.owner_id
    INTO old_station_owner_id
    FROM "StationTable" s
    JOIN "RoutineTable" r ON r.station_id = s.id
    WHERE r.id = OLD.routine_id;

    IF old_station_owner_id = new_station_owner_id AND cost_unit_delta = 0 THEN
        RETURN NEW;
    END IF;

    IF old_station_owner_id <> new_station_owner_id THEN
        UPDATE "UserAccountTable"
        SET
            routine_task_cost_unit_count = GREATEST(0, routine_task_cost_unit_count - OLD.cost_unit),
            updated_at = NOW()
        WHERE user_id = old_station_owner_id;

        UPDATE "UserAccountTable"
        SET
            routine_task_cost_unit_count = routine_task_cost_unit_count + new_cost_unit,
            updated_at = NOW()
        WHERE user_id = new_station_owner_id
        RETURNING routine_task_cost_unit_count INTO updated_routine_task_cost_unit_count;
    ELSE
        UPDATE "UserAccountTable"
        SET
            routine_task_cost_unit_count = GREATEST(0, routine_task_cost_unit_count + cost_unit_delta),
            updated_at = NOW()
        WHERE user_id = new_station_owner_id
        RETURNING routine_task_cost_unit_count INTO updated_routine_task_cost_unit_count;
    END IF;

    IF updated_routine_task_cost_unit_count > max_routine_task_cost_unit_count THEN
        RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % routine task cost units. Current cost unit count: %.',
            new_station_owner_plan, max_routine_task_cost_unit_count, updated_routine_task_cost_unit_count
        USING ERRCODE = 'check_violation';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_updated_routine_task ON "RoutineTaskTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_updated_routine_task
    BEFORE UPDATE OF routine_id, payload, cost_unit
    ON "RoutineTaskTable"
    FOR EACH ROW
    EXECUTE FUNCTION trigger_function_accounting_updated_routine_task();
