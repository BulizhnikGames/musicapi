-- name: LikeSong :one
INSERT INTO likes (user_id, song_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUsersLikedSongs :many
SELECT songs.* FROM songs
JOIN likes ON likes.song_id = songs.id
WHERE likes.user_id = $1
ORDER BY songs.name ASC;

-- name: GetSongsLikeCount :one
SELECT COUNT(*) FROM likes
WHERE song_id = $1;

-- name: UnlikeSong :exec
DELETE FROM likes
WHERE user_id = $1 AND song_id = $2;