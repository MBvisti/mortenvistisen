-- name: QueryArticleTagConnectionByID :one
select * from article_tag_connections where id=$1;

-- name: QueryArticleTagConnection :many
select * from article_tag_connections;

-- name: QueryTagsByArticleID :many
select t.id, t.created_at, t.updated_at, t.title
from tags t
inner join article_tag_connections atc on t.id = atc.tag_id
where atc.article_id = $1;

-- name: InsertArticleTagConnection :one
insert into
    article_tag_connections (id, article_id, tag_id)
values
    ($1, $2, $3)
returning *;

-- name: UpdateArticleTagConnection :one
update article_tag_connections
    set article_id=$2, tag_id=$3
where id = $1
returning *;

-- name: DeleteArticleTagConnection :exec
delete from article_tag_connections where id=$1;

-- name: CountArticleTagConnection :one
select count(*) from article_tag_connections;

-- name: UpsertArticleTagConnection :one
insert into
    article_tag_connections (id, article_id, tag_id)
values
    ($1, $2, $3)
on conflict (id) do update set article_id=excluded.article_id, tag_id=excluded.tag_id
returning *;
