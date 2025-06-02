package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		header := strings.Fields(authHeader)
		if len(header) != 2 || header[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := header[1]
		jwksKeySet, err := jwk.Fetch(r.Context(), os.Getenv("JWKS_URL"))
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		_, err = jwt.Parse([]byte(token), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
