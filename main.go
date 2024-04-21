package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go-final-project/config"
	"go-final-project/datetask"
	"go-final-project/store"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var taskStore store.TaskStore

type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func main() {
	// Initialize DB
	dbPath := config.GetDBFileAppPath()
	db, err := store.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskStore = store.NewTaskStore(db)

	//Initialize routing
	r := chi.NewRouter()

	// Define a directory for static files.
	workDir, _ := os.Getwd()
	webDir := http.Dir(workDir + "/web")

	// Set route for serving static files.
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// Set api routes
	r.Get("/api/nextdate", nextDateHandler)
	r.Post("/api/task", addTask)

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
	_, err = datetask.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := datetask.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task store.Task
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("read body, err=%w", err)
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf("unmarshal task, err=%w", err)
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("task=", task)

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	err = checkTask(task)

	if err != nil {
		makeHttpResponse(w, Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	taskDate, _ := time.Parse("20060102", task.Date)

	if task.Repeat != "" {
		if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
			task.Date, err = datetask.NextDate(time.Now(), task.Date, task.Repeat)
		}
	} else {
		if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
			task.Date = time.Now().Format("20060102")
		}
	}

	if err != nil {
		makeHttpResponse(w, Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := taskStore.Add(task)
	if err != nil {
		makeHttpResponse(w, Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, Response{ID: id}, http.StatusCreated)
}

func checkTask(t store.Task) error {
	if t.Title == "" {
		return errors.New("task title is not specified")
	}

	_, err := time.Parse("20060102", t.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	return nil
}

func makeHttpResponse(w http.ResponseWriter, response Response, status int) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.WriteHeader(status)
		w.Write(jsonResponse)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
