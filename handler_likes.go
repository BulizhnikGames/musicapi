package main

import (
	"encoding/json"
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (apiCfg *apiConfig) handlerLikeSongByID(w http.ResponseWriter, r *http.Request, user database.User) {
	songIDStr := chi.URLParam(r, "songID")
	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	like, err := apiCfg.DB.LikeSong(r.Context(), database.LikeSongParams{
		UserID: user.ID,
		SongID: songID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't like song: %v", err))
		return
	}

	responseWithJSON(w, 201, like)
}

func (apiCfg *apiConfig) handlerLikeSong(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name   string `json:"name"`
		Album  string `json:"album"`
		Artist string `json:"artist"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	artist, err := apiCfg.DB.GetArtistByName(r.Context(), params.Artist)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Artist not found: %v", err))
		return
	}

	songs, err := apiCfg.DB.GetSongsByNameAndArtist(r.Context(), database.GetSongsByNameAndArtistParams{
		Name:     params.Name,
		ArtistID: artist.ID,
	})

	var res *database.Song
	for _, song := range songs {
		album, err := apiCfg.DB.GetAlbumByID(r.Context(), song.AlbumID)
		if err != nil {
			continue
		}
		if album.Name == params.Album {
			res = &song
			break
		}
	}
	if res == nil {
		responseWithError(w, 404, fmt.Sprintf("Song not found"))
		return
	}

	like, err := apiCfg.DB.LikeSong(r.Context(), database.LikeSongParams{
		UserID: user.ID,
		SongID: (*res).ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't like song: %v", err))
		return
	}

	responseWithJSON(w, 201, like)
}

func (apiCfg *apiConfig) handlerGetUsersLikes(w http.ResponseWriter, r *http.Request, user database.User) {
	dbSongs, err := apiCfg.DB.GetUsersLikedSongs(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't get user's likes: %v", err))
		return
	}

	songs, err := apiCfg.databaseSongsToSongs(r.Context(), dbSongs)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating song's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 200, songs)
}
