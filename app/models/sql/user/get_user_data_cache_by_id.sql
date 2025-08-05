-- name: GetUserDataCacheById
SELECT 
    u.public_id,
    u.name, 
    u.display_name,
    u.email,
    u.role,
    u.plan,
    u.status,
    ui.avatar_url,
    us.language,
    us.general_setting_code,
    us.privacy_setting_code,
    u.created_at,
    NOW() AS updated_at
FROM "UserTable" u
LEFT JOIN "UserInfoTable" ui ON u.id = ui.user_id
LEFT JOIN "UserSettingTable" us ON u.id = us.user_id
WHERE u.id = ?