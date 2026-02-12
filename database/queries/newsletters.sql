-- name: QueryNewsletterByID :one
select * from newsletters where id=$1;

-- name: QueryNewsletters :many
select * from newsletters;

-- name: QueryPublishedNewsletters :many
select * from newsletters where is_published=true order by released_at desc;

-- name: InsertNewsletter :one
insert into
    newsletters (created_at, updated_at, title, slug, meta_title, meta_description, is_published, released_at, content)
values
    (now(), now(), $1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: UpdateNewsletter :one
update newsletters
    set updated_at=now(), title=$2, slug=$3, meta_title=$4, meta_description=$5, is_published=$6, released_at=$7, content=$8
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
    newsletters (created_at, updated_at, title, slug, meta_title, meta_description, is_published, released_at, content)
values
    (now(), now(), $1, $2, $3, $4, $5, $6, $7)
on conflict (id) do update set updated_at=now(), title=excluded.title, slug=excluded.slug, meta_title=excluded.meta_title, meta_description=excluded.meta_description, is_published=excluded.is_published, released_at=excluded.released_at, content=excluded.content
returning *;
