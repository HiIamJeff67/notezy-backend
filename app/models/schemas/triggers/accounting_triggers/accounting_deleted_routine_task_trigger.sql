CREATE OR REPLACE FUNCTION trigger_function_accounting_deleted_routine_task()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_routine_task: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH owner_deltas AS (
        SELECT
            owner_id,
            count(*) as total_delta
        FROM old_table
        GROUP BY owner_id
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
