package main

import (
	_ "github.com/mattn/go-sqlite3"
	"go-final-project/api"
	"go-final-project/config"
	"go-final-project/store"
	"log"
	"net/http"
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

	api.SetTaskStore(store.NewTaskStore(db))

	//Initialize routing
	r := api.GetRouter()

	//Start web server.
	port := strconv.Itoa(config.GetPort())
	log.Println("Start listening on the port=", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Start web server, err = %w", err)
		return
	}
}
