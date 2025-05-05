-- name: DeletePost :exec
DELETE FROM posts WHERE post_uuid = ?;