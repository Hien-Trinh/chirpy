package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (a *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := a.db.CreateChirp(getCleanedBody(params.Body))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create chirp: %s", err))
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (a *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get chirp: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (a *apiConfig) handlerChirpsGetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	chirp, err := a.db.GetChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get chirp: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)

}

func getCleanedBody(body string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		lower_word := strings.ToLower(word)
		if _, ok := profaneWords[lower_word]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
