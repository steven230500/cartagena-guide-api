-- name: ListPlans :many
SELECT * FROM plans ORDER BY title;

-- name: CreatePlan :one
INSERT INTO plans (title, description, type, price, price_amount, date, time, location, neighborhood)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdatePlan :one
UPDATE plans SET
    title = $2, description = $3, type = $4, price = $5, price_amount = $6,
    date = $7, time = $8, location = $9, neighborhood = $10
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans WHERE id = $1;
