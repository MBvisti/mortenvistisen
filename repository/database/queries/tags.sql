-- name: GetTagsForPost :many
SELECT
    tags.*
FROM
    tags
LEFT JOIN
    posts_tags ON posts_tags.tag_id = tags.id
WHERE
    posts_tags.post_id = $1;

-- name: QueryAllTags :many
select * from tags;

-- name: InsertTag :exec
insert into tags (id, name)
values ($1, $2);

-- name: AssociateTagWithPost :exec
insert into posts_tags (id, post_id, tag_id) values ($1, $2, $3);
