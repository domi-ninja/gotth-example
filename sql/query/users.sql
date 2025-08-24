-- name: CreateUser :one
insert into users (id, created_at, email, password_hash)
values (?, ?, ?, ?)
returning *;

-- name: GetUserByEmail :one
SELECT *
FROM users 
WHERE email = ?;

-- name: GetUserById :one
SELECT *
FROM users 
WHERE id = ?;

-- name: GetUserByEmailWithPassword :one
SELECT *
FROM users 
WHERE email = ?;
