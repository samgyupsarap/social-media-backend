-- name: CreateUser :exec
INSERT INTO users (user_uuid, full_name, email, user_name, password, profile_picture)
VALUES (?, ?, ?, ?, ?, ?);

-- name: CreatePost :exec
INSERT INTO posts (post_uuid, post_content, post_tags, user_uuid, likes)
VALUES (?, ?, ?, ?, ?);

-- name: CreateComment :exec
INSERT INTO comments (comment_uuid, comment_content, user_uuid, post_uuid)
VALUES (?, ?, ?, ?);