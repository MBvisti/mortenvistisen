-- name: GetLatestPosts :many
SELECT
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
FROM
    posts
WHERE
    released_at IS NOT NULL AND draft = false
ORDER BY
    released_at DESC;

-- name: GetPostBySlug :one
SELECT
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
FROM
    posts
WHERE
    slug = $1;

-- name: GetFiveRandomPosts :many
SELECT
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
FROM
    posts
WHERE
    released_at IS NOT NULL AND draft = false AND posts.id != $1
ORDER BY
    random()
limit 5;

-- name: QueryPosts :many
select
	*
from
	posts
limit 
	coalesce(sqlc.narg('limit')::int, null)
offset 
	coalesce(sqlc.narg('offset')::int, 0);

-- name: QueryPostsInPages :many
SELECT
    posts.*
FROM
	posts
LIMIT
    $1
OFFSET
    $2;

-- name: QueryAllPosts :many
SELECT
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
FROM
    posts
LIMIT
    7
OFFSET
    $1;

-- name: QueryAllFilenames :many
select filename from posts;

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
select * from posts where id = $1;

-- name: QueryPostBySlug :one
select * from posts where slug = $1;

-- name: QueryPostsCount :one
select count(id) from posts;
