package main

import (
	"encoding/json"
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateSong(w http.ResponseWriter, r *http.Request, user database.User) {
	if !user.IsArtist {
		responseWithError(w, 403, "You need to be an artist to create song")
		return
	}

	type parameters struct {
		Name    string   `json:"name"`
		Album   string   `json:"album"`
		Artists []string `json:"artists"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Artists == nil {
		responseWithError(w, 400, "Artists are required")
		return
	}

	album, err := apiCfg.DB.GetAlbumByNameAndArtist(r.Context(), database.GetAlbumByNameAndArtistParams{
		Name:     params.Album,
		ArtistID: user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error getting album (name: %s, artistID: %s) for this (%s) song: %v",
			params.Album, user.ID, params.Name, err))
		return
	}

	artistsIDs := []uuid.UUID{user.ID}
	for _, artist := range params.Artists {
		artistRaw, err := apiCfg.DB.GetArtistByName(r.Context(), artist)
		if err != nil {
			responseWithError(w, 400, "One or more artists could not be found")
			return
		} else {
			artistsIDs = append(artistsIDs, artistRaw.ID)
		}
	}

	dbSong, err := apiCfg.DB.CreateSong(r.Context(), database.CreateSongParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		AlbumID:   album.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error creating song: %v", err))
		return
	}
	for _, artist := range artistsIDs {
		_, err = apiCfg.DB.CreateArtistSongLink(r.Context(), database.CreateArtistSongLinkParams{
			ArtistID: artist,
			SongID:   dbSong.ID,
		})
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("Error creating song: %v", err))
			err = apiCfg.DB.DeleteSongByID(r.Context(), dbSong.ID)
			if err != nil {
				log.Println("!Error deleting song!:", err)
			}
			return
		}
	}

	song, err := apiCfg.databaseSongToSong(r.Context(), dbSong)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Error creating song's custom JSON: %v", err))
		return
	}

	responseWithJSON(w, 201, song)
}

func (apiCfg *apiConfig) handlerGetSongs(w http.ResponseWriter, r *http.Request) {
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

	var dbSongs []database.Song
	if params.Name != "" {
		dbSongs, err = apiCfg.DB.GetSongsByName(r.Context(), params.Name)
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("Error getting songs with name (%s): %v", params.Name, err))
			return
		}
	} else {
		dbSongs, err = apiCfg.DB.GetAllSongs(r.Context())
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("Error getting all songs: %v", err))
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
