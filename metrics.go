package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resWriter, req)
	})
}

func (cfg *apiConfig) handlerMetrics(resWriter http.ResponseWriter, req *http.Request) {
	htmlRes := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	resWriter.WriteHeader(200)
	resWriter.Write([]byte(htmlRes))
}