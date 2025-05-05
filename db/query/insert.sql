-- name: CreateUser :exec
INSERT INTO users (user_uuid, full_name, email, user_name, password, profile_picture)
VALUES (?, ?, ?, ?, ?, ?);

-- name: CreatePost :exec
INSERT INTO posts (post_uuid, post_content, post_tags, user_uuid, likes)
VALUES (?, ?, ?, ?, ?);