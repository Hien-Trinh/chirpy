package main

import (
	"net/http"
	"strings"

	"github.com/Hien-Trinh/chirpy/internal/auth"
)

func (a *apiConfig) handlerRefreshPost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	refresh_token, err := a.db.GetRefreshTokensByToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token")
		return
	}

	token_signed, err := auth.CreateJWT(a.jwtSecret, refresh_token.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token_signed,
	})
}
