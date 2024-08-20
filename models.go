package main

import (
	"context"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/google/uuid"
)

type Song struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"Name"`
	Likes   int64     `json:"Likes"`
	Album   string    `json:"Album"`
	Artists []string  `json:"Artists"`
}

type Album struct {
	ID     uuid.UUID `json:"ID"`
	Name   string    `json:"Name"`
	Artist string    `json:"Artist"`
	Songs  []Song    `json:"Songs"`
}

type Playlist struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"Name"`
	OwnerID uuid.UUID `json:"Owner"`
	Songs   []Song    `json:"Songs"`
}

func (apiCfg *apiConfig) databaseSongToSong(ctx context.Context, song database.Song) (Song, error) {
	likes, err := apiCfg.DB.GetSongsLikeCount(ctx, song.ID)
	if err != nil {
		return Song{}, err
	}

	album, err := apiCfg.DB.GetAlbumByID(ctx, song.AlbumID)
	if err != nil {
		return Song{}, err
	}

	artists, err := apiCfg.DB.GetArtistsOfSong(ctx, song.ID)
	if err != nil {
		return Song{}, err
	}

	return Song{
		ID:      song.ID,
		Name:    song.Name,
		Likes:   likes,
		Album:   album.Name,
		Artists: artists,
	}, nil
}

func (apiCfg *apiConfig) databaseSongsToSongs(ctx context.Context, songs []database.Song) ([]Song, error) {
	res := []Song{}
	for _, dbSong := range songs {
		song, err := apiCfg.databaseSongToSong(ctx, dbSong)
		if err != nil {
			return []Song{}, err
		}
		res = append(res, song)
	}
	return res, nil
}

func (apiCfg *apiConfig) databaseAlbumToAlbum(ctx context.Context, album database.Album) (Album, error) {
	artist, err := apiCfg.DB.GetAlbumsArtist(ctx, album.ID)
	if err != nil {
		return Album{}, err
	}

	dbSongs, err := apiCfg.DB.GetAlbumsSongs(ctx, album.ID)
	if err != nil {
		return Album{}, err
	}

	songs, err := apiCfg.databaseSongsToSongs(ctx, dbSongs)
	if err != nil {
		return Album{}, err
	}

	return Album{
		ID:     album.ID,
		Name:   album.Name,
		Artist: artist.Name,
		Songs:  songs,
	}, nil
}

func (apiCfg *apiConfig) databaseAlbumsToAlbums(ctx context.Context, albums []database.Album) ([]Album, error) {
	res := []Album{}
	for _, dbAlbum := range albums {
		album, err := apiCfg.databaseAlbumToAlbum(ctx, dbAlbum)
		if err != nil {
			return []Album{}, err
		}
		res = append(res, album)
	}
	return res, nil
}

func (apiCfg *apiConfig) databasePlaylistToPlaylist(ctx context.Context, playlist database.Playlist) (Playlist, error) {
	dbSongs, err := apiCfg.DB.GetSongsInPlaylist(ctx, playlist.ID)
	if err != nil {
		return Playlist{}, err
	}

	songs, err := apiCfg.databaseSongsToSongs(ctx, dbSongs)
	if err != nil {
		return Playlist{}, err
	}

	return Playlist{
		ID:      playlist.ID,
		Name:    playlist.Name,
		OwnerID: playlist.OwnerID,
		Songs:   songs,
	}, nil
}

func (apiCfg *apiConfig) databasePlaylistsToPlaylists(ctx context.Context, playlists []database.Playlist) ([]Playlist, error) {
	res := []Playlist{}
	for _, dbPlaylist := range playlists {
		playlist, err := apiCfg.databasePlaylistToPlaylist(ctx, dbPlaylist)
		if err != nil {
			return []Playlist{}, err
		}
		res = append(res, playlist)
	}
	return res, nil
}
