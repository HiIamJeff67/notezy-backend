UPDATE "BlockGroupTable"
SET deleted_at = NOW()
WHERE id IN ?
AND deleted_at IS NULL
AND NOT EXISTS (
    SELECT 1 
    FROM "BlockTable" 
    WHERE block_group_id = "BlockGroupTable".id 
    AND deleted_at IS NULL
)