-- name: UpdateUser :exec
UPDATE users SET full_name = ?, email = ?, user_name = ?, password = ?, profile_picture = ? WHERE user_uuid = ?;

-- name: UpdatePost :exec
UPDATE posts SET post_content = ?, post_tags = ? WHERE post_uuid = ?;

-- name: UpdatePostLikes :exec
UPDATE posts SET likes = ? WHERE post_uuid = ?;