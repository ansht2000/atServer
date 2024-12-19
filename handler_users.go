package main

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/ansht2000/atServer/internal/auth"
	"github.com/ansht2000/atServer/internal/database"
	"github.com/google/uuid"
)

type parametersUsers struct {
	Password string `json:"password"`
	Email string `json:"email"`
	ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
}
type returnValueUsers struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
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

	expiryTime := params.ExpiresInSeconds
	// check if ExpiresInSeconds is unset or over an hour
	if expiryTime == nil {
		// set to the default 1 hour
		var oneHour = 3600
		expiryTime = &oneHour
	}
	*expiryTime = int(math.Min(3600, float64(*expiryTime)))

	tokString, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(*expiryTime) * time.Second)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error making authorization token", err)
		return
	}

	resVals := returnValueUsers{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tokString,
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
	}
	respondWithJSON(resWriter, http.StatusCreated, resVal)
}