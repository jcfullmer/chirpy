package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jcfullmer/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from auth header", err)
		return
	}
	tokenDB, err := cfg.db.RefreshTokenLookup(context.Background(), token)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error checking database for refresh token", err)
		return
	}
	if tokenDB.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "token expired", err)
		return
	}
	if tokenDB.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	newToken, err := auth.MakeJWT(tokenDB.UserID, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new token", err)
		return
	}
	type tokenStruct struct {
		Token string `json:"token"`
	}
	newtokenS := tokenStruct{
		Token: newToken,
	}
	respondWithJSON(w, http.StatusOK, newtokenS)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from auth header", err)
	}
	err = cfg.db.RevokeToken(context.Background(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking token", err)
	}
	w.WriteHeader(http.StatusNoContent)

}
