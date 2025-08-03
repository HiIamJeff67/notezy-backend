-- name: UpdateAuthCodeForSendingValidationEmail
UPDATE "UserAccountTable" ua
SET
    auth_code = ?, 
    auth_code_expired_at = ?, 
    block_auth_code_until = ?
FROM "UserTable" u
WHERE ua.user_id = u.id AND u.email = ? AND block_auth_code_until < now()
RETURNING u.name, u.user_agent, block_auth_code_until, now();
