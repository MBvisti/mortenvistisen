-- name: GetTagsForPost :many
SELECT
    tags.*
FROM
    tags
LEFT JOIN
    posts_tags ON posts_tags.tag_id = tags.id
WHERE
    posts_tags.post_id = $1;
