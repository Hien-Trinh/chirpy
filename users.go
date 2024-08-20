package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Hien-Trinh/chirpy/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (a *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	users, err := a.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get users: %s", err))
		return
	}

	for _, user := range users {
		if user.Email == params.Email {
			respondWithError(w, http.StatusBadRequest, "User already exists")
			return
		}
	}

	hashed_password := passwordHash(params.Password)

	user, err := a.db.CreateUser(params.Email, hashed_password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	user_without_password := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Email: user.Email,
	}

	respondWithJSON(w, 201, user_without_password)
}

func (a *apiConfig) handlerUsersPut(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	user, err := auth.GetUserByJWT(a.db, a.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get user: %s", err))
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	hashed_password := passwordHash(params.Password)

	user_updated, err := a.db.UpdateUser(user.Id, params.Email, hashed_password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't update user: %s", err))
		return
	}

	user_without_password := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user_updated.Id,
		Email: user_updated.Email,
	}

	respondWithJSON(w, 200, user_without_password)
}

func passwordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash)
}
