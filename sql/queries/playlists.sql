-- name: CreatePlayList :one
INSERT INTO playlists (id, created_at, updated_at, name, owner_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;