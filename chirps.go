package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ansht2000/atServer/internal/database"
	"github.com/google/uuid"
)

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

	chirpParams := database.CreateChirpParams{
		Body: params.Body,
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