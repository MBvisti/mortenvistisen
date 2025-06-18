-- name: QueryNewsletterByID :one
select * from newsletters where id=$1;

-- name: QueryNewsletterByTitle :one
select * from newsletters where title=$1;

-- name: QueryNewsletterBySlug :one
select * from newsletters where slug=$1;

-- name: QueryNewsletters :many
select * from newsletters order by created_at desc;

-- name: QueryPublishedNewsletters :many
select * from newsletters where is_published=true order by released_at desc;

-- name: QueryDraftNewsletters :many
select * from newsletters where is_published=false order by created_at desc;

-- name: QueryNewslettersPaginated :many
select * from newsletters 
order by created_at desc 
limit $1 offset $2;

-- name: CountNewsletters :one
select count(*) from newsletters;

-- name: InsertNewsletter :one
insert into
    newsletters (id, created_at, updated_at, title, slug, content)
values
    ($1, $2, $3, $4, $5, $6)
returning *;

-- name: UpdateNewsletter :one
update newsletters
    set updated_at=$2, title=$3, slug=$4, content=$5, is_published=$6
where id = $1
returning *;

-- name: UpdateNewsletterContent :one
update newsletters
    set updated_at=$2, content=$3
where id = $1
returning *;

-- name: PublishNewsletter :one
update newsletters
    set updated_at=$2, is_published=$3, released_at=$4
where id = $1
returning *;

-- name: DeleteNewsletter :exec
delete from newsletters where id=$1;

-- name: QueryNewslettersReadyToSend :many
select * from newsletters where send_status='ready_to_send' order by created_at asc;

-- name: MarkNewsletterReadyToSend :one
update newsletters
    set updated_at=$2, send_status='ready_to_send', total_recipients=$3
where id = $1
returning *;

-- name: UpdateNewsletterSendStatus :one
update newsletters
    set updated_at=$2, send_status=$3, sending_started_at=$4, sending_completed_at=$5, emails_sent=$6
where id = $1
returning *;