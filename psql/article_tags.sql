-- name: QueryArticleTagByID :one
select * from article_tags where id=$1;

-- name: QueryArticleTagByTitle :one
select * from article_tags where title=$1;

-- name: QueryArticleTags :many
select * from article_tags order by title asc;

-- name: QueryArticleTagsByArticleID :many
select at.* from article_tags at
join article_tag_connections atc on at.id = atc.tag_id
where atc.article_id = $1
order by at.title asc;

-- name: InsertArticleTag :one
insert into article_tags (id, created_at, updated_at, title)
values ($1, $2, $3, $4)
returning *;

-- name: UpdateArticleTag :one
update article_tags
    set updated_at=$2, title=$3
where id = $1
returning *;

-- name: DeleteArticleTag :exec
delete from article_tags where id=$1;

-- name: CountArticleTags :one
select count(*) from article_tags;