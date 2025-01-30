-- name: QueryFirstUser :one
select * from users order by created_at asc limit 1;

-- name: QueryUserByID :one
select * from users where id=$1;

-- name: QueryUserByMail :one
select * from users where mail=$1;

-- name: QueryUsers :many
select * from users;

-- name: InsertUser :one
insert into
    users (id, created_at, updated_at, mail, password)
values
    ($1, $2, $3, $4, $5)
returning *;

-- name: UpdateUser :one
update users
    set updated_at=$2, mail=$3
where id = $1
returning *;

-- name: DeleteUser :exec
delete from users where id=$1;

-- name: ChangeUserPassword :exec
update users set updated_at=$2, password=$3 where id=$1;

-- name: VerifyUserMail :exec
update users set updated_at=$2, mail_verified_at=$3 where mail=$1;

-- name: UpdateUserIsAdmin :one
UPDATE users 
SET 
    updated_at = $2
WHERE id = $1
RETURNING *;
