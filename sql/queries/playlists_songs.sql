-- name: AddSongToPlaylist :one
INSERT INTO playlists_songs (playlist_id, song_id, add_time)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSongsInPlaylist :many
SELECT songs.* FROM songs
JOIN playlists_songs ON playlists_songs.song_id = songs.id
WHERE playlists_songs.playlist_id = $1
ORDER BY add_time ASC;

-- name: RemoveSongFromPlaylist :exec
DELETE FROM playlists_songs
WHERE playlist_id = $1 AND song_id = $2;