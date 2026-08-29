-- name: GetRouteProgress :one
SELECT * FROM user_route_progress WHERE user_id = $1 AND route_id = $2;

-- name: UpsertRouteProgress :one
INSERT INTO user_route_progress (user_id, route_id, current_step)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, route_id)
DO UPDATE SET current_step = $3, updated_at = now()
RETURNING *;

-- name: CountCompletedRoutes :one
SELECT count(*) FROM user_route_progress p
JOIN (
    SELECT route_id, count(*) AS total_steps FROM route_steps GROUP BY route_id
) s ON s.route_id = p.route_id
WHERE p.user_id = $1 AND p.current_step >= s.total_steps - 1;
