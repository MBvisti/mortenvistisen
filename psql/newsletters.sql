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
    newsletters (id, created_at, updated_at, title, slug, content, is_published, released_at)
values
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: UpdateNewsletter :one
update newsletters
    set updated_at=$2, title=$3, slug=$4, content=$5, is_published=$6, released_at=$7
where id = $1
returning *;

-- name: UpdateNewsletterContent :one
update newsletters
    set updated_at=$2, content=$3
where id = $1
returning *;

-- name: DeleteNewsletter :exec
delete from newsletters where id=$1;
