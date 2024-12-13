package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(resWriter http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		Valid bool `json:"valid"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding request data", err)
		return
	}

	const maxBodyLength = 140
	if len(params.Body) > maxBodyLength {
		respondWithError(resWriter, http.StatusBadRequest, "content body is too long", nil)
		return
	}

	resVal := returnValue{Valid: true}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}