-- name: UpdateAuthCodeForSendingValidationEmail
UPDATE "UserAccountTable" ua
SET
    auth_code = ?, 
    auth_code_expired_at = ?
FROM "UserTable" u
WHERE ua.user_id = u.id AND u.email = ?
RETURNING u.name, u.user_agent;
