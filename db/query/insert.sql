-- name: CreateUser :exec
INSERT INTO users (user_uuid, full_name, email, user_name, password, profile_picture)
VALUES (?, ?, ?, ?, ?, ?)