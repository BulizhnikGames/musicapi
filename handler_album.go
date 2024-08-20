package main

import (
	"encoding/json"
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateAlbum(w http.ResponseWriter, r *http.Request, user database.User) {
	if !user.IsArtist {
		responseWithError(w, 403, "You need to be an artist to create album")
		return
	}

	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	dbAlbum, err := apiCfg.DB.CreateAlbum(r.Context(), database.CreateAlbumParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		ArtistID:  user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error creating album: %v", err))
		return
	}

	album, err := apiCfg.databaseAlbumToAlbum(r.Context(), dbAlbum)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating album's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 201, album)
}

func (apiCfg *apiConfig) handlerGetAlbums(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string `json:"name"`
		Artist string `json:"artist"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Artist != "" {
		apiCfg.handlerGetArtistsAlbums(w, r, params.Name, params.Artist)
		return
	}

	var dbAlbums []database.Album
	if params.Name != "" {
		dbAlbums, err = apiCfg.DB.GetAlbumsByName(r.Context(), params.Name)
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Error getting albums with name (%s): %v", params.Name, err))
			return
		}
	} else {
		dbAlbums, err = apiCfg.DB.GetAllAlbums(r.Context())
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Error getting all albums: %v", err))
			return
		}
	}

	albums, err := apiCfg.databaseAlbumsToAlbums(r.Context(), dbAlbums)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating album's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 200, albums)
}

func (apiCfg *apiConfig) handlerDeleteAlbum(w http.ResponseWriter, r *http.Request, user database.User) {
	albumIDStr := chi.URLParam(r, "albumID")
	albumID, err := uuid.Parse(albumIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse album id: %v", err))
		return
	}

	album, err := apiCfg.DB.GetAlbumByID(r.Context(), albumID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Error getting album with id %v: %v", albumID, err))
		return
	}
	if album.ArtistID != user.ID {
		responseWithError(w, 403, "You need to be an artist of album to delete it")
		return
	}

	err = apiCfg.DB.DeleteAlbum(r.Context(), albumID)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error deleting album: %v", err))
	}

	responseWithJSON(w, 200, struct{}{})
}
