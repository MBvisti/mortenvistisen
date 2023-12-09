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
    posts.read_time,
    array_agg(tags)::uuid[] AS tag_ids
FROM
    posts
INNER JOIN
    posts_tags ON posts.id = posts_tags.post_id
RIGHT JOIN
    tags ON tags.id = posts_tags.tag_id
WHERE
    released_at IS NOT NULL && draft = false
ORDER BY
    released_at DESC
LIMIT
    5;
