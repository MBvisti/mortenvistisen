-- name: QueryNewsletters :many
select
    newsletters.id as newsletter_id,
    newsletters.created_at as newsletter_created_at,
    newsletters.updated_at as newsletter_updated_at,
    newsletters.title as newsletter_title,
    newsletters.edition as newsletter_edition,
    newsletters.released as newsletter_released,
    newsletters.released_at as newsletter_released_at,
    newsletters.body as newsletter_body,
    newsletters.associated_article_id as newsletter_associated_article_id,
    posts.id as post_id,
    posts.created_at as post_created_at,
    posts.updated_at as post_updated_at,
    posts.title as post_title,
    posts.header_title as post_header_title,
    posts.filename as post_filename,
    posts.slug as post_slug,
    posts.excerpt as post_excerpt,
    posts.draft as post_draft,
    posts.released_at as post_released_at,
    posts.read_time as post_read_time
from newsletters
join posts on posts.id = newsletters.associated_article_id
where
    (
        newsletters.released = sqlc.narg('is_released')::bool
        or sqlc.narg('is_released')::bool is null
    )
limit coalesce(sqlc.narg('limit')::int, null)
offset coalesce(sqlc.narg('offset')::int, 0)
;

-- name: QueryNewsletterInPages :many
select newsletters.*
from newsletters
limit 7
offset $1
;

-- name: QueryNewslettersCount :one
select count(id)
from newsletters
;

-- name: QueryReleasedNewslettersCount :one
select count(id) as newsletters_count
from newsletters
where released = true
;

-- name: InsertNewsletter :one
insert into newsletters
	(id, created_at, updated_at, title, edition, released, released_at, body, associated_article_id)
values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: QueryNewsletterByID :one
select *
from newsletters
where id = $1
;

-- name: UpdateNewsletter :one
update newsletters
	set updated_at = $1, title = $2, edition = $3, released = $4, released_at = $5, body = $6, associated_article_id = $7
where id = $8
returning *;


-- name: CountNewsletters :one
select count(id)
from newsletters
;

-- name: CountReleasedNewsletters :one
select count(id)
from newsletters
where released=true
;

