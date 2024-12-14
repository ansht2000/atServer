package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(resWriter http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	errRes := errorResponse{Error: msg}
	respondWithJSON(resWriter, code, errRes)
}

func respondWithJSON(resWriter http.ResponseWriter, code int, payload interface{}) {
	resWriter.Header().Set("Content-Type", "application/json")
	resWriter.WriteHeader(code)
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		res = []byte("An unexpected error occurred")
		resWriter.Header().Set("Content-Type", "text/plaintext")
		resWriter.WriteHeader(500)
	}
	resWriter.Write(res)
}