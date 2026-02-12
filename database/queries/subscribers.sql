-- name: QuerySubscriberByID :one
select * from subscribers where id=$1;

-- name: QuerySubscriberByEmail :one
select * from subscribers
where lower(email) = lower($1)
order by id desc
limit 1;

-- name: QuerySubscribers :many
select * from subscribers;

-- name: InsertSubscriber :one
insert into
    subscribers (created_at, updated_at, email, subscribed_at, referer, is_verified)
values
    (now(), now(), $1, $2, $3, $4)
returning *;

-- name: UpdateSubscriber :one
update subscribers
    set updated_at=now(), email=$2, subscribed_at=$3, referer=$4, is_verified=$5
where id = $1
returning *;

-- name: DeleteSubscriber :exec
delete from subscribers where id=$1;

-- name: QueryPaginatedSubscribers :many
select * from subscribers
order by created_at desc
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountSubscribers :one
select count(*) from subscribers;

-- name: UpsertSubscriber :one
insert into
    subscribers (created_at, updated_at, email, subscribed_at, referer, is_verified)
values
    (now(), now(), $1, $2, $3, $4)
on conflict (id) do update set updated_at=now(), email=excluded.email, subscribed_at=excluded.subscribed_at, referer=excluded.referer, is_verified=excluded.is_verified
returning *;
