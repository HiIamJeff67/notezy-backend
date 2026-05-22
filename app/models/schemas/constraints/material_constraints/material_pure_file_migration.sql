ALTER TABLE "MaterialTable"
    ADD COLUMN IF NOT EXISTS content_type "MaterialContentType";

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    ADD COLUMN IF NOT EXISTS parse_media_type VARCHAR(128);

-- ============================== SQL Separator ==============================

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'MaterialTable'
          AND column_name = 'type'
    ) THEN
        UPDATE "MaterialTable"
        SET content_type = CASE type::text
            WHEN 'Textbook' THEN 'text/plain'::"MaterialContentType"
            WHEN 'Notebook' THEN 'text/plain'::"MaterialContentType"
            WHEN 'LearningCards' THEN 'text/html'::"MaterialContentType"
            WHEN 'Workflow' THEN 'application/json'::"MaterialContentType"
            ELSE 'text/plain'::"MaterialContentType"
        END
        WHERE content_type IS NULL;
    ELSE
        UPDATE "MaterialTable"
        SET content_type = 'text/plain'::"MaterialContentType"
        WHERE content_type IS NULL;
    END IF;
END
$$;

-- ============================== SQL Separator ==============================

UPDATE "MaterialTable"
SET parse_media_type = ''
WHERE parse_media_type IS NULL;

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    ALTER COLUMN content_type SET DEFAULT 'text/plain'::"MaterialContentType";

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    ALTER COLUMN content_type SET NOT NULL;

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    ALTER COLUMN parse_media_type SET DEFAULT '';

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    ALTER COLUMN parse_media_type SET NOT NULL;

-- ============================== SQL Separator ==============================

ALTER TABLE "MaterialTable"
    DROP COLUMN IF EXISTS type;

-- ============================== SQL Separator ==============================

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_type
        WHERE typname = 'MaterialType'
    ) THEN
        DROP TYPE "MaterialType";
    END IF;
END
$$;
