-- name: InsertSiteView :exec
INSERT INTO site_views (
    session_id,
    visitor_id,
    created_at,
    url_path,
    url_query,
    referrer_path,
    referrer_query,
    referrer_domain,
    page_title,
    event_type,
    event_name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: QueryViewsByDate :many
SELECT 
    sqlc.embed(site_views),
	sqlc.embed(site_sessions)
FROM site_views
JOIN site_sessions on 
	site_sessions.id=site_views.session_id
WHERE 
    site_views.created_at >= sqlc.arg(start_date)::timestamp
	AND site_views.created_at <= sqlc.arg(end_date)::timestamp
    AND event_type = 'page_view';

-- name: QueryTrafficCountsByDate :many
select 
	date_trunc('hour', created_at)::timestamp AS hour,
	count(distinct visitor_id) as visitor_count,
	count(id) as views
from 
	site_views 
where 
    created_at >= sqlc.arg(start_date)::timestamp
	AND created_at <= sqlc.arg(end_date)::timestamp
    AND event_type = 'page_view'
group by date_trunc('hour', created_at);
