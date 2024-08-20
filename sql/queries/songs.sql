-- name: CreateSong :one
INSERT INTO songs (id, created_at, updated_at, name, album_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAlbumsSongs :many
SELECT * FROM songs
WHERE album_id = $1
ORDER BY updated_at ASC;

-- name: GetSongsByName :many
SELECT * FROM songs
WHERE name = $1
ORDER BY updated_at DESC;

-- name: DeleteSongByID :exec
DELETE FROM songs
WHERE id = $1;

-- name: GetAllSongs :many
SELECT * FROM songs
ORDER BY name ASC;

-- name: GetSongsByNameAndArtist :many
SELECT songs.* FROM songs
JOIN artists_songs ON songs.id = artists_songs.song_id
WHERE songs.name = $1 AND artists_songs.artist_id = $2
ORDER BY songs.updated_at DESC;