-- name: QueryAllSubscribers :many
select * from subscribers;

-- name: DeleteSubscriber :exec
delete from subscribers where id=$1;

-- name: InsertSubscriber :one
insert into subscribers 
	(id, created_at, updated_at, email, subscribed_at, referer, is_verified) 
values 
	($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: ConfirmSubscriberEmail :exec
update subscribers set is_verified=true, updated_at=$2 where id=$1;

-- name: QuerySubscriber :one
select * from subscribers where id = $1;

-- name: QueryVerifiedSubscribers :many
select * from subscribers where is_verified = true;
