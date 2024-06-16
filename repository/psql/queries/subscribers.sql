-- name: InsertSubscriber :one
insert into subscribers 
	(id, created_at, updated_at, email, subscribed_at, referer, is_verified) 
values 
	($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: QuerySubscribers :many
select *
from subscribers
where is_verified = coalesce(sqlc.narg('is_verified')::bool, null)
limit coalesce(sqlc.narg('limit')::int, null)
offset coalesce(sqlc.narg('offset')::int, 0)
;

-- name: QuerySubscriberByID :one
select *
from subscribers
where id = $1
;

-- name: QuerySubscriberByEmail :one
select *
from subscribers
where email = $1
;

-- name: QueryVerifiedSubscribers :many
select *
from subscribers
where is_verified = true
;

-- name: QuerySubscriberCount :one
select count(id)
from subscribers
;

-- name: QuerySubscriberCountByStatus :one
select count(id)
from subscribers
where is_verified = $1
;

-- name: QueryNewSubscribersInCurrentMonth :many
select *
from subscribers
where date_trunc('month', created_at) = date_trunc('month', current_timestamp)
;

-- name: UpdateSubscriber :one
update subscribers set updated_at = $1, email = $2, subscribed_at = $3, referer = $4, is_verified = $5 where id = $1
returning *
;

-- name: DeleteSubscriber :exec
delete from subscribers
where id = $1
;

