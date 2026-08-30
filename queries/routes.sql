-- name: ListRoutes :many
SELECT * FROM routes ORDER BY title;

-- name: GetRouteBySlug :one
SELECT * FROM routes WHERE slug = $1;

-- name: CreateRoute :one
INSERT INTO routes (
    slug, title, title_en, description, description_en, duration, distance, difficulty, category, image,
    highlights, audio_guide, offline, price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: UpdateRoute :one
UPDATE routes SET
    slug = $2, title = $3, title_en = $4, description = $5, description_en = $6, duration = $7, distance = $8,
    difficulty = $9, category = $10, image = $11, highlights = $12,
    audio_guide = $13, offline = $14, price = $15, updated_at = now()
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
