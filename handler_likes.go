package main

import (
	"encoding/json"
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (apiCfg *apiConfig) handlerLikeSong(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		ID uuid.UUID `json:"song_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	like, err := apiCfg.DB.LikeSong(r.Context(), database.LikeSongParams{
		UserID: user.ID,
		SongID: params.ID,
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

func (apiCfg *apiConfig) handlerUnlikeSong(w http.ResponseWriter, r *http.Request, user database.User) {
	songIDStr := chi.URLParam(r, "songID")
	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse songID: %v", err))
	}

	err = apiCfg.DB.UnlikeSong(r.Context(), database.UnlikeSongParams{
		UserID: user.ID,
		SongID: songID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't unlike song: %v", err))
	}

	responseWithJSON(w, 200, struct{}{})
}
