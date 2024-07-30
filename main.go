package main

import (
	"log"
	"net/http"

	"github.com/closknight/Chirpy/internal/database"
)

func main() {
	const port = "8080"
	database_path := "database.json"
	db, err := database.NewDB(database_path)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	apiCfg := &apiConfig{fileServerHits: 0}
	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fs))
	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/chirps", db.HandleNewChirp)
	mux.HandleFunc("GET /api/chirps", db.HandleGetChirps)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileServerHits int
}
