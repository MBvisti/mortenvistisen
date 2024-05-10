-- name: QueryAllNewsletters :many
select * from newsletters;

-- name: QueryNewsletterInPages :many
select 
	newsletters.*
from newsletters 
limit
	7
offset
	$1;

-- name: QueryReleasedNewslettersCount :one
select 
	count(id) as newsletters_count 
from 
	newsletters
where released = true;

-- name: InsertNewsletter :one
insert into newsletters
	(id, created_at, updated_at, title, edition, released, released_at, body, associated_article_id)
values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: QueryNewsletterByID :one
select * from newsletters where id = $1;
