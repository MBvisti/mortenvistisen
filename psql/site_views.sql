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
