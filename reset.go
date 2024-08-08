package main

import "net/http"

func (a *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset"))
}
