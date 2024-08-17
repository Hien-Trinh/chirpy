package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := a.db.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	respondWithJSON(w, 201, user)
}

// func (a *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request) {
// 	chirps, err := a.db.GetChirps()
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get chirp: %s", err))
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, chirps)

// }

// func (a *apiConfig) handlerUsersGetById(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.PathValue("id"))
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID: %s", err))
// 		return
// 	}

// 	chirp, err := a.db.GetChirpById(id)
// 	if err != nil {
// 		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get chirp: %s", err))
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, chirp)

// }
