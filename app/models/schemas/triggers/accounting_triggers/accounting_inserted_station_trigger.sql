CREATE OR REPLACE FUNCTION trigger_function_accounting_inserted_station()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP <> 'INSERT') THEN
        RAISE EXCEPTION 'Invalid operation for trigger_function_accounting_inserted_station: %. Expected INSERT.', TG_OP
        USING ERRCODE = 'program_limit_exceeded';
    END IF;

    WITH owner_deltas AS (
        SELECT
            owner_id,
            count(*) as total_delta
        FROM new_table
        GROUP BY owner_id
    )
    UPDATE "UserAccountTable" ua
    SET
        station_count = station_count + od.total_delta,
        updated_at = NOW()
    FROM owner_deltas od
    WHERE ua.user_id = od.owner_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================== SQL Separator ==============================

DROP TRIGGER IF EXISTS trigger_accounting_inserted_station ON "StationTable"

-- ============================== SQL Separator ==============================

CREATE TRIGGER trigger_accounting_inserted_station
    AFTER INSERT
    ON "StationTable"
    REFERENCING NEW TABLE AS new_table
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_function_accounting_inserted_station();
