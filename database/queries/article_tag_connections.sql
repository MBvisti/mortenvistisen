-- name: QueryArticleTagConnectionByID :one
select * from article_tag_connections where id=$1;

-- name: QueryArticleTagConnection :many
select * from article_tag_connections;

-- name: InsertArticleTagConnection :one
insert into
    article_tag_connections (article_id, tag_id)
values
    ($1, $2)
returning *;

-- name: UpdateArticleTagConnection :one
update article_tag_connections
    set article_id=$2, tag_id=$3
where id = $1
returning *;

-- name: DeleteArticleTagConnection :exec
delete from article_tag_connections where id=$1;

-- name: UpsertArticleTagConnection :one
insert into
    article_tag_connections (article_id, tag_id)
values
    ($1, $2)
on conflict (id) do update set article_id=excluded.article_id, tag_id=excluded.tag_id
returning *;
