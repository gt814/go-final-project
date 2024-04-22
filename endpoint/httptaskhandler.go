package endpoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-final-project/datetask"
	"go-final-project/store"
	"log"
	"net/http"
	"time"
)

var taskStore store.TaskStore

func SetTaskStore(ts store.TaskStore) {
	taskStore = ts
}

type TaskIdResponse struct {
	ID int64 `json:"id"`
}

type TasksResponse struct {
	Tasks []store.Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
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

func AddTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task store.Task
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("read body, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf("unmarshal task, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	err = checkTask(task)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
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
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := taskStore.Add(task)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, TaskIdResponse{ID: id}, http.StatusCreated)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []store.Task

	tasks, err := taskStore.GetAll()

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []store.Task{}
	}

	makeHttpResponse(w, TasksResponse{Tasks: tasks}, http.StatusOK)
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

func makeHttpResponse(w http.ResponseWriter, response any, status int) {
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
