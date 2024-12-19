package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/ansht2000/atServer/internal/auth"
)

var ErrRefreshTokenExpired = errors.New("refresh token has expired")

type returnValueRefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(resWriter http.ResponseWriter, req *http.Request) {
	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "could not get refresh token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), userToken)
	if err == sql.ErrNoRows {
		respondWithError(resWriter, http.StatusUnauthorized, "refresh token does not exist", err)
		return
	} else if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error encountered while getting refresh token", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid  {
		respondWithError(resWriter, http.StatusUnauthorized, "refresh token has expired", ErrRefreshTokenExpired)
		return
	}

	newJWT, err := auth.MakeJWT(refreshToken.UserID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error making authorization token", err)
		return
	}

	resVal := returnValueRefreshToken{
		Token: newJWT,
	}
	respondWithJSON(resWriter, http.StatusOK, resVal)
}

func (cfg *apiConfig) handlerRevoke(resWriter http.ResponseWriter, req *http.Request) {
	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "could not get refresh token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), userToken)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "could not revoke refresh token", err)
	}
	respondWithJSON(resWriter, http.StatusNoContent, nil)
}