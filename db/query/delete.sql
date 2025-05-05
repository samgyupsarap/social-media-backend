-- name: DeletePost :exec
DELETE FROM posts WHERE post_uuid = ?;

-- name: DeleteComment :exec
DELETE FROM comments WHERE comment_uuid = ?;