package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func profanityFilter(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		lowerCase := strings.ToLower(word)
		if _, ok := badWords[lowerCase]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func handlerValidateChirp(resWriter http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		CleanedBody string `json:"cleaned_body"`
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

	cleanedMessage := profanityFilter(params.Body)
	resVal := returnValue{CleanedBody: cleanedMessage}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}