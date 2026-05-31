CREATE OR REPLACE FUNCTION trigger_function_accounting_deleted_routine_tag()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'DELETE') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_deleted_routine_tag: %. Expected DELETE.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    -- Count only owned tags: permission = 'Owner'. Shared tags do not consume owner quota.
    WITH owner_tag_deltas AS (
        SELECT
            user_id,
            count(*) as total_delta
        FROM old_table
        WHERE permission = 'Owner'
        GROUP BY user_id
    )
    UPDATE "UserAccountTable" ua
    SET
        routine_tag_count = GREATEST(0, routine_tag_count - otd.total_delta),
        updated_at = NOW()
    FROM owner_tag_deltas otd
    WHERE ua.user_id = otd.user_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_deleted_routine_tag ON "UsersToRoutineTagsTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_deleted_routine_tag
    AFTER DELETE
    ON "UsersToRoutineTagsTable"
    REFERENCING OLD TABLE AS old_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_deleted_routine_tag();
