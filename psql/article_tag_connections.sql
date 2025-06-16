-- name: QueryArticleTagConnectionByID :one
select * from article_tag_connections where id=$1;

-- name: QueryArticleTagConnectionsByArticleID :many
select * from article_tag_connections where article_id=$1;

-- name: QueryArticleTagConnectionsByTagID :many
select * from article_tag_connections where tag_id=$1;

-- name: QueryArticleTagConnection :one
select * from article_tag_connections where article_id=$1 and tag_id=$2;

-- name: InsertArticleTagConnection :one
insert into article_tag_connections (id, article_id, tag_id)
values ($1, $2, $3)
returning *;

-- name: DeleteArticleTagConnection :exec
delete from article_tag_connections where id=$1;

-- name: DeleteArticleTagConnectionByArticleAndTag :exec
delete from article_tag_connections where article_id=$1 and tag_id=$2;

-- name: DeleteArticleTagConnectionsByArticleID :exec
delete from article_tag_connections where article_id=$1;

-- name: DeleteArticleTagConnectionsByTagID :exec
delete from article_tag_connections where tag_id=$1;