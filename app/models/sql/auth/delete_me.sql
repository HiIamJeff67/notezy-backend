DELETE FROM "UserTable" u
USING "UserAccountTable" ua
WHERE u.id = ?
    AND ua.user_id = u.id
    AND (
        u.role = 'Guest'
        OR (
            AND ua.auth_code = ?
            AND ua.auth_code_expired_at > NOW()
            AND u.role <> 'Guest'
        )
    )
RETURNING NOW() AS deleted_at;