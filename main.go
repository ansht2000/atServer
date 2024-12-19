package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ansht2000/atServer/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	secretKey string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not establish connection to db: %v\n", err)
	}
	dbQueries := database.New(db)
	secretKey := os.Getenv("SECRET")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		secretKey: secretKey,
	}

	serveMux := http.NewServeMux()
	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serveMux.Handle("/app/", fileserverHandler)

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpsFromID)

	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)

	server := &http.Server{Handler: serveMux, Addr: ":" + port}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}