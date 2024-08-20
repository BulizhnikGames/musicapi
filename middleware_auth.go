package main

import (
	"fmt"
	"github.com/BulizhnikGames/musicapi/internal/auth"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := auth.GetEmail(r.Header)
		if err != nil {
			responseWithError(w, 401, fmt.Sprintf("auth error: %v", err))
			return
		}
		password, err := auth.GetPassword(r.Header)
		if err != nil {
			responseWithError(w, 401, fmt.Sprintf("auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByEmailAndPassword(r.Context(), database.GetUserByEmailAndPasswordParams{
			Email:    email,
			Password: password,
		})
		if err != nil {
			responseWithError(w, 401, fmt.Sprintf("Couldn't get user (auth error): %v", err))
			return
		}

		handler(w, r, user)
	}
}
