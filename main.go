package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{fileServerHits: 0}
	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fs))
	mux.HandleFunc("/healthz", HandleHealthz)
	mux.HandleFunc("/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("/reset", apiCfg.HandleReset)
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
