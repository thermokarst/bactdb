package handlers

import (
	"net/http"
	"os"
	"strings"
)

func corsHandler(h http.Handler) http.Handler {
	cors := func(w http.ResponseWriter, r *http.Request) {
		domains := os.Getenv("DOMAINS")
		allowedDomains := strings.Split(domains, ",")
		if origin := r.Header.Get("Origin"); origin != "" {
			for _, s := range allowedDomains {
				if s == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				}
			}
		}
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(cors)
}
