package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-final-project/datetask"
	"go-final-project/store"
	"log"
	"net/http"
	"strconv"
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

	task, err = checkAndEnrichTask(task)
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

func GetTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	task, err := taskStore.Get(id)

	if task.ID == 0 {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
		return
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, task, http.StatusOK)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
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

	if task.ID == 0 {
		makeHttpResponse(w, ErrorResponse{Error: "ID not specified"}, http.StatusBadRequest)
		return
	}

	task, err = checkAndEnrichTask(task)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	t, err := taskStore.Get(task.ID)

	if t.ID == 0 {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	err = taskStore.Edit(task)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	makeHttpResponse(w, "{}", http.StatusOK)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	t, err := taskStore.Get(id)

	if t.ID == 0 {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	if t.Repeat == "" {
		err = taskStore.Delete(id)

		if err != nil {
			makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		}
	} else {
		t.Date, err = datetask.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			makeHttpResponse(w, ErrorResponse{Error: "invalid date format"}, http.StatusInternalServerError)
			return
		}

		err = taskStore.Edit(t)
		if err != nil {
			makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		}
	}

	makeHttpResponse(w, "{}", http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	t, err := taskStore.Get(id)

	if t.ID == 0 {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	err = taskStore.Delete(id)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	makeHttpResponse(w, "{}", http.StatusOK)
}

func checkAndEnrichTask(t store.Task) (store.Task, error) {
	if t.Title == "" {
		return t, errors.New("task title is not specified")
	}

	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	} else {
		taskDate, err := time.Parse("20060102", t.Date)

		if err != nil {
			return t, errors.New("invalid date format")
		}

		if t.Repeat != "" {
			if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
				t.Date, err = datetask.NextDate(time.Now(), t.Date, t.Repeat)

				if err != nil {
					return t, errors.New("invalid date format")
				}
			}
		} else {
			if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
				t.Date = time.Now().Format("20060102")
			}
		}

	}
	return t, nil
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
