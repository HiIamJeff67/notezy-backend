CREATE OR REPLACE FUNCTION trigger_function_accounting_inserted_routine_task()
RETURNS TRIGGER AS $$
DECLARE
    r RECORD;
BEGIN
    IF (TG_OP <> 'INSERT') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_inserted_routine_task: %. Expected INSERT.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    FOR r IN
        WITH station_deltas AS (
            SELECT
                station_id,
                count(*) as total_delta
            FROM new_table
            GROUP BY station_id
        ),
        owner_deltas AS (
            SELECT
                s.owner_id,
                sum(sd.total_delta) as total_delta
            FROM station_deltas sd
            JOIN "StationTable" s ON s.id = sd.station_id
            GROUP BY s.owner_id
        ),
        updated_accounts AS (
            UPDATE "UserAccountTable" ua
            SET
                routine_task_count = routine_task_count + od.total_delta,
                updated_at = NOW()
            FROM owner_deltas od
            WHERE ua.user_id = od.owner_id
            RETURNING ua.user_id
        ),
        station_counts AS (
            SELECT
                s.id,
                s.owner_id,
                count(rt.id) as routine_task_count
            FROM "StationTable" s
            JOIN station_deltas sd ON sd.station_id = s.id
            LEFT JOIN "RoutineTaskTable" rt ON rt.station_id = s.id
            GROUP BY s.id, s.owner_id
        )

        SELECT
            sc.id, u.plan, sc.routine_task_count, pl.max_routine_task_count_per_station
        FROM station_counts sc
        JOIN "UserTable" u ON sc.owner_id = u.id
        JOIN "PlanLimitationTable" pl ON u.plan = pl.key
        LEFT JOIN updated_accounts ua ON ua.user_id = sc.owner_id
        WHERE sc.routine_task_count > pl.max_routine_task_count_per_station
    LOOP
        RAISE EXCEPTION 'Quota exceeded: Plan "%" allows maximum % routine tasks per station. Current count: %.',
            r.plan, r.max_routine_task_count_per_station, r.routine_task_count
        USING ERRCODE = 'check_violation';
    END LOOP;

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
