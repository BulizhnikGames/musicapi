// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: artists_songs.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createArtistSongLink = `-- name: CreateArtistSongLink :one
INSERT INTO artists_songs (artist_id, song_id)
VALUES ($1, $2)
RETURNING artist_id, song_id
`

type CreateArtistSongLinkParams struct {
	ArtistID uuid.UUID
	SongID   uuid.UUID
}

func (q *Queries) CreateArtistSongLink(ctx context.Context, arg CreateArtistSongLinkParams) (ArtistsSong, error) {
	row := q.db.QueryRowContext(ctx, createArtistSongLink, arg.ArtistID, arg.SongID)
	var i ArtistsSong
	err := row.Scan(&i.ArtistID, &i.SongID)
	return i, err
}

const getArtistsOfSong = `-- name: GetArtistsOfSong :many
SELECT users.name FROM users
JOIN artists_songs ON users.id = artists_songs.artist_id
WHERE artists_songs.song_id = $1
ORDER BY users.name ASC
`

func (q *Queries) GetArtistsOfSong(ctx context.Context, songID uuid.UUID) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getArtistsOfSong, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSongsByArtist = `-- name: GetSongsByArtist :many
SELECT songs.id, songs.created_at, songs.updated_at, songs.name, songs.album_id FROM songs
JOIN artists_songs ON songs.id = artists_songs.song_id
WHERE artists_songs.artist_id = $1
ORDER BY songs.name ASC
`

func (q *Queries) GetSongsByArtist(ctx context.Context, artistID uuid.UUID) ([]Song, error) {
	rows, err := q.db.QueryContext(ctx, getSongsByArtist, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Song
	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.AlbumID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
