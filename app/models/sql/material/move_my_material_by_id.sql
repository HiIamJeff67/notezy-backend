-- name: MoveMyMaterialById
UPDATE "MaterialTable" m
SET 
    "parent_sub_shelf_id" = ?, 
    "updated_at" = NOW()
FROM "SubShelfTable" ss
JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id
WHERE
    m.id = ?
    AND m.deleted_at IS NULL
    AND m.parent_sub_shelf_id = ss.id
    AND uts.user_id = ?
    AND uts.permission IN ?
    AND EXISTS (
        SELECT 1
        FROM "SubShelfTable" dest_ss
        JOIN "UsersToShelvesTable" dest_uts ON dest_ss.root_shelf_id = dest_uts.root_shelf_id
        WHERE
            dest_ss.id = ?
            AND dest_uts.user_id = ?
            AND dest_uts.permission IN ?
    )