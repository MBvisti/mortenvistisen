-- name: QueryNewsletters :many
select
    newsletters.id as newsletter_id,
    newsletters.created_at as newsletter_created_at,
    newsletters.updated_at as newsletter_updated_at,
    newsletters.title as newsletter_title,
    newsletters.content as newsletter_content,
    newsletters.released as newsletter_released,
    newsletters.released_at as newsletter_released_at
from newsletters
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
	(id, created_at, updated_at, title, content, released, released_at)
values 
	($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: QueryNewsletterByID :one
select *
from newsletters
where id = $1
;

-- name: UpdateNewsletter :one
update newsletters
	set updated_at = $1, title = $2, content = $3, released = $4, released_at = $5
where id = $6
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

