package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	req := reqParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json", err)
		return
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.db.UpgradeUser(context.Background(), userID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
