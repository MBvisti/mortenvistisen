-- name: QueryAnalyticsByID :one
select * from analytics where id = $1;

-- name: QueryAnalyticsByWebsiteID :many
select * from analytics where website_id = $1;

-- name: QueryAnalytics :many
select * from analytics;

-- name: InsertAnalytic :exec
insert into analytics (
    id,
    website_id,
    type,
    url,
    path,
    referrer,
    title,
    timestamp,
    screen,
    language,
    visitor_id,
    session_id,
    scroll_depth
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: DeleteAnalytics :exec
delete from analytics where id = $1;

-- name: DeleteAnalyticsByWebsiteID :exec
delete from analytics where website_id = $1;

-- name: QueryAnalyticsByVisitorID :many
select * from analytics where visitor_id = $1;

-- name: QueryAnalyticsBySessionID :many
select * from analytics where session_id = $1;

-- name: QueryAnalyticsByDateRange :many
select * from analytics 
where website_id = $1 
and timestamp between $2 and $3
order by timestamp desc;

