-- name: SelectUserByUUID :one
SELECT 
    user_uuid,
    user_name,
    password,
    email,
    full_name
    FROM users 
    WHERE user_uuid = ?;

-- name: SelectUserByUserName :one
SELECT 
    user_uuid,
    user_name,
    password,
    email,
    full_name
FROM users 
WHERE user_name = ?;

-- name: SelectUserByEmail :one
SELECT 
    user_uuid,
    user_name,
    password,
    email,
    full_name
FROM users
WHERE email = ?;

-- name: ShowPosts :many
SELECT 
    post_uuid,
    post_content,
    post_tags,
    user_uuid,
    likes
FROM posts; 

    