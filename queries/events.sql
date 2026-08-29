-- name: ListEvents :many
SELECT * FROM events ORDER BY start_date;

-- name: GetEventBySlug :one
SELECT * FROM events WHERE slug = $1;

-- name: CreateEvent :one
INSERT INTO events (
    title, slug, start_date, end_date, type, venue, related_business_id,
    lat, lng, image, tags, description, content
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events SET
    title = $2, slug = $3, start_date = $4, end_date = $5, type = $6, venue = $7,
    lat = $8, lng = $9, image = $10, tags = $11, description = $12, content = $13,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;
