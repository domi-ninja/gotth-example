-- name: HealthchecksRead :one
SELECT *
FROM healthchecks 
WHERE id = ?;

-- name: HealthchecksCreate :one
INSERT INTO healthchecks (id, created_at)
VALUES (?, ?)
RETURNING *;