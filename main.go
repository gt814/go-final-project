package main

import (
	"github.com/go-chi/chi/v5"
	"go-final-project/config"
	"go-final-project/store"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	db, err := store.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	// Define a directory for static files.
	workDir, _ := os.Getwd()
	webDir := http.Dir(workDir + "/web")

	// Set route for serving static files.
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	//Start web server.
	port := strconv.Itoa(config.GetPort())
	log.Println("Start listening on the port=", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Start web server, err = %w", err)
		return
	}
}
