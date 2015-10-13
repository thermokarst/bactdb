package handlers

import "net/http"

func jsonHandler(h http.Handler) http.Handler {
	j := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(j)
}
