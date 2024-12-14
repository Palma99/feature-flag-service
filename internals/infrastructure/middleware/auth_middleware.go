package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
)

func CheckAuthMiddleware(jwtService *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get token from auth header
			bearerToken := r.Header.Get("Authorization")
			parts := strings.Split(bearerToken, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token := parts[1]

			if payload, err := jwtService.ValidateToken(token); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				ctx := context.WithValue(r.Context(), context_keys.UserIDKey, payload.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func CheckPublicKeyAuthMiddleware(keyService *services.KeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			publicKey := r.Header.Get("X-Public-Key")

			if publicKey == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("public key not found"))
				return
			}

			if !keyService.IsPublicKey(publicKey) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("public key is not valid"))
				return
			}

			ctx := context.WithValue(r.Context(), context_keys.PublicKeyKey, publicKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
