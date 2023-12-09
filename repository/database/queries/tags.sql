-- name: GetTags :many
SELECT
    tags.*
FROM
    tags
WHERE
    tags.id = ANY($1);
