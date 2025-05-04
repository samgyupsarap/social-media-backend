-- name: UpdateUser :exec
UPDATE users SET full_name = ?, email = ?, user_name = ?, password = ?, profile_picture = ? WHERE user_uuid = ?;
