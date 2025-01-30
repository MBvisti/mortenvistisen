-- name: QueryTagByID :one
SELECT id, name
FROM tags
WHERE id = $1;

-- name: QueryTags :many
SELECT id, name
FROM tags
ORDER BY name ASC;

-- name: QueryTagsByPostID :many
SELECT t.id, t.name
FROM tags t
JOIN posts_tags pt ON pt.tag_id = t.id
WHERE pt.post_id = $1
ORDER BY t.name ASC;

-- name: InsertTag :one
INSERT INTO tags (
    id, name
) VALUES (
    $1, $2
)
RETURNING id;

-- name: UpdateTag :exec
UPDATE tags
SET name = $2
WHERE id = $1;

-- name: DeleteTag :exec
DELETE FROM tags 
WHERE id = $1;

-- name: DeleteTagFromPosts :exec
DELETE FROM posts_tags 
WHERE tag_id = $1;

