package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ansht2000/atServer/internal/auth"
	"github.com/ansht2000/atServer/internal/database"
	"github.com/google/uuid"
)

var ErrInvalidAPIKey = errors.New("invalid api key")

type parametersUsers struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type parametersWebhook struct {
	Event string `json:"event"`
	Data struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

type returnValueUsers struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerLoginUser(resWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parametersUsers{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding login request data", err)
		return
	}

	user, err := cfg.db.GetUserFromEmail(req.Context(), params.Email)
	if err == sql.ErrNoRows {
		respondWithError(resWriter, http.StatusUnauthorized, "incorrect email or password", err)
		return
	} else if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error retrieving user data", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	tokString, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error making authorization token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error making refresh token", err)
		return
	}

	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	resVals := returnValueUsers{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tokString,
		RefreshToken: refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(resWriter, http.StatusOK, resVals)
}

func (cfg *apiConfig) handlerCreateUser(resWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parametersUsers{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding creating user request data", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating password hash", err)
		return
	}

	createUserParams := database.CreateUserParams{
		HashedPassword: hashedPass,
		Email: params.Email,
	}

	user, err := cfg.db.CreateUser(req.Context(), createUserParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating user", err)
		return
	}

	resVal := returnValueUsers{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(resWriter, http.StatusCreated, resVal)
}

func (cfg *apiConfig) handlerUpdateUser(resWriter http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "error getting authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secretKey)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "error validating access token", err)
		return
	}

	defer req.Body.Close()
	params := parametersUsers{}
	decoder := json.NewDecoder(req.Body)
	if err = decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding request data", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error hashing new password", err)
		return
	}

	updateParams := database.UpdateUserEmailPasswordByIDParams{
		HashedPassword: hashedPass,
		Email: params.Email,
		ID: userID,
	}
	user, err := cfg.db.UpdateUserEmailPasswordByID(req.Context(), updateParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error updating email and password", err)
	}

	resVal := returnValueUsers{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}

func (cfg *apiConfig) handlerUpgradeUser(resWriter http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "could not find api key", err)
		return
	}

	if apiKey != cfg.apiKey {
		respondWithError(resWriter, http.StatusUnauthorized, "invalid api key", ErrInvalidAPIKey)
		return
	}

	defer req.Body.Close()
	params := parametersWebhook{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding request data", err)
		return
	}

	if params.Event != "user.upgraded" {
		resWriter.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserByID(req.Context(), params.Data.UserID)
	if err == sql.ErrNoRows {
		respondWithError(resWriter, http.StatusNotFound, "user not found", err)
		return
	} else if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error retrieving user data", err)
		return
	}

	resWriter.WriteHeader(http.StatusNoContent)
}
