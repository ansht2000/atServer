package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ansht2000/atServer/internal/database"
	"github.com/ansht2000/atServer/internal/auth"
	"github.com/google/uuid"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

var ErrBodyLengthTooLong = errors.New("body length is too long")

type parametersChirps struct {
	Body string `json:"body"`
}
type returnValueChirps struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
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
		return "", ErrBodyLengthTooLong
	}
	cleanedMessage := profanityFilter(body)
	return cleanedMessage, nil
}

func (cfg *apiConfig) handlerCreateChirp(resWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parametersChirps{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding request data", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "error getting authorization header", err)
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, cfg.secretKey)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "error validating user token", err)
		return
	}

	cleanedMessage, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body: cleanedMessage,
		UserID: userID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating chirp", err)
		return
	}

	resVal := returnValueChirps{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(resWriter, http.StatusCreated, resVal)
}

func (cfg *apiConfig) handlerGetChirps(resWriter http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error getting chirps", err)
		return
	}

	resVals := make([]returnValueChirps, len(chirps))
	for i, chirp := range chirps {
		resVals[i] = returnValueChirps{
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

	resVal := returnValueChirps{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}