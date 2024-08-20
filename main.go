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

	v1Router.Post("/albums", apiCfg.middlewareAuth(apiCfg.handlerCreateAlbum))
	v1Router.Get("/albums", apiCfg.handlerGetAlbums) //Should also return all additional data

	v1Router.Post("/songs", apiCfg.middlewareAuth(apiCfg.handlerCreateSong))
	v1Router.Get("/songs", apiCfg.handlerGetSongs) //Should also return all additional data

	v1Router.Get("/albums/{artistName}", apiCfg.handlerGetArtistsAlbums)
	v1Router.Get("/songs/{artistName}", apiCfg.handlerGetArtistsSongs)

	v1Router.Post("/likes/{songID}", apiCfg.middlewareAuth(apiCfg.handlerLikeSongByID))
	v1Router.Post("/likes", apiCfg.middlewareAuth(apiCfg.handlerLikeSong))
	v1Router.Get("/likes", apiCfg.middlewareAuth(apiCfg.handlerGetUsersLikes))

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
