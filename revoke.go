package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (a *apiConfig) handlerRevokePost(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	refresh_token, err := a.db.GetRefreshTokensByToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token")
		return
	}

	err = a.db.RevokeRefreshToken(refresh_token.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't revoke token: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
