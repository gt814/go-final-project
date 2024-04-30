package main

import (
	"go-final-project/api"
	"go-final-project/config"
	"go-final-project/service"
	"go-final-project/store"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize DB
	dbPath := config.GetDBFileAppPath()
	db, err := store.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service.SetTaskStore(store.NewTaskStore(db))

	//Initialize routing
	r := api.GetRouter()

	//Start web server.
	log.Println("Start listening on the port=", config.GetPort())
	if err := http.ListenAndServe(":"+config.GetPort(), r); err != nil {
		log.Fatal("Start web server, err = %w", err)
		return
	}
}
