-- name: GetLatestPosts :many
SELECT
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
FROM
    posts
WHERE
    released_at IS NOT NULL AND draft = false
ORDER BY
    released_at DESC;

-- name: GetPostBySlug :one
SELECT
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
FROM
    posts
WHERE
    slug = $1;
