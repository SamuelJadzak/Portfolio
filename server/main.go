package main

import (
	"database/sql"
	"example/data-access/env"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"example/data-access/auth"
	"example/data-access/routes"
	"fmt"
)

func NewServer(postStore *routes.PostStore, authStore *auth.AuthStore) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(
		mux,
		postStore,
		authStore,
	)
	var handler http.Handler = mux
	handler = authMiddleWare(handler)
	handler = timeoutMiddleWare(handler)
	return handler
}

func timeoutMiddleWare(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 2*time.Second, "timeout")
}

func authMiddleWare(next http.Handler) http.Handler {
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
			_, claims, err := auth.ValidateToken(token.Value, "refresh")
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
			_, claims, err := auth.ValidateToken(token, "access")
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

func main() {
	env.Load()
	connStr := buildConnectionStr()
	fmt.Printf("Connecting with: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	postStore := &routes.PostStore{Db: db, Location: env.PostLocation.GetValue()}

	authStore := &auth.AuthStore{Db: db, Location: env.AuthLocation.GetValue()}
	auth.InitAdminPwd(authStore)

	srv := NewServer(postStore, authStore)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(env.ServerHost.GetValue(), env.ServerPort.GetValue()),
		Handler: srv,
	}
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}

func buildConnectionStr() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=disable",
		env.PostgresHost.GetValue(),
		env.PostgresPort.GetValue(),
		env.PostgresDatabase.GetValue(),
		env.PostgresUser.GetValue())
}
