package main

import (
	"database/sql"
	"example/data-access/env"

	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	env.Load()

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=disable",
		env.PostgresHost.GetValue(),
		env.PostgresPort.GetValue(),
		env.PostgresDatabase.GetValue(),
		env.PostgresUser.GetValue())

	fmt.Printf("Connecting with: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
