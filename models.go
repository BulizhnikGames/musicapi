package main

import (
	"context"
	"github.com/BulizhnikGames/musicapi/internal/database"
)

type Song struct {
	Name    string   `json:"Name"`
	Likes   int64    `json:"Likes"`
	Album   string   `json:"Album"`
	Artists []string `json:"Artists"`
}

type Album struct {
	Name   string `json:"Name"`
	Artist string `json:"Artist"`
	Songs  []Song `json:"Songs"`
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
		Name:   album.Name,
		Artist: artist,
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
