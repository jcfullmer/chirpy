package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jcfullmer/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON request", err)
		return
	}
	validatedBody, err := validate_chirp(params.Body)
	if err == fmt.Errorf("Chirp is too long") {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt process chirp", err)
		return
	}
	dbEntry := database.CreateChirpParams{
		Body:   validatedBody,
		UserID: params.User_id,
	}
	chirpDB, err := cfg.db.CreateChirp(context.Background(), dbEntry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creatign user in databs", err)
		return
	}
	c := Chirp{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body:      chirpDB.Body,
		User_id:   chirpDB.UserID,
	}
	respondWithJSON(w, http.StatusCreated, c)
	log.Printf("New chirp created: %v", c.ID)

}

func validate_chirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", fmt.Errorf("Chirp is too long")
	}
	profanityCheck := profanityChecker(body)
	return profanityCheck, nil
}

func profanityChecker(msg string) string {
	if len(msg) == 0 {
		return ""
	}
	words := strings.Split(msg, " ")

	result := []string{}
	for _, word := range words {
		bwCheck := badWordCheck(word)
		result = append(result, bwCheck)
	}
	return strings.Join(result, " ")
}

func badWordCheck(word string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	for _, badWord := range badWords {
		if strings.ToLower(word) == badWord {
			return "****"
		}

	}
	return word
}
