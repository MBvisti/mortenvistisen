-- name: QueryLatestPosts :many
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

-- name: QueryPostBySlug :one
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
    posts.read_time,
    ARRAY_AGG(tags.name)::text[] as tags
FROM
    posts
JOIN 
    posts_tags ON posts_tags.post_id = posts.id
JOIN
    tags on tags.id = posts_tags.tag_id
WHERE
    slug = $1
GROUP BY
    posts.id;

-- name: QueryAllPost :many
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
    posts.read_time,
    ARRAY_AGG(tags.name)::text[] as tags
FROM 
    posts
JOIN 
    posts_tags ON posts_tags.post_id = posts.id
JOIN
    tags on tags.id = posts_tags.tag_id
GROUP BY
    posts.id;
