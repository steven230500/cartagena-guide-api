-- name: ListParishes :many
SELECT * FROM parishes ORDER BY name;

-- name: ListSchedulesByParish :many
SELECT * FROM parish_schedules WHERE parish_id = $1 ORDER BY position;

-- name: CreateParish :one
INSERT INTO parishes (name, address, neighborhood, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateParishSchedule :one
INSERT INTO parish_schedules (parish_id, day, times, position)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateParish :one
UPDATE parishes SET name = $2, address = $3, neighborhood = $4, phone = $5
WHERE id = $1
RETURNING *;

-- name: DeleteParish :exec
DELETE FROM parishes WHERE id = $1;

-- name: DeleteSchedulesByParish :exec
DELETE FROM parish_schedules WHERE parish_id = $1;
