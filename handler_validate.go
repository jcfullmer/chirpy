package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	profanityCheck := profanityChecker(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: profanityCheck,
	})

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
