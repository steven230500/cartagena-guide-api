-- name: ListEvents :many
SELECT * FROM events ORDER BY start_date;

-- name: GetEventBySlug :one
SELECT * FROM events WHERE slug = $1;

-- name: CreateEvent :one
INSERT INTO events (
    title, title_en, slug, start_date, end_date, type, venue, related_business_id,
    lat, lng, image, tags, description, description_en, content
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events SET
    title = $2, title_en = $3, slug = $4, start_date = $5, end_date = $6, type = $7, venue = $8,
    lat = $9, lng = $10, image = $11, tags = $12, description = $13, description_en = $14, content = $15,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;
