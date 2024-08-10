package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	db, err := NewDB("database.json")
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /admin/metrics/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect to /admin/metrics if slash is present
		http.Redirect(w, r, "/admin/metrics", http.StatusMovedPermanently)
	})

	mux.HandleFunc("POST /api/chirps", db.handlerChirpsPost)
	mux.HandleFunc("GET /api/chirps", db.handlerChirpsGet)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
