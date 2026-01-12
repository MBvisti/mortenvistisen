-- name: QueryNewsletterByID :one
select * from newsletters where id=$1;

-- name: QueryNewsletterBySlug :one
select * from newsletters where slug=$1;

-- name: QueryNewsletters :many
select * from newsletters;

-- name: InsertNewsletter :one
insert into
    newsletters (id, created_at, updated_at, title, meta_title, meta_description, is_published, released_at, slug, content)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: UpdateNewsletter :one
update newsletters
    set updated_at=now(), title=$2, meta_title=$3, meta_description=$4, is_published=$5, released_at=$6, slug=$7, content=$8
where id = $1
returning *;

-- name: DeleteNewsletter :exec
delete from newsletters where id=$1;

-- name: QueryPaginatedNewsletters :many
select * from newsletters
order by created_at desc
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountNewsletters :one
select count(*) from newsletters;

-- name: UpsertNewsletter :one
insert into
    newsletters (id, created_at, updated_at, title, meta_title, meta_description, is_published, released_at, slug, content)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8)
on conflict (id) do update set updated_at=now(), title=excluded.title, meta_title=excluded.meta_title, meta_description=excluded.meta_description, is_published=excluded.is_published, released_at=excluded.released_at, slug=excluded.slug, content=excluded.content
returning *;
