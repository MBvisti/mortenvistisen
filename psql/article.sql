-- name: QueryArticleByID :one
SELECT 
    p.id, p.created_at, p.updated_at, p.title, p.filename, 
    p.slug, p.excerpt, p.draft, p.released_at as release_date, 
    p.read_time
FROM posts p
WHERE p.id = $1;

-- name: QueryArticles :many
SELECT 
    p.id, p.created_at, p.updated_at, p.title, p.filename, 
    p.slug, p.excerpt, p.draft, p.released_at as release_date, 
    p.read_time
FROM posts p
ORDER BY p.created_at DESC;

-- name: QueryArticlesPage :many
SELECT 
    p.id, p.created_at, p.updated_at, p.title, p.filename, 
    p.slug, p.excerpt, p.draft, p.released_at as release_date, 
    p.read_time
FROM posts p
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: QueryArticlesCount :one
SELECT COUNT(*) FROM posts;

-- name: QueryArticleTags :many
SELECT t.id, t.name
FROM tags t
JOIN posts_tags pt ON pt.tag_id = t.id
WHERE pt.post_id = $1;

-- name: InsertArticle :one
INSERT INTO posts (
    id, created_at, updated_at, title, filename,
    slug, excerpt, draft, released_at, read_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;

-- name: InsertArticleTag :exec
INSERT INTO posts_tags (
    id, post_id, tag_id
) VALUES (
    $1, $2, $3
);

-- name: UpdateArticle :exec
UPDATE posts
SET 
    updated_at = $2,
    title = $3,
    filename = $4,
    slug = $5,
    excerpt = $6,
    draft = $7,
    released_at = $8,
    read_time = $9
WHERE id = $1;

-- name: DeleteArticle :exec
DELETE FROM posts WHERE id = $1;

-- name: DeleteArticleTags :exec
DELETE FROM posts_tags WHERE post_id = $1;

-- name: QueryArticleBySlug :one
SELECT 
    p.id, p.created_at, p.updated_at, p.title, p.filename, 
    p.slug, p.excerpt, p.draft, p.released_at as release_date, 
    p.read_time
FROM posts p
WHERE p.slug = $1;
