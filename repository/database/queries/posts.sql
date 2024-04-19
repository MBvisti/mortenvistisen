-- name: GetLatestPosts :many
SELECT
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
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
    posts.header_title,
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

-- name: GetFiveRandomPosts :many
SELECT
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
FROM
    posts
WHERE
    released_at IS NOT NULL AND draft = false AND posts.id != $1
ORDER BY
    random()
limit 5;

-- name: QueryAllPosts :many
SELECT
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
FROM
    posts;
