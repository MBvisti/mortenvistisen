-- name: QueryArticleByID :one
select * from articles where id=$1;

-- name: QueryArticleByTitle :one
select * from articles where title=$1;

-- name: QueryArticleBySlug :one
select * from articles where slug=$1;

-- name: QueryArticles :many
select * from articles order by created_at desc;

-- name: QueryPublishedArticles :many
select * from articles where is_published=true order by first_published_at desc;

-- name: QueryDraftArticles :many
select * from articles where is_published=false order by created_at desc;

-- name: QueryArticlesPaginated :many
select * from articles 
order by created_at desc 
limit $1 offset $2;

-- name: CountArticles :one
select count(*) from articles;

-- name: InsertArticle :one
insert into
    articles (id, created_at, updated_at, title, excerpt, meta_title, meta_description, slug, image_link, content, read_time)
values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning *;

-- name: UpdateArticle :one
update articles
    set updated_at=$2, title=$3, excerpt=$4, meta_title=$5, meta_description=$6, slug=$7, image_link=$8, content=$9, read_time=$10, is_published=$11
where id = $1
returning *;

-- name: UpdateArticleContent :one
update articles
    set updated_at=$2, content=$3
where id = $1
returning *;

-- name: UpdateArticleMetadata :one
update articles
    set updated_at=$2, title=$3, excerpt=$4, meta_title=$5, meta_description=$6, slug=$7, image_link=$8
where id = $1
returning *;

-- name: PublishArticle :one
update articles
    set updated_at=$2, is_published=$3, first_published_at=$4
where id = $1
returning *;

-- -- name: UnpublishArticle :one
-- update articles
--     set updated_at=$2, is_published=false
-- where id = $1
-- returning *;

-- name: DeleteArticle :exec
delete from articles where id=$1;
