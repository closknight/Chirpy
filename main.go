package main

import (
	"log"
	"net/http"

	"github.com/closknight/Chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const port = "8080"
	database_path := "database.json"
	db, err := database.NewDB(database_path)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := &apiConfig{
		fileServerHits: 0,
		DB:             db,
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
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
