-- name: QueryPosts :many
select posts.*
from posts
order by posts.released_at desc
limit coalesce(sqlc.narg('limit')::int, null)
offset coalesce(sqlc.narg('offset')::int, 0)
;

-- name: GetLatestPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
from posts
where released_at is not null and draft = false
order by released_at desc
;

-- name: GetPostBySlug :one
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
from posts
where slug = $1
;

-- name: GetFiveRandomPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time
from posts
where released_at is not null and draft = false and posts.id != $1
order by random()
limit 5
;

-- name: QueryPostsInPages :many
select posts.*
from posts
limit $1
offset $2
;

-- name: QueryAllPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.header_title,
    posts.filename,
    posts.slug,
    posts.excerpt,
    posts.draft,
    posts.released_at,
    posts.read_time,
    (select count(id) from posts) as total_posts_count
from posts
limit 7
offset $1
;

-- name: QueryAllFilenames :many
select filename
from posts
;

-- name: InsertPost :one
insert into posts (id, created_at, updated_at, title, header_title, filename, slug, excerpt, draft, released_at, read_time)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning id;

-- name: UpdatePost :one
update posts
    set updated_at = $1, title = $2, header_title = $3, slug = $4, excerpt = $5, draft = $6, released_at = $7, read_time = $8
where id = $9
returning *;

-- name: QueryPostByID :one
select *
from posts
where id = $1
;

-- name: QueryPostBySlug :one
select *
from posts
where slug = $1
;

-- name: QueryPostsCount :one
select count(id)
from posts
;

