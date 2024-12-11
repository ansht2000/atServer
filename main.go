package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resWriter, req)
	})
}

func (cfg *apiConfig) handlerMetrics(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	resWriter.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("/healthz", handlerReadiness)
	serveMux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("/reset", apiCfg.handlerReset)
	server := http.Server{Handler: serveMux, Addr: ":" + port}
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}