-- name: QueryTagByID :one
select * from tags where id=$1;

-- name: QueryTags :many
select * from tags;

-- name: InsertTag :one
insert into
    tags (id, created_at, updated_at, title)
values
    ($1, now(), now(), $2)
returning *;

-- name: UpdateTag :one
update tags
    set updated_at=now(), title=$2
where id = $1
returning *;

-- name: DeleteTag :exec
delete from tags where id=$1;

-- name: QueryPaginatedTags :many
select * from tags
order by created_at desc
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountTags :one
select count(*) from tags;

-- name: UpsertTag :one
insert into
    tags (id, created_at, updated_at, title)
values
    ($1, now(), now(), $2)
on conflict (id) do update set updated_at=now(), title=excluded.title
returning *;
