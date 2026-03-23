package middleware

import (
	"example/data-access/auth"
	"fmt"
	"net/http"
	"time"
)

func TimeoutMiddleWare(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 2*time.Second, "timeout")
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			next.ServeHTTP(w, r)
			return
		case "/api/refresh":
			token, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			claims, err := auth.ValidateToken(token.Value, "refresh")

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if err := checkClaims(claims, "refresh"); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		default:
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			claims, err := auth.ValidateToken(token, "access")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if err := checkClaims(claims, "access"); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
	})
}

func checkClaims(claims *auth.CustomClaims, tokenType string) error {
	if claims.Type != tokenType {
		return fmt.Errorf("invalid token type")
	}
	if claims.Username == "" {
		return fmt.Errorf("invalid username")
	}
	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}
	return nil
}
