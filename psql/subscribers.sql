-- name: QuerySubscriberByID :one
select * from subscribers where id=$1;

-- name: QuerySubscriberByEmail :one
select * from subscribers where email=$1;

-- name: QuerySubscribers :many
select * from subscribers order by created_at desc;

-- name: QueryVerifiedSubscribers :many
select * from subscribers where is_verified = true order by created_at desc;

-- name: InsertSubscriber :one
insert into
    subscribers (id, created_at, updated_at, email, subscribed_at, referer, is_verified)
values
    ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: UpdateSubscriber :one
update subscribers
    set updated_at=$2, email=$3, referer=$4
where id = $1
returning *;

-- name: DeleteSubscriber :exec
delete from subscribers where id=$1;

-- name: VerifySubscriber :exec
update subscribers set updated_at=$2, is_verified=$3 where id=$1;

-- name: CountSubscribers :one
select count(*) from subscribers;

-- name: CountVerifiedSubscribers :one
select count(*) from subscribers where is_verified = true;