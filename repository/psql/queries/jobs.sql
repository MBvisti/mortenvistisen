-- name: QueryJobs :many
select *
from jobs
limit coalesce(sqlc.narg('limit')::int, null)
offset coalesce(sqlc.narg('offset')::int, 0)
;

