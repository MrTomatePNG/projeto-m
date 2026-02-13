-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
RETURNING *;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE name = $1 LIMIT 1;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = $1 LIMIT 1;

-- name: ListAllTags :many
SELECT * FROM tags
ORDER BY name ASC;

-- name: GetOrCreateTag :one
INSERT INTO tags (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;
