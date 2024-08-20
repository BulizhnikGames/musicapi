-- name: CreateAlbum :one
INSERT INTO albums (id, created_at, updated_at, name, artist_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetArtistsAlbums :many
SELECT * FROM albums
WHERE artist_id = $1
ORDER BY updated_at DESC;

-- name: GetAlbumsByName :many
SELECT * FROM albums
WHERE name = $1
ORDER BY updated_at DESC;

-- name: GetAlbumByNameAndArtist :one
SELECT * FROM albums
WHERE name = $1 AND artist_id = $2
LIMIT 1;

-- name: GetAlbumByID :one
SELECT * FROM albums
WHERE id = $1
LIMIT 1;

-- name: GetAlbumsArtist :one
SELECT users.id, users.name FROM users
JOIN albums ON albums.artist_id = users.id
WHERE albums.id = $1;

-- name: GetAllAlbums :many
SELECT * FROM albums
ORDER BY name ASC;

-- name: DeleteAlbum :exec
DELETE FROM albums
WHERE id = $1;