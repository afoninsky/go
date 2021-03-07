package mux

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware returns JWT authentication middleware (HS256)
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) < 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := strings.Join(parts[1:], " ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
