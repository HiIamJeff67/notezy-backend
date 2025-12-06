-- name: GetMyBlockGroupAndItsBlocksById
SELECT
    bg.id AS block_group_id,
    bg.block_pack_id,
    bg.prev_block_group_id,
    bg.sync_block_group_id,
    bg.mega_byte_size,
    bg.deleted_at AS block_group_deleted_at,
    bg.updated_at AS block_group_updated_at,
    bg.created_at AS block_group_created_at,
    b.id AS block_id,
    b.parent_block_id,
    b.type AS block_type,
    b.props AS block_props,
    b.content AS block_content,
    b.deleted_at AS block_deleted_at,
    b.updated_at AS block_updated_at,
    b.created_at AS block_created_at
FROM "BlockGroupTable" bg
-- left join the block table so that if the block group exists but the blocks under it does not exist
-- then we can still get the meta data of the block group
LEFT JOIN "BlockTable" b ON bg.id = b.block_group_id
JOIN "BlockPackTable" bp ON bg.block_pack_id = bp.id
JOIN "SubShelfTable" ss ON bp.parent_sub_shelf_id = ss.id
JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id
WHERE
    bg.id = $1
    AND uts.user_id = $2
    AND uts.permission = ANY($3::"AccessControlPermission"[])
    AND (
        CASE
            WHEN $4 = 0 THEN (
                b.deleted_at IS NOT NULL 
                AND bg.deleted_at IS NOT NULL
            )
            WHEN $4 = 2 THEN (
                b.deleted_at IS NULL 
                AND bg.deleted_at IS NULL
            )
            ELSE true
        END
    );