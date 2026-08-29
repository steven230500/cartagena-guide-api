-- name: ListRoutes :many
SELECT * FROM routes ORDER BY title;

-- name: GetRouteBySlug :one
SELECT * FROM routes WHERE slug = $1;

-- name: CreateRoute :one
INSERT INTO routes (
    slug, title, description, duration, distance, difficulty, category, image,
    highlights, audio_guide, offline, price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateRoute :one
UPDATE routes SET
    slug = $2, title = $3, description = $4, duration = $5, distance = $6,
    difficulty = $7, category = $8, image = $9, highlights = $10,
    audio_guide = $11, offline = $12, price = $13, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRoute :exec
DELETE FROM routes WHERE id = $1;

-- name: ListStepsByRoute :many
SELECT * FROM route_steps WHERE route_id = $1 ORDER BY position;

-- name: CreateRouteStep :one
INSERT INTO route_steps (route_id, title, description, audio_url, duration, lat, lng, image, position)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteStepsByRoute :exec
DELETE FROM route_steps WHERE route_id = $1;
