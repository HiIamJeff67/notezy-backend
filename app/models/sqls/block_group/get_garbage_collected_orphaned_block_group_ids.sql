SELECT bg.id
FROM "BlockGroupTable" AS bg
INNER JOIN "BlockPackTable" bp ON bg.block_pack_id = bp.id
INNER JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id
WHERE bg.id IN ?
AND bg.deleted_at IS NULL
AND EXISTS (
    SELECT 1
    FROM "UsersToShelvesTable" uts
    WHERE uts.root_shelf_id = ss.root_shelf_id
    AND uts.user_id = ?
    AND uts.permission IN ?
)
AND NOT EXISTS (
    SELECT 1
    FROM "BlockTable" b
    WHERE b.block_group_id = bg.id
    AND b.deleted_at IS NULL
)
