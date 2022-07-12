package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/jorgefg4/todolist/pkg/database"

	"github.com/jorgefg4/todolist/pkg/service"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	// Declaramos de variables necesarias
	var db *sql.DB
	var ctx context.Context = context.Background()

	//Llamamos a la capa Service para crear el Server
	ph := database.NewPostgres(db, ctx)
	svc := service.NewService(ph)
	s, err := svc.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	//Cabeceras CORS:
	handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}}).Handler(s.Router())

	//Se pone a escuchar en el puerto TCP 8000 de localhost y llama al handler
	log.Fatal(http.ListenAndServe(":8000", handler))

	// Si nada falla, indica que el server esta en funcionamiento
	//fmt.Printf("Server running\n")
	log.Println("Server running")
}
