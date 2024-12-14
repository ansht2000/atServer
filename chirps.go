package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ansht2000/atServer/internal/database"
	"github.com/google/uuid"
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

func validateChirp(body string) (string, error) {
	const maxBodyLength = 140
	if len(body) > maxBodyLength {
		return "", errors.New("body length is too long")
	}
	cleanedMessage := profanityFilter(body)
	return cleanedMessage, nil
}

func (cfg *apiConfig) handlerCreateChirp(resWriter http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnValue struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding request data", err)
	}

	cleanedMessage, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, err.Error(), err)
	}

	chirpParams := database.CreateChirpParams{
		Body: cleanedMessage,
		UserID: params.UserID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating chirp", err)
		return
	}

	resVal := returnValue{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(resWriter, http.StatusCreated, resVal)
}

func (cfg *apiConfig) handlerGetChirps(resWriter http.ResponseWriter, req *http.Request) {
	type returnValue struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error getting chirps", err)
		return
	}

	resVals := make([]returnValue, len(chirps))
	for i, chirp := range chirps {
		resVals[i] = returnValue{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
	}
	respondWithJSON(resWriter, http.StatusOK, resVals)
}

func (cfg *apiConfig) handlerGetChirpsFromID(resWriter http.ResponseWriter, req *http.Request) {
	type returnValue struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error parsing id from request", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err == sql.ErrNoRows {
		respondWithError(resWriter, http.StatusNotFound, "chirp not found", err)
		return
	} else if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error getting chirp", err)
		return
	}

	resVal := returnValue{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}