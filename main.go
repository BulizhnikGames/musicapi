package main

import (
	"database/sql"
	"github.com/BulizhnikGames/musicapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT must be set in .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set in .env file")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/r", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Delete("/users/{userID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteUser))

	v1Router.Post("/albums", apiCfg.middlewareAuth(apiCfg.handlerCreateAlbum))
	v1Router.Get("/albums", apiCfg.handlerGetAlbums)
	v1Router.Delete("/albums/{albumID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteAlbum))

	v1Router.Post("/songs", apiCfg.middlewareAuth(apiCfg.handlerCreateSong))
	v1Router.Get("/songs", apiCfg.handlerGetSongs)
	v1Router.Delete("/songs/{songID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteSong))

	v1Router.Post("/likes", apiCfg.middlewareAuth(apiCfg.handlerLikeSong))
	v1Router.Get("/likes", apiCfg.middlewareAuth(apiCfg.handlerGetUsersLikes))
	v1Router.Delete("/likes/{songID}", apiCfg.middlewareAuth(apiCfg.handlerUnlikeSong))

	v1Router.Post("/users/playlists", apiCfg.middlewareAuth(apiCfg.handlerCreatePlaylist))
	v1Router.Get("/users/playlists", apiCfg.middlewareAuth(apiCfg.handlerGetUsersPlaylists))
	v1Router.Delete("/users/playlists/{playlistID}", apiCfg.middlewareAuth(apiCfg.handlerDeletePlaylist))
	v1Router.Post("/users/playlists/songs", apiCfg.middlewareAuth(apiCfg.handlerAddSongToPlaylist))

	v1Router.Get("/playlists/{playlistID}", apiCfg.handlerGetSongsInPlaylist)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server listening on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server crashed with error: ", err)
	}
}
