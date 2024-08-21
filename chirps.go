package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hien-Trinh/chirpy/internal/auth"
)

func (a *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := auth.GetUserByJWT(a.db, a.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get user: %s", err))
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := a.db.CreateChirp(user.Id, getCleanedBody(params.Body))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create chirp: %s", err))
		return
	}

	respondWithJSON(w, 201, chirp)
}

// handlerChirpsGet returns all chirps
func (a *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get chirp: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

// handlerChirpsGetById returns a chirp by ID
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

// handlerChirpsDeleteById deletes a chirp by ID
func (a *apiConfig) handlerChirpsDeleteById(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := auth.GetUserByJWT(a.db, a.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get user: %s", err))
		return
	}

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

	if chirp.AuthorId != user.Id {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps")
		return
	}

	chirp, err = a.db.DeleteChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't delete chirp: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

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
