package main

import (
	"github.com/go-chi/chi/v5"
	"go-final-project/config"
	"go-final-project/store"
	"go-final-project/tasks"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

	// Set api routes
	r.Get("/api/nextdate", nextDateHandler)

	//Start web server.
	port := strconv.Itoa(config.GetPort())
	log.Println("Start listening on the port=", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Start web server, err = %w", err)
		return
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", dateParam)
	if err != nil {
		http.Error(w, "Invalid 'date' parameter", http.StatusBadRequest)
		return
	}

	log.Println("Параметры запуска NextDate now, dateParam, repeatParam", now, dateParam, repeatParam)
	_, err = tasks.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := tasks.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
