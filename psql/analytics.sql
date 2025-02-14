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
    scroll_depth,
	real_ip
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: QueryAnalyticsByDateRange :many
select 
	* 
from analytics 
where 
	website_id = $1 
and 
	timestamp between $2 and $3
order by 
	timestamp desc;

-- name: QueryDailyVisits :one
SELECT 
    COUNT(DISTINCT visitor_id) as visit_count
FROM analytics 
WHERE 
    website_id = $1 
    AND DATE(timestamp) = DATE($2);

-- name: QueryDailyViews :one
SELECT 
    COUNT(*) as view_count
FROM analytics 
WHERE 
    website_id = $1 
    AND DATE(timestamp) = DATE($2)
	AND type='pageview';
