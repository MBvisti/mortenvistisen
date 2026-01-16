-- name: QueryArticleByID :one
select * from articles where id=$1;

-- name: QueryArticleBySlug :one
select * from articles where slug=$1;

-- name: QueryArticles :many
select * from articles;

-- name: QueryPublishedArticles :many
select * from articles
where published = true
order by first_published_at desc;

-- name: InsertArticle :one
insert into
    articles (id, created_at, updated_at, first_published_at, title, excerpt, meta_title, meta_description, slug, image_link, read_time, content, published)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning *;

-- name: UpdateArticle :one
update articles
    set updated_at=now(), first_published_at=$2, title=$3, excerpt=$4, meta_title=$5, meta_description=$6, slug=$7, image_link=$8, read_time=$9, content=$10, published=$11
where id = $1
returning *;

-- name: DeleteArticle :exec
delete from articles where id=$1;

-- name: QueryPaginatedArticles :many
select * from articles
order by created_at desc
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountArticles :one
select count(*) from articles;

-- name: UpsertArticle :one
insert into
    articles (id, created_at, updated_at, first_published_at, title, excerpt, meta_title, meta_description, slug, image_link, read_time, content, published)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
on conflict (id) do update set updated_at=now(), first_published_at=excluded.first_published_at, title=excluded.title, excerpt=excluded.excerpt, meta_title=excluded.meta_title, meta_description=excluded.meta_description, slug=excluded.slug, image_link=excluded.image_link, read_time=excluded.read_time, content=excluded.content, published=excluded.published
returning *;
