-- name: QueryProjectByID :one
select * from projects where id=$1;

-- name: QueryProjectBySlug :one
select * from projects where slug=$1;

-- name: QueryProjects :many
select * from projects;

-- name: QueryPublishedProjects :many
select * from projects where published=true order by started_at desc;

-- name: InsertProject :one
insert into
    projects (
      created_at,
      updated_at,
      published,
      title,
      slug,
      started_at,
      status,
      description,
      content,
      project_url
    )
values
    (
      now(),
      now(),
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8
    )
returning *;

-- name: UpdateProject :one
update projects
set
    updated_at=now(),
    published=$2,
    title=$3,
    slug=$4,
    started_at=$5,
    status=$6,
    description=$7,
    content=$8,
    project_url=$9
where id = $1
returning *;

-- name: DeleteProject :exec
delete from projects where id=$1;

-- name: QueryPaginatedProjects :many
select * from projects
order by created_at desc
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountProjects :one
select count(*) from projects;

-- name: UpsertProject :one
insert into
    projects (
      created_at,
      updated_at,
      published,
      title,
      slug,
      started_at,
      status,
      description,
      content,
      project_url
    )
values
    (
      now(),
      now(),
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8
    )
on conflict (slug) do update
set
    updated_at=now(),
    published=excluded.published,
    title=excluded.title,
    started_at=excluded.started_at,
    status=excluded.status,
    description=excluded.description,
    content=excluded.content,
    project_url=excluded.project_url
returning *;
