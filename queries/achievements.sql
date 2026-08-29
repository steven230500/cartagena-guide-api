-- name: ListAchievements :many
SELECT * FROM achievements ORDER BY threshold;

-- name: CreateAchievement :one
INSERT INTO achievements (
    code, title, description, icon, criteria_type, threshold
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateAchievement :one
UPDATE achievements SET
    code = $2, title = $3, description = $4, icon = $5,
    criteria_type = $6, threshold = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAchievement :exec
DELETE FROM achievements WHERE id = $1;
