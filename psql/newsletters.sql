-- name: QueryAllNewsletters :many
SELECT 
    id, created_at, updated_at, title,
    content, released_at, released, slug
FROM newsletters
ORDER BY created_at DESC;
;

-- name: QueryNewsletterBySlug :one
SELECT 
    id, created_at, updated_at, title,
    content, released_at, released, slug
FROM newsletters
WHERE slug=$1;

-- name: QueryNewsletters :many
SELECT 
    id, created_at, updated_at, title,
    content, released_at, released, slug
FROM newsletters
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: QueryNewslettersCount :one
SELECT COUNT(*) FROM newsletters;

-- name: QueryNewsletterByID :one
SELECT 
    id, created_at, updated_at, title,
    content, released_at, released, slug
FROM newsletters
WHERE id = $1;

-- name: InsertNewsletter :one
INSERT INTO newsletters (
    id, created_at, updated_at, title,
    content, released_at, released, slug
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;

-- name: UpdateNewsletter :exec
UPDATE newsletters
SET 
    updated_at = $2,
    title = $3,
    content = $4,
    released_at = $5,
    released = $6,
	slug = $7
WHERE id = $1;

-- name: DeleteNewsletter :exec
DELETE FROM newsletters 
WHERE id = $1;
