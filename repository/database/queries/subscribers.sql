-- name: QueryAllSubscribers :many
select * from subscribers;

-- name: DeleteSubscriber :exec
delete from subscribers where id=$1;
