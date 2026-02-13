-- name: AddTagToPost :exec
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2);

-- name: RemoveTagFromPost :exec
DELETE FROM post_tags 
WHERE post_id = $1 AND tag_id = $2;

-- name: GetPostTags :many
SELECT t.* FROM tags t
JOIN post_tags pt ON t.id = pt.tag_id
WHERE pt.post_id = $1;

-- name: GetPostsByTag :many
SELECT p.* FROM posts p
JOIN post_tags pt ON p.id = pt.post_id
WHERE pt.tag_id = $1 AND p.status = 'completed'
ORDER BY p.created_at DESC;

-- name: SetPostTags :exec
DELETE FROM post_tags WHERE post_id = $1;
