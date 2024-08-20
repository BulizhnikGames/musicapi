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

func (apiCfg *apiConfig) handlerCreatePlaylist(w http.ResponseWriter, r *http.Request, user database.User) {
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

	dbPlaylist, err := apiCfg.DB.CreatePlayList(r.Context(), database.CreatePlayListParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		OwnerID:   user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error creating playlist: %v", err))
		return
	}

	playlist, err := apiCfg.databasePlaylistToPlaylist(r.Context(), dbPlaylist)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating playlist's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 201, playlist)
}

func (apiCfg *apiConfig) handlerGetUsersPlaylists(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		PlaylistIDStr string `json:"playlist"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.PlaylistIDStr == "" {
		dbPlaylists, err := apiCfg.DB.GetUsersPlaylists(r.Context(), user.ID)
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find %s's playlists: %v", user.Name, err))
			return
		}

		playlists, err := apiCfg.databasePlaylistsToPlaylists(r.Context(), dbPlaylists)
		if err != nil {
			responseWithError(w, 500, fmt.Sprintf("Error creating playlist's custom JSON: %v", err))
			return
		}

		responseWithJSON(w, 200, playlists)
	} else {
		playlistID, err := uuid.Parse(params.PlaylistIDStr)
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("Couldn't parse playlist id: %v", err))
			return
		}

		dbPlaylist, err := apiCfg.DB.GetPlaylistByID(r.Context(), playlistID)
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't find playlist with name id %s: %v", playlistID, err))
			return
		}

		playlist, err := apiCfg.databasePlaylistToPlaylist(r.Context(), dbPlaylist)
		if err != nil {
			responseWithError(w, 500, fmt.Sprintf("Error creating playlist's custom JSON: %v", err))
			return
		}

		responseWithJSON(w, 200, playlist)
	}
}

func (apiCfg *apiConfig) handlerGetSongsInPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistIDStr := chi.URLParam(r, "playlistID")
	playlistID, err := uuid.Parse(playlistIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse playlist id: %v", err))
		return
	}

	dbPlaylist, err := apiCfg.DB.GetPlaylistByID(r.Context(), playlistID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find playlist with id: %v", err))
		return
	}

	playlist, err := apiCfg.databasePlaylistToPlaylist(r.Context(), dbPlaylist)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating playlist's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 200, playlist)
}

func (apiCfg *apiConfig) handlerDeletePlaylist(w http.ResponseWriter, r *http.Request, user database.User) {
	playlistIDStr := chi.URLParam(r, "playlistID")
	playlistID, err := uuid.Parse(playlistIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse playlist id: %v", err))
		return
	}

	type parameters struct {
		SongID uuid.UUID `json:"song_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.SongID != uuid.Nil {
		apiCfg.handlerRemoveSongFromPlaylist(w, r, user, playlistID, params.SongID)
		return
	}

	ownerID, err := apiCfg.DB.GetPlaylistsOwnerID(r.Context(), playlistID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find playlist with id: %v", err))
		return
	}
	if ownerID != user.ID {
		responseWithError(w, 403, fmt.Sprintf("You do not own this playlist: %v", user))
		return
	}

	err = apiCfg.DB.DeletePlaylist(r.Context(), playlistID)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error deleting playlist: %v", err))
	}

	responseWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handlerRemoveSongFromPlaylist(w http.ResponseWriter, r *http.Request,
	user database.User, playlistID uuid.UUID, songID uuid.UUID) {
	ownerID, err := apiCfg.DB.GetPlaylistsOwnerID(r.Context(), playlistID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find playlist with id: %v", err))
		return
	}
	if ownerID != user.ID {
		responseWithError(w, 403, fmt.Sprintf("You do not own this playlist: %v", user))
		return
	}

	err = apiCfg.DB.RemoveSongFromPlaylist(r.Context(), database.RemoveSongFromPlaylistParams{
		PlaylistID: playlistID,
		SongID:     songID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error removing song from play list: %v", err))
	}

	responseWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handlerAddSongToPlaylist(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		SongID     uuid.UUID `json:"song_id"`
		PlaylistID uuid.UUID `json:"playlist_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	playlist, err := apiCfg.DB.GetPlaylistByID(r.Context(), params.PlaylistID)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't find playlist with id %s: %v", params.PlaylistID, err))
		return
	}
	if playlist.OwnerID != user.ID {
		responseWithError(w, 403, fmt.Sprintf("You do not own this playlist"))
		return
	}

	playlist_song_link, err := apiCfg.DB.AddSongToPlaylist(r.Context(), database.AddSongToPlaylistParams{
		PlaylistID: params.PlaylistID,
		SongID:     params.SongID,
		AddTime:    time.Now(),
	})

	responseWithJSON(w, 201, playlist_song_link)
}
