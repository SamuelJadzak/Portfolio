package main

import (
	"database/sql"
	"example/data-access/env"
	"log"
	"net"
	"net/http"
	"os"

	"example/data-access/routes"
	"fmt"
)

// TODO: need to also investigate passing context
// TODO: auth
func NewServer(postStore *routes.PostStore) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(
		mux,
		postStore,
	)
	var handler http.Handler = mux
	// handler = someMiddleware(handler)
	return handler
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
	postStore := &routes.PostStore{Db: db}

	srv := NewServer(postStore)
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

// func someMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		next.ServeHTTP(w, r)
// 	})
// }
