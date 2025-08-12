-- name: CreatePost :one
insert into posts (id, created_at, title, body, author)
values (?, ?, ?, ?, ? )
returning *;

-- name: GetPostsPage :many
SELECT *
FROM posts 
ORDER BY created_at DESC
  LIMIT @pageSize OFFSET @pagingOffset;

-- name: GetPostById :one
SELECT *
FROM posts 
WHERE id = ?;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = ?;