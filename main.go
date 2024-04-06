package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-final-project/config"
	"net/http"
	"os"
	"strconv"
)

func main() {
	r := chi.NewRouter()

	// Define a directory for static files.
	workDir, _ := os.Getwd()
	webDir := http.Dir(workDir + "/web")

	// Set route for serving static files.
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	//Start web server.
	port := strconv.Itoa(config.GetPort())
	fmt.Println("start listening on the port=", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		err = fmt.Errorf("start web server, err = %w", err)
		fmt.Println(err)
		return
	}
}
