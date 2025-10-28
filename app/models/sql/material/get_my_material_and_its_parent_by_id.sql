-- name: GetMyMaterialAndItsParentById
SELECT
    m.id,
    m.name,
    m.type,
    m.size,
    m.content_key,
    m.deleted_at,
    m.updated_at,
    m.created_at,
    ss.root_shelf_id AS root_shelf_id,
    ss.id AS parent_sub_shelf_id,
    ss.name AS parent_sub_shelf_name,
    ss.prev_sub_shelf_id AS parent_sub_shelf_prev_sub_shelf_id,
    ss.path AS parent_sub_shelf_path,
    ss.deleted_at AS parent_sub_shelf_deleted_at,
    ss.updated_at AS parent_sub_shelf_updated_at,
    ss.created_at AS parent_sub_shelf_created_at
FROM "MaterialTable" m
LEFT JOIN "SubShelfTable" ss ON m.parent_sub_shelf_id = ss.id
LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id
WHERE m.id = ? AND uts.user_id = ? AND uts.permission IN ?
    AND (
        CASE
            WHEN ? = 0 THEN m.deleted_at  IS NOT NULL
            WHEN ? = 2 THEN m.deleted_at IS NULL
            ELSE true
        END
    )