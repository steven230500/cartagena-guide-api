-- name: ListFavoriteBusinesses :many
SELECT b.* FROM businesses b
JOIN user_favorites f ON f.business_id = b.id
WHERE f.user_id = $1
ORDER BY f.created_at DESC;

-- name: AddFavorite :exec
INSERT INTO user_favorites (user_id, business_id) VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveFavorite :exec
DELETE FROM user_favorites WHERE user_id = $1 AND business_id = $2;
