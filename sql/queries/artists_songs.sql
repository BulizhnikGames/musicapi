-- name: CreateArtistSongLink :one
INSERT INTO artists_songs (artist_id, song_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetSongsByArtist :many
SELECT songs.* FROM songs
JOIN artists_songs ON songs.id = artists_songs.song_id
WHERE artists_songs.artist_id = $1
ORDER BY songs.name ASC;

-- name: GetArtistsOfSong :many
SELECT users.name FROM users
JOIN artists_songs ON users.id = artists_songs.artist_id
WHERE artists_songs.song_id = $1
ORDER BY users.name ASC;