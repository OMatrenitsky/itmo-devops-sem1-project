package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"project_sem/internal/db"
	"project_sem/internal/handlers"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/v0/prices", handlers.PostPrices(database)).
		Methods(http.MethodPost)

	r.HandleFunc("/api/v0/prices", handlers.GetPrices(database)).
		Methods(http.MethodGet)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
