-- name: QuerySubscriberByID :one
SELECT 
    id, created_at, updated_at, email, 
    subscribed_at, referer, is_verified
FROM subscribers
WHERE id = $1;

-- name: QuerySubscriberByEmail :one
SELECT 
    id, created_at, updated_at, email, 
    subscribed_at, referer, is_verified
FROM subscribers
WHERE email = $1;

-- name: QuerySubscribers :many
SELECT 
    id, created_at, updated_at, email, 
    subscribed_at, referer, is_verified
FROM subscribers
ORDER BY created_at DESC;

-- name: QuerySubscribersPage :many
SELECT 
    id, created_at, updated_at, email, 
    subscribed_at, referer, is_verified
FROM subscribers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: QuerySubscribersCount :one
SELECT COUNT(*) FROM subscribers;

-- name: QueryVerifiedSubscribers :many
SELECT 
    id, created_at, updated_at, email, 
    subscribed_at, referer, is_verified
FROM subscribers
WHERE is_verified = true
ORDER BY created_at DESC;

-- name: InsertSubscriber :one
INSERT INTO subscribers (
    id, created_at, updated_at, email,
    subscribed_at, referer, is_verified
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateSubscriber :one
UPDATE subscribers
SET 
    updated_at = $2,
    email = $3,
    subscribed_at = $4,
    referer = $5,
    is_verified = $6
WHERE id = $1
RETURNING *;


-- name: UpdateSubscriberVerification :one
UPDATE subscribers
SET 
    updated_at = $2,
    is_verified = $3
WHERE id = $1
RETURNING *;


-- name: DeleteSubscriber :exec
DELETE FROM subscribers 
WHERE id = $1;

-- name: DeleteSubscriberByEmail :exec
DELETE FROM subscribers 
WHERE email = $1;

