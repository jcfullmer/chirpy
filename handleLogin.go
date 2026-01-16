package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jcfullmer/chirpy/internal/auth"
	"github.com/jcfullmer/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string        `json:"email"`
		Password         string        `json:"password"`
		ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{
		Email:    "",
		Password: "",
	}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	u, err := cfg.db.LoginUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	loginBool, err := auth.CheckPasswordHash(params.Password, u.HashedPassword)
	switch loginBool {
	case true:
		token, err := auth.MakeJWT(u.ID, cfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
			return
		}
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
			return
		}
		refreshTokenParams := database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: u.ID,
		}
		refreshTokenDB, err := cfg.db.CreateRefreshToken(context.Background(), refreshTokenParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error adding refresh token to database", err)
			return
		}
		user := User{
			ID:           u.ID,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
			Email:        u.Email,
			Token:        token,
			RefreshToken: refreshTokenDB,
		}
		respondWithJSON(w, http.StatusOK, user)
		log.Printf("Logged in user %s", user.Email)
	default:
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
	}
}
