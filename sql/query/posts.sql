-- name: CreatePost :one
insert into posts (id, created_at, title, body, user_id)
values (?, ?, ?, ?, ? )
returning *;

-- name: GetPostsPage :many
SELECT posts.*, users.email
FROM posts JOIN users ON posts.user_id = users.id
ORDER BY posts.created_at DESC
  LIMIT @pageSize OFFSET @pagingOffset;

-- name: GetPostById :one
SELECT posts.*, users.email
FROM posts JOIN users ON posts.user_id = users.id
WHERE posts.id = ?;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = ?;