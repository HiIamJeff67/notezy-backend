-- name: ResetEmail
UPDATE "UserTable" u
SET 
    email = ?
FROM "UserAccountTable" ua
WHERE ua.auth_code = ?
    AND u.user_id = ?
    AND u.user_id = ua.user_id
RETURNING u.updated_at;