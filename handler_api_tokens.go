package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jcfullmer/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerReresh(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		token string `json:"token"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from auth header", err)
	}
	expires_at, err := cfg.db.RefreshTokenLookup(context.Background(), token)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
	} else if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error checking databse for refresh token", err)
	}
	if expires_at.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "token expired", err)
	}
	respondWithJSON(w, http.StatusOK, token)
}
