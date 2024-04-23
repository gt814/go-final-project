package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go-final-project/config"
	"go-final-project/endpoint"
	"go-final-project/store"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Initialize DB
	dbPath := config.GetDBFileAppPath()
	db, err := store.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	endpoint.SetTaskStore(store.NewTaskStore(db))
	//taskStore = store.NewTaskStore(db)

	//Initialize routing
	r := chi.NewRouter()

	// Define a directory for static files.
	workDir, _ := os.Getwd()
	webDir := http.Dir(workDir + "/web")

	// Set route for serving static files.
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// Set api routes
	r.Get("/api/nextdate", endpoint.NextDateHandler)
	r.Post("/api/task", endpoint.AddTask)
	r.Get("/api/task", endpoint.GetTask)
	r.Put("/api/task", endpoint.EditTask)
	r.Delete("/api/task", endpoint.DeleteTask)
	r.Post("/api/task/done", endpoint.DoneTask)
	r.Get("/api/tasks", endpoint.GetTasks)

	//Start web server.
	port := strconv.Itoa(config.GetPort())
	log.Println("Start listening on the port=", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Start web server, err = %w", err)
		return
	}
}
