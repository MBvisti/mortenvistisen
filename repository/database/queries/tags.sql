-- name: GetTags :many
SELECT
    tags.*
FROM
    tags
WHERE
    tags.id in (select unnest(@tag_ids::uuid[]));
