package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Hien-Trinh/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	apiCfg.db = db

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiCfg.jwtSecret = os.Getenv("JWT_SECRET")

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

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsPost)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerChirpsGetById)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersPost)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersPut)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginPost)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
