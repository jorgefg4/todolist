package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jorgefg4/todolist/pkg/database"

	"github.com/jorgefg4/todolist/pkg/service"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// parseURI build the connection string to the postgres database
func parseURI() string {
	return "postgresql://" + os.Getenv("USER_DB") + ":" + os.Getenv("PASSWORD_DB") +
		"@" + os.Getenv("HOST_DB") + ":" + os.Getenv("PORT_DB") + "/" +
		os.Getenv("NAME_DB") + "?sslmode=disable"
}

func main() {
	// Declaration of required variables
	var db *sql.DB
	var ctx context.Context = context.Background()

	// Call to the Service layer to create the Server
	ph := database.NewPostgres(db, ctx)
	svc := service.NewService(ph)
	s, err := svc.NewServer(parseURI())
	if err != nil {
		log.Fatal(err)
	}

	// CORS headers:
	handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}}).Handler(s.Router())

	// It listens on TCP port 8000 of localhost and calls the handler
	log.Fatal(http.ListenAndServe(":8000", handler))

	// If nothing fails, it indicates that the server is up and running
	fmt.Printf("Server running\n")
}
