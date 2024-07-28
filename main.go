package main

import (
	"log"
	"net/http"
)

const MAX_CHIRP_LENGTH = 140

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{fileServerHits: 0}
	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fs))
	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/validate_chirp", HandleValidateChirp)
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
