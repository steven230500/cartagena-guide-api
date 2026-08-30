-- name: ListBusinesses :many
SELECT * FROM businesses
WHERE (sqlc.narg('type')::text IS NULL OR type = sqlc.narg('type'))
  AND (sqlc.narg('barrio')::text IS NULL OR barrio = sqlc.narg('barrio'))
  AND (sqlc.narg('q')::text IS NULL OR name ILIKE '%' || sqlc.narg('q') || '%')
ORDER BY name;

-- name: GetBusinessBySlug :one
SELECT * FROM businesses WHERE slug = $1;

-- name: GetBusinessByID :one
SELECT * FROM businesses WHERE id = $1;

-- name: ListBusinessesByOwner :many
SELECT * FROM businesses WHERE owner_id = $1 ORDER BY name;

-- name: SetBusinessOwner :one
UPDATE businesses SET owner_id = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: CreateBusiness :one
INSERT INTO businesses (
    name, slug, type, subtype, barrio, lat, lng, image, tags, description, description_en,
    hours, price_hint, price_typical_note, phone, web, email, instagram
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: UpdateBusiness :one
UPDATE businesses SET
    name = $2, slug = $3, type = $4, subtype = $5, barrio = $6, lat = $7, lng = $8,
    image = $9, tags = $10, description = $11, description_en = $12, hours = $13, price_hint = $14,
    price_typical_note = $15, phone = $16, web = $17, email = $18, instagram = $19,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBusiness :exec
DELETE FROM businesses WHERE id = $1;
