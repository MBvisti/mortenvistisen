-- name: StoreToken :exec
insert into tokens
    (id, created_at, hash, expires_at, scope, user_id) values ($1, $2, $3, $4, $5, $6) 
returning *;

-- name: QueryTokenByHash :one
select * from tokens where hash=$1;

-- name: DeleteToken :exec
delete from tokens where id=$1;

-- name: DeleteSubscriberTokenBySubscriberID :exec
delete from subscriber_tokens where subscriber_id=$1;

-- name: InsertSubscriberToken :exec
insert into subscriber_tokens
    (id, created_at, hash, expires_at, scope, subscriber_id) values ($1, $2, $3, $4, $5, $6);

-- name: QuerySubscriberTokenByHash :one
select * from subscriber_tokens where hash=$1;
