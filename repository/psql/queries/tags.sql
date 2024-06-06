-- name: GetTagsForPost :many
select tags.*
from tags
left join posts_tags on posts_tags.tag_id = tags.id
where posts_tags.post_id = $1
;

-- name: QueryAllTags :many
select *
from tags
;

-- name: InsertTag :exec
insert into tags (id, name)
values ($1, $2);

-- name: AssociateTagWithPost :exec
insert into posts_tags (id, post_id, tag_id) values ($1, $2, $3);

