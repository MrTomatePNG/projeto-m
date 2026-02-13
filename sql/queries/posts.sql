-- name: ListPostsWithTags :many
SELECT
    p.id,
    p.comment,
    p.media_url,
    p.thumb_url,
    p.status,
    u.username AS author_name,
    COALESCE(
        json_agg(
            json_build_object('id', t.id, 'name', t.name)
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'
    ) AS tags
FROM posts p
JOIN users u ON p.user_id = u.id
LEFT JOIN post_tags pt ON p.id = pt.post_id
LEFT JOIN tags t ON pt.tag_id = t.id
WHERE p.status = 'completed'
  AND p.created_at < $1  -- cursor: último created_at da página anterior
GROUP BY p.id, u.username, p.created_at
ORDER BY p.created_at DESC
LIMIT $2;


-- name: CreatePost :one
INSERT INTO posts (user_id,media_type,media_hash)
VALUES ($1,$2,$3)
RETURNING * ;

-- name: UpdatePostProgress :exec
UPDATE posts
SET status = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2 AND status = 'processing';

-- name: UpdatePostMedia :exec
UPDATE posts
SET media_url = $2, 
    thumb_url = $3,
    status = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetPostByID :one
SELECT p.*, u.username as author_name
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.id = $1;

-- name: GetPostsByUserID :many
SELECT * FROM posts
WHERE user_id = $1 AND status = 'completed'
ORDER BY created_at DESC
LIMIT $2;

-- name: GetPendingPosts :many
SELECT * FROM posts
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1;

-- name: GetProcessingPosts :many
SELECT * FROM posts
WHERE status = 'processing'
AND updated_at < NOW() - INTERVAL '30 minutes'
ORDER BY updated_at ASC;
