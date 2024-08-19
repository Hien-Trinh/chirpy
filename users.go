package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	token_parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't parse token: %s", err))
		return
	}

	claims, ok := token_parsed.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse claims")
		return
	}

	token_expiration_time, err := claims.GetExpirationTime()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get expiration time: %s", err))
		return
	}

	if token_expiration_time.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired")
		return
	}

	subject, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get subject: %s", err))
		return
	}

	user_id, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse user ID")
		return
	}

	user, err := a.db.GetUserById(user_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get user: %s", err))
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
