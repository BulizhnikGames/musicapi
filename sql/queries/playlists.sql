-- name: CreatePlayList :one
INSERT INTO playlists (id, created_at, updated_at, name, owner_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUsersPlaylists :many
SELECT * FROM playlists
WHERE owner_id = $1
ORDER BY name ASC;

-- name: GetPlaylistByID :one
SELECT * FROM playlists
WHERE id = $1;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE id = $1;

-- name: GetPlaylistsOwnerID :one
SELECT owner_id FROM playlists
WHERE id = $1;