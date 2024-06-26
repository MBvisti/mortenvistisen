-- name: QueryTagsByPost :many
select tags.*
from tags
left join posts_tags on posts_tags.tag_id = tags.id
where posts_tags.post_id = $1
;

-- name: QueryAllTags :many
select *
from tags
;

-- name: QueryTagsByIDs :many
select *
from tags
where id in (sqlc.arg(tag_ids)::uuid[])
;

-- name: InsertTag :exec
insert into tags (id, name)
values ($1, $2);

-- name: AssociateTagWithPost :exec
insert into posts_tags (id, post_id, tag_id) values ($1, $2, $3);

-- name: DeleteTagsFromPost :exec
delete from posts_tags
where post_id = $1
;

