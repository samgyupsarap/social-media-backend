-- name: SelectUserByUUID :one
SELECT * FROM users WHERE user_uuid = ?;

-- name: SelectUserByUserName :one
SELECT 
    user_uuid,
    user_name,
    password,
    email,
    full_name
FROM users 
WHERE user_name = ?;