-- name: CreateBusinessClaim :one
INSERT INTO business_claims (business_id, user_id, message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBusinessClaim :one
SELECT * FROM business_claims WHERE id = $1;

-- name: ResolveBusinessClaim :one
UPDATE business_claims SET status = $2, resolved_at = now() WHERE id = $1
RETURNING *;

-- name: ListBusinessClaims :many
SELECT
    bc.id, bc.business_id, bc.user_id, bc.message, bc.status, bc.created_at, bc.resolved_at,
    b.name AS business_name, b.slug AS business_slug, u.email AS user_email
FROM business_claims bc
JOIN businesses b ON b.id = bc.business_id
JOIN users u ON u.id = bc.user_id
WHERE (sqlc.narg('status')::text IS NULL OR bc.status = sqlc.narg('status'))
  AND (sqlc.narg('user_id')::uuid IS NULL OR bc.user_id = sqlc.narg('user_id'))
ORDER BY bc.created_at DESC;
