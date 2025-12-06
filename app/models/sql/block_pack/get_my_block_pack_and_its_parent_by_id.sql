-- name: GetMyBlockPackAndItsParentById
SELECT
    bp.id,
    bp.name,
    bp.icon,
    bp.header_background_url,
    bp.block_count,
    bp.deleted_at,
    bp.updated_at,
    bp.created_at,
    ss.root_shelf_id AS root_shelf_id
    ss.id AS parent_sub_shelf_id, 
    ss.name AS parent_sub_shelf_name,
    ss.prev_sub_shelf_id AS parent_sub_shelf_prev_sub_shelf_id,
    ss.path AS parent_sub_shelf_path,
    ss.deleted_at AS parent_sub_shelf_deleted_at,
    ss.updated_at AS parent_sub_shelf_updated_at,
    ss.created_at AS parent_sub_shelf_created_at
FROM "BlockPackTable" bp
JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id
JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id
WHERE bp.id = $1 AND uts.user_id = $2 AND uts.permission = ANY($3::"AccessControlPermission"[])
    AND (
        CASE
            WHEN $4 = 0 THEN m.deleted_at IS NOT NULL
            WHEN $4 = 2 THEN m.deleted_at IS NULL
            ELSE true
        END
    )