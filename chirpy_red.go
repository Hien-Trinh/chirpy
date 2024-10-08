package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (a *apiConfig) handlerChirpyRedPost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if token != a.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = a.db.UpdateUserChirpyRed(params.Data.UserId, true)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't update user: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
