-- name: ValidateEmailByAuthCode
UPDATE UserTable u
SET u.Role = 'Normal'
FROM UserAccountTable ua
WHERE ua.user_id = u.id
    AND u.id = ?
    AND ua.auth_code = ?
    AND ua.auth_code_expired_at > NOW()
RETURNING u.updated_at;