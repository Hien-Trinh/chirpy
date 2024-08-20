package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *apiConfig) handlerRefreshPost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	refresh_tokens, err := a.db.GetRefreshTokens()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get refresh tokens: %s", err))
		return
	}

	for _, refresh_token := range refresh_tokens {
		if refresh_token.Token == token {
			if refresh_token.ExpiresAt.Before(time.Now().UTC()) {
				respondWithError(w, http.StatusUnauthorized, "Token has expired")
				return
			}

			token_expiry := time.Hour

			claims := jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(token_expiry).UTC()),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				Issuer:    "chirpy",
				Subject:   strconv.Itoa(refresh_token.UserID),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token_signed, err := token.SignedString([]byte(a.jwtSecret))
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create token: %s", err))
				return
			}

			respondWithJSON(w, 200, struct {
				Token string `json:"token"`
			}{Token: token_signed})
			return
		}
	}

	respondWithError(w, http.StatusUnauthorized, "Invalid token")
}
