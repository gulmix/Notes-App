package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gulmix/Notes-App/internal/handlers"
	"github.com/gulmix/Notes-App/internal/repo"
)

func main() {
	db, err := repo.New("postgres://postgres:psq@localhost:5432/notes")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()
	r := mux.NewRouter()
	r.HandleFunc("/notes", handlers.CreateHandler).Methods("POST")
	r.HandleFunc("/notes/{id}", handlers.GetHandler).Methods("GET")
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
