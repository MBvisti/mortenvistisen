-- name: InsertSiteSession :one
INSERT INTO site_sessions (
    id,
    created_at,
    hostname,
    browser,
    os,
    device,
    screen,
    lang,
    country,
    subdivision1,
    subdivision2,
    city,
	finger
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: QuerySiteSession :one
SELECT * FROM site_sessions 
WHERE id = $1;
