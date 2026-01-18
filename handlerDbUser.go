package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jcfullmer/chirpy/internal/auth"
	"github.com/jcfullmer/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	CreateUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	}
	newUser, err := cfg.db.CreateUser(context.Background(), CreateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a user", err)
		return
	}
	u := User{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: false,
	}
	respondWithJSON(w, http.StatusCreated, u)
	log.Printf("New User created with email: %s", u.Email)
}

func (cfg *apiConfig) handlerUpdateLogin(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "you must input an email and password", fmt.Errorf("empty email or password"))
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	updatedUser, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.CreatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	})
}
