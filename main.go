package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serverMux := http.NewServeMux()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))
	serverMux.HandleFunc("/healthz", healthzHandler)
	serverMux.HandleFunc("/metrics", apiCfg.metricsHandler)
	serverMux.HandleFunc("/reset", apiCfg.resetHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Sett Content-Type header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// TODO: Hent verdien fra atomic.Int32
	hits := cfg.fileserverHits.Load()

	// TODO: Skriv response i formatet "Hits: x\n"
	fmt.Fprintf(w, "Hits: %d\n", hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Sett fileserverHits tilbake til 0
	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
