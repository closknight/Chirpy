package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/closknight/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	const port = "8080"
	database_path := "database.json"

	if *dbg {
		err := os.Remove(database_path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}

	db, err := database.NewDB(database_path)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := &apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}
	mux := http.NewServeMux()
	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fs))
	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChipByID)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandleUpdateUser)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
