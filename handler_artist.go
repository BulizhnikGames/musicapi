package main

import (
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"net/http"
)

func (apiCfg *apiConfig) handlerGetArtistsAlbums(w http.ResponseWriter, r *http.Request, albumName, artistName string) {
	artist, err := apiCfg.DB.GetArtistByName(r.Context(), artistName)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find artist: %v", err))
		return
	}

	if albumName == "" {
		dbAlbums, err := apiCfg.DB.GetArtistsAlbums(r.Context(), artist.ID)
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find %s's albums: %v", artist.Name, err))
			return
		}

		albums, err := apiCfg.databaseAlbumsToAlbums(r.Context(), dbAlbums)
		if err != nil {
			responseWithError(w, 500, fmt.Sprintf("Error creating album's custom JSON: %v", err))
			return
		}

		responseWithJSON(w, 200, albums)
	} else {
		dbAlbum, err := apiCfg.DB.GetAlbumByNameAndArtist(r.Context(), database.GetAlbumByNameAndArtistParams{
			Name:     albumName,
			ArtistID: artist.ID,
		})
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find %s's album with name %s: %v", artist.Name, albumName, err))
			return
		}

		album, err := apiCfg.databaseAlbumToAlbum(r.Context(), dbAlbum)
		if err != nil {
			responseWithError(w, 500, fmt.Sprintf("Error creating album's custom JSON: %v", err))
			return
		}

		responseWithJSON(w, 200, album)
	}
}

func (apiCfg *apiConfig) handlerGetArtistsSongs(w http.ResponseWriter, r *http.Request, songName, artistName string) {
	artist, err := apiCfg.DB.GetArtistByName(r.Context(), artistName)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find artist: %v", err))
		return
	}

	var dbSongs []database.Song
	if songName == "" {
		dbSongs, err = apiCfg.DB.GetSongsByArtist(r.Context(), artist.ID)
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find %s's songs: %v", artist.Name, err))
			return
		}
	} else {
		dbSongs, err = apiCfg.DB.GetSongsByNameAndArtist(r.Context(), database.GetSongsByNameAndArtistParams{
			Name:     songName,
			ArtistID: artist.ID,
		})
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find %s's songs with name %s: %v", artist.Name, songName, err))
			return
		}
	}

	songs, err := apiCfg.databaseSongsToSongs(r.Context(), dbSongs)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating song's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 200, songs)
}
