package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (a *apiConfig) handlerLoginPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
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
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)) == nil {
				time_till_expiry := 24 * time.Hour
				if params.Expires_in_seconds > 0 {
					time_till_expiry = time.Duration(params.Expires_in_seconds) * time.Second
				}

				claims := jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time_till_expiry).UTC()),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
					Issuer:    "chirpy",
					Subject:   strconv.Itoa(user.Id),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				token_signed, err := token.SignedString([]byte(a.jwtSecret))
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create token: %s", err))
					return
				}

				user_without_password := struct {
					Id    int    `json:"id"`
					Email string `json:"email"`
					Token string `json:"token"`
				}{
					Id:    user.Id,
					Email: user.Email,
					Token: token_signed,
				}

				respondWithJSON(w, 200, user_without_password)
				return
			}
		}
	}

	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
}
