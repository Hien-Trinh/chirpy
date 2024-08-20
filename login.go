package main

import (
	"crypto/rand"
	"encoding/hex"
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
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)) == nil {
				token_expiry := time.Hour
				refresh_token_expiry := time.Hour * 24 * 60

				claims := jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(token_expiry).UTC()),
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

				refresh_token_int, err := rand.Read(make([]byte, 32))
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create refresh token: %s", err))
					return
				}
				refresh_token_string := hex.EncodeToString([]byte(strconv.Itoa(refresh_token_int)))
				_, err = a.db.CreateRefreshToken(refresh_token_string, time.Now().Add(refresh_token_expiry).UTC())
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create refresh token: %s", err))
					return
				}

				user_without_password := struct {
					Id           int    `json:"id"`
					Email        string `json:"email"`
					Token        string `json:"token"`
					RefreshToken string `json:"refresh_token"`
				}{
					Id:           user.Id,
					Email:        user.Email,
					Token:        token_signed,
					RefreshToken: refresh_token_string,
				}

				respondWithJSON(w, 200, user_without_password)
				return
			}
		}
	}

	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
}
