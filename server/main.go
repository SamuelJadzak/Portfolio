package main

import (
	"database/sql"
	"example/data-access/env"
	"log"
	"net"
	"net/http"
	"os"

	"example/data-access/auth"
	"example/data-access/middleware"
	"example/data-access/routes"
	"fmt"

	"github.com/rs/cors"
)

func NewServer(postStore *routes.PostStore, authStore *auth.AuthStore) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(
		mux,
		postStore,
		authStore,
	)
	var handler http.Handler = mux
	handler = middleware.AuthMiddleWare(handler)
	handler = middleware.TimeoutMiddleWare(handler)
	// handler = cors.Default().Handler(handler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"}, // Your frontend domain
		// Set Access-Control-Allow-Methods here
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(handler)
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
