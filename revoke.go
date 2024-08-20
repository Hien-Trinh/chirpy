package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (a *apiConfig) handlerRevokePost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	refresh_tokens, err := a.db.GetRefreshTokens()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get refresh tokens: %s", err))
		return
	}

	for _, refresh_token := range refresh_tokens {
		if refresh_token.Token == token {
			err := a.db.RevokeRefreshToken(refresh_token.Id)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't revoke token: %s", err))
				return
			}

			respondWithJSON(w, http.StatusNoContent, struct{}{})
			return
		}
	}

	respondWithError(w, http.StatusUnauthorized, "Invalid token")
}
