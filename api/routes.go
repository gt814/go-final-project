package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func GetRouter() *chi.Mux {
	//Initialize routing
	r := chi.NewRouter()

	// Define a directory for static files.
	workDir, _ := os.Getwd()
	webDir := http.Dir(workDir + "/web")

	// Set route for serving static files.
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// Set api routes
	path := "/api/"
	r.Get(path+"nextdate", NextDateHandler)
	r.Post(path+"task", AddTask)
	r.Get(path+"task", GetTask)
	r.Put(path+"task", EditTask)
	r.Delete(path+"task", DeleteTask)
	r.Post(path+"task/done", DoneTask)
	r.Get(path+"tasks", GetTasks)

	return r
}
