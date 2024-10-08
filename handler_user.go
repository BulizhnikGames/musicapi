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

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Is_artist bool   `json:"is_artist"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Is_artist {
		_, err = apiCfg.DB.GetArtistByName(r.Context(), params.Name)
		if err == nil {
			responseWithError(w, 400, fmt.Sprintf("Artist %s already exists", params.Name))
			return
		}
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
		IsArtist:  params.Is_artist,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	responseWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, 200, user)
}

func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse user id: %v", err))
		return
	}
	if userID != user.ID {
		responseWithError(w, 403, fmt.Sprintf("User %s is not authorized to delete", user.Name))
		return
	}

	err = apiCfg.DB.DeleteUser(r.Context(), userID)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't delete user: %v", err))
		return
	}

	responseWithJSON(w, 200, struct{}{})
}
